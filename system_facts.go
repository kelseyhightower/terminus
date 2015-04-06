package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// SystemFacts holds the system facts.
type SystemFacts struct {
	Architecture   string
	BootID         string
	Hostname       string
	Interfaces     Interfaces
	Kernel         string
	MachineID      string
	OSRelease      OSRelease
	Virtualization string

	mu sync.Mutex
}

// OSRelease holds the OS release facts.
type OSRelease struct {
	Name       string
	ID         string
	PrettyName string
	Version    string
	VersionID  string
}

// Interfaces holds the network facts.
type Interfaces map[string]Interface

// Interface holds facts for a single interface.
type Interface struct {
	Name         string
	Index        int
	HardwareAddr string
	IpAddresses  []string
}

func getSystemFacts() *SystemFacts {
	facts := new(SystemFacts)
	var wg sync.WaitGroup

	wg.Add(3)
	go facts.getOSRelease(&wg)
	go facts.getInterfaces(&wg)
	go facts.getHostnamectl(&wg)

	wg.Wait()
	return facts
}

func (f *SystemFacts) getOSRelease(wg *sync.WaitGroup) {
	defer wg.Done()
	osReleaseFile, err := os.Open("/etc/os-release")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer osReleaseFile.Close()

	f.mu.Lock()
	defer f.mu.Unlock()
	scanner := bufio.NewScanner(osReleaseFile)
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), "=")
		key := columns[0]
		value := strings.Trim(strings.TrimSpace(columns[1]), `"`)
		switch key {
		case "NAME":
			f.OSRelease.Name = value
		case "ID":
			f.OSRelease.ID = value
		case "PRETTY_NAME":
			f.OSRelease.PrettyName = value
		case "VERSION":
			f.OSRelease.Version = value
		case "VERSION_ID":
			f.OSRelease.VersionID = value
		}
	}
	return
}

func (f *SystemFacts) getInterfaces(wg *sync.WaitGroup) {
	defer wg.Done()
	ls, err := net.Interfaces()
	if err != nil {
		log.Println(err.Error())
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	m := make(Interfaces)
	for _, i := range ls {
		ipaddreses := make([]string, 0)
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, ip := range addrs {
			ipaddreses = append(ipaddreses, ip.String())
		}
		m[i.Name] = Interface{
			Name:         i.Name,
			Index:        i.Index,
			HardwareAddr: i.HardwareAddr.String(),
			IpAddresses:  ipaddreses,
		}
	}
	f.Interfaces = m
	return
}

func (f *SystemFacts) getHostnamectl(wg *sync.WaitGroup) {
	defer wg.Done()

	var out bytes.Buffer
	cmd := exec.Command("/usr/bin/hostnamectl")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), ":")
		key := strings.TrimSpace(columns[0])
		value := strings.TrimSpace(columns[1])
		switch key {
		case "Architecture":
			f.Architecture = value
		case "Static hostname":
			f.Hostname = value
		case "Machine ID":
			f.MachineID = value
		case "Boot ID":
			f.BootID = value
		case "Virtualization":
			f.Virtualization = value
		case "Kernel":
			f.Kernel = value
		}
	}
	return
}
