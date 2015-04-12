package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kelseyhightower/terminus/facts"
	"golang.org/x/sys/unix"
)

// SystemFacts holds the system facts.
type SystemFacts struct {
	Architecture string      `json:"architecture"`
	BootID       string      `json:"bootID"`
	Date         Date        `json:"date"`
	Domainname   string      `json:"domainName"`
	Hostname     string      `json:"hostName"`
	Network      Network     `json:"network"`
	Kernel       Kernel      `json:"kernel"`
	MachineID    string      `json:"machineID"`
	Memory       Memory      `json:"memory"`
	OSRelease    OSRelease   `json:"osRelease"`
	Swap         Swap        `json:"swap"`
	Uptime       int64       `json:"uptime"`
	LoadAverage  LoadAverage `json:"loadAverage"`

	mu sync.Mutex
}

// Holds the load average facts.
type LoadAverage struct {
	One  uint64 `json:"one"`
	Five uint64 `json:"five"`
	Ten  uint64 `json:"ten"`
}

// Date holds the date facts.
type Date struct {
	Unix int64  `json:"unix"`
	UTC  string `json:"UTC"`
}

// Swap holds the swap facts.
type Swap struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

// OSRelease holds the OS release facts.
type OSRelease struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	PrettyName string `json:"prettyName"`
	Version    string `json:"version"`
	VersionID  string `json:"versionID"`
}

// Kernel holds the kernel facts.
type Kernel struct {
	Name    string `json:"name"`
	Release string `json:"release"`
	Version string `json:"version"`
}

// Memory holds the memory facts.
type Memory struct {
	Total    uint64 `json:"total"`
	Free     uint64 `json:"free"`
	Shared   uint64 `json:"shared"`
	Buffered uint64 `json:"buffered"`
}

// Network holds the network facts.
type Network struct {
	Interfaces Interfaces `json:"interfaces"`
}

// Interfaces holds the interface facts.
type Interfaces map[string]Interface

// Interface holds facts for a single interface.
type Interface struct {
	Name         string   `json:"name"`
	Index        int      `json:"index"`
	HardwareAddr string   `json:"hardwareAddr"`
	IpAddresses  []string `json:"ipAddresses"`
}

func getFacts() *facts.Facts {
	f := facts.New()
	systemFacts := getSystemFacts()
	f.Add("System", systemFacts)
	processExternalFacts(externalFactsDir, f)
	return f
}

func getSystemFacts() *SystemFacts {
	facts := new(SystemFacts)
	var wg sync.WaitGroup

	wg.Add(7)
	go facts.getOSRelease(&wg)
	go facts.getInterfaces(&wg)
	go facts.getBootID(&wg)
	go facts.getMachineID(&wg)
	go facts.getUname(&wg)
	go facts.getSysInfo(&wg)
	go facts.getDate(&wg)

	wg.Wait()
	return facts
}

func (f *SystemFacts) getDate(wg *sync.WaitGroup) {
	defer wg.Done()

	now := time.Now()
	f.Date.Unix = now.Unix()
	f.Date.UTC = now.UTC().String()

	return
}

func (f *SystemFacts) getSysInfo(wg *sync.WaitGroup) {
	defer wg.Done()

	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		log.Println(err.Error())
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.Memory.Total = info.Totalram
	f.Memory.Free = info.Freeram
	f.Memory.Shared = info.Sharedram
	f.Memory.Buffered = info.Bufferram

	f.Swap.Total = info.Totalswap
	f.Swap.Free = info.Freeswap

	f.Uptime = info.Uptime

	f.LoadAverage.One = info.Loads[0]
	f.LoadAverage.Five = info.Loads[1]
	f.LoadAverage.Ten = info.Loads[2]

	return
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

func (f *SystemFacts) getMachineID(wg *sync.WaitGroup) {
	defer wg.Done()
	machineIDFile, err := os.Open("/etc/machine-id")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer machineIDFile.Close()
	data, err := ioutil.ReadAll(machineIDFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.MachineID = strings.TrimSpace(string(data))
	return
}

func (f *SystemFacts) getBootID(wg *sync.WaitGroup) {
	defer wg.Done()
	bootIDFile, err := os.Open("/proc/sys/kernel/random/boot_id")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer bootIDFile.Close()
	data, err := ioutil.ReadAll(bootIDFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.BootID = strings.TrimSpace(string(data))
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
	f.Network.Interfaces = m
	return
}

func (f *SystemFacts) getUname(wg *sync.WaitGroup) {
	defer wg.Done()

	var buf unix.Utsname
	err := unix.Uname(&buf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.Domainname = charsToString(buf.Domainname)
	f.Architecture = charsToString(buf.Machine)
	f.Hostname = charsToString(buf.Nodename)
	f.Kernel.Name = charsToString(buf.Sysname)
	f.Kernel.Release = charsToString(buf.Release)
	f.Kernel.Version = charsToString(buf.Version)
	return
}

func processExternalFacts(dir string, f *facts.Facts) {
	d, err := os.Open(dir)
	if err != nil {
		log.Println(err)
		return
	}
	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}

	executableFacts := make([]string, 0)
	staticFacts := make([]string, 0)

	for _, fi := range files {
		name := filepath.Join(dir, fi.Name())
		if isExecutable(fi) {
			executableFacts = append(executableFacts, name)
			continue
		}
		if strings.HasSuffix(name, ".json") {
			staticFacts = append(staticFacts, name)
		}
	}

	var wg sync.WaitGroup
	for _, p := range staticFacts {
		p := p
		wg.Add(1)
		go factsFromFile(p, f, &wg)
	}
	for _, p := range executableFacts {
		p := p
		wg.Add(1)
		go factsFromExec(p, f, &wg)
	}
	wg.Wait()
}

func factsFromFile(path string, f *facts.Facts, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}
	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
		return
	}
	f.Add(strings.TrimSuffix(filepath.Base(path), ".json"), result)
}

func factsFromExec(path string, f *facts.Facts, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := exec.Command(path).Output()
	if err != nil {
		log.Println(err)
		return
	}
	var result interface{}
	err = json.Unmarshal(out, &result)
	if err != nil {
		log.Println(err)
		return
	}
	f.Add(filepath.Base(path), result)
}

func isExecutable(fi os.FileInfo) bool {
	if m := fi.Mode(); !m.IsDir() && m&0111 != 0 {
		return true
	}
	return false
}

func charsToString(ca [65]int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}
