package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kelseyhightower/terminus/facts"
	"golang.org/x/sys/unix"
)

// Constants
const (
	// LINUX_SYSINFO_LOADS_SCALE has been described elsewhere as a "magic" number.
	// It reverts the calculation of "load << (SI_LOAD_SHIFT - FSHIFT)" done in the original load calculation.
	LINUX_SYSINFO_LOADS_SCALE float64 = 65536.0
)

// SystemFacts holds the system facts.
type SystemFacts struct {
	Architecture string
	BootID       string
	Date         Date
	Domainname   string
	Hostname     string
	Network      Network
	Kernel       Kernel
	MachineID    string
	Memory       Memory
	OSRelease    OSRelease
	Swap         Swap
	Uptime       int64
	LoadAverage  LoadAverage
	FileSystems  FileSystems

	mu sync.Mutex
}

// Holds the load average facts.
type LoadAverage struct {
	One  string
	Five string
	Ten  string
}

// Date holds the date facts.
type Date struct {
	Unix int64
	UTC  string
}

// Swap holds the swap facts.
type Swap struct {
	Total uint64
	Free  uint64
}

// OSRelease holds the OS release facts.
type OSRelease struct {
	Name       string
	ID         string
	PrettyName string
	Version    string
	VersionID  string
}

// Kernel holds the kernel facts.
type Kernel struct {
	Name    string
	Release string
	Version string
}

// Memory holds the memory facts.
type Memory struct {
	Total    uint64
	Free     uint64
	Shared   uint64
	Buffered uint64
}

// Network holds the network facts.
type Network struct {
	Interfaces Interfaces
}

// Interfaces holds the interface facts.
type Interfaces map[string]Interface

// Interface holds facts for a single interface.
type Interface struct {
	Name         string
	Index        int
	HardwareAddr string
	IpAddresses  []string
	Ip4Addresses []Ip4Address
	Ip6Addresses []Ip6Address
}

type Ip4Address struct {
	CIDR    string
	Ip      string
	Netmask string
}

type Ip6Address struct {
	CIDR   string
	Ip     string
	Prefix int
}

// FileSystems holds the Filesystem facts.
type FileSystems map[string]FileSystem

// FileSystem holds facts for a filesystem (man fstab).
type FileSystem struct {
	Device     string
	MountPoint string
	Type       string
	Options    []string
	DumpFreq   uint64
	PassNo     uint64
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

	wg.Add(8)
	go facts.getOSRelease(&wg)
	go facts.getInterfaces(&wg)
	go facts.getBootID(&wg)
	go facts.getMachineID(&wg)
	go facts.getUname(&wg)
	go facts.getSysInfo(&wg)
	go facts.getDate(&wg)
	go facts.getFileSystems(&wg)

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

	f.LoadAverage.One = fmt.Sprintf("%.2f", float64(info.Loads[0])/LINUX_SYSINFO_LOADS_SCALE)
	f.LoadAverage.Five = fmt.Sprintf("%.2f", float64(info.Loads[1])/LINUX_SYSINFO_LOADS_SCALE)
	f.LoadAverage.Ten = fmt.Sprintf("%.2f", float64(info.Loads[2])/LINUX_SYSINFO_LOADS_SCALE)

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
		if len(columns) > 1 {
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
		ip4addrs := make([]Ip4Address, 0)
		ip6addrs := make([]Ip6Address, 0)

		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, ip := range addrs {
			cidr := ip.String()
			ipaddreses = append(ipaddreses, cidr)
			ip, ipnet, _ := net.ParseCIDR(cidr)
			if ip.To4() != nil {
				ip4addrs = append(ip4addrs, Ip4Address{cidr, ip.String(), toNetmask(ipnet.Mask)})
				continue
			}
			if ip.To16() != nil {
				ones, _ := ipnet.Mask.Size()
				ip6addrs = append(ip6addrs, Ip6Address{cidr, ip.String(), ones})
			}
		}
		m[i.Name] = Interface{
			Name:         i.Name,
			Index:        i.Index,
			HardwareAddr: i.HardwareAddr.String(),
			IpAddresses:  ipaddreses,
			Ip4Addresses: ip4addrs,
			Ip6Addresses: ip6addrs,
		}
	}
	f.Network.Interfaces = m
	return
}

func toNetmask(m net.IPMask) string {
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
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

func (f *SystemFacts) getFileSystems(wg *sync.WaitGroup) {
	defer wg.Done()

	mtab, err := ioutil.ReadFile("/etc/mtab")
	if err != nil {
		log.Println(err.Error())
		return
	}

	fsMap := make(FileSystems)

	f.mu.Lock()
	defer f.mu.Unlock()

	s := bufio.NewScanner(bytes.NewBuffer(mtab))
	for s.Scan() {
		line := s.Text()
		if string(line[0]) == "#" {
			continue
		}
		fields := strings.Fields(s.Text())
		fs := FileSystem{}
		fs.Device = fields[0]
		fs.MountPoint = fields[1]
		fs.Type = fields[2]
		fs.Options = strings.Split(fields[3], ",")
		dumpFreq, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fs.DumpFreq = dumpFreq

		passNo, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fs.PassNo = passNo

		fsMap[fs.Device] = fs
	}

	f.FileSystems = fsMap
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
