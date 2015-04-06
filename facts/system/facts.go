package system

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

// Facts holds the system facts.
type Facts struct {
	Architecture   string     `json:"architecture"`
	BootID         string     `json:"boot_id"`
	Hostname       string     `json:"hostname"`
	Interfaces     Interfaces `json:"interfaces"`
	Kernel         string     `json:"kernel"`
	MachineID      string     `json:"machine_id"`
	OSRelease      OSRelease  `json:"os_release"`
	Virtualization string     `json:"virtualization"`

	mu sync.Mutex
}

// OSRelease holds the OS release facts.
type OSRelease struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	PrettyName string `json:"pretty_name"`
	Version    string `json:"version"`
	VersionID  string `json:"version_id"`
}

// Interfaces holds the network facts.
type Interfaces map[string]Interface

// Interface holds facts for a single interface.
type Interface struct {
	Name         string   `json:"name"`
	Index        int      `json:"index"`
	HardwareAddr string   `json:"hardware_address"`
	IpAddresses  []string `json:"ip_addresses"`
}

func Run() *Facts {
	facts := new(Facts)
	var wg sync.WaitGroup

	wg.Add(3)
	go facts.getOSRelease(&wg)
	go facts.getInterfaces(&wg)
	go facts.getHostnamectl(&wg)

	wg.Wait()
	return facts
}

func (f *Facts) getOSRelease(wg *sync.WaitGroup) {
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

func (f *Facts) getInterfaces(wg *sync.WaitGroup) {
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

func (f *Facts) getHostnamectl(wg *sync.WaitGroup) {
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
