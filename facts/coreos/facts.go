package coreos

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

type Facts struct {
	Architecture   string      `json:"architecture" yaml:"architecture"`
	BootId         string      `json:"boot_id" yaml:"boot_id"`
	Id             string      `json:"id" yaml:"id"`
	Interfaces     []Interface `json:"interfaces yaml:"interfaces"`
	Hostname       string      `json:"hostname" yaml:"hostname"`
	Kernel         string      `json:"kernel" yaml:"kernel"`
	MachineId      string      `json:"machine_id" yaml:"machine_id"`
	Name           string      `json:"name" yaml:"name"`
	PrettyName     string      `json:"pretty_name" yaml:"pretty_name"`
	Version        string      `json:"version" yaml:"version"`
	VersionId      string      `json:"version_id" yaml:"version_id"`
	Virtualization string      `json:"virtualization" yaml:"virtualization"`
}

type Interface struct {
	Index        int      `json:"index"`
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardware_address"`
	IPAddress    []string `json:"ip_address"`
}

func Run() *Facts {
	facts := &Facts{}
	facts.osRelease()
	facts.interfaces()
	facts.hostnamectl()
	return facts
}

func (f *Facts) osRelease() error {
	osReleaseFile, err := os.Open("/etc/os-release")
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer osReleaseFile.Close()
	scanner := bufio.NewScanner(osReleaseFile)
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), "=")
		key := columns[0]
		value := strings.Trim(strings.TrimSpace(columns[1]), `"`)
		switch key {
		case "NAME":
			f.Name = value
		case "ID":
			f.Id = value
		case "PRETTY_NAME":
			f.PrettyName = value
		case "VERSION":
			f.Version = value
		case "VERSION_ID":
			f.VersionId = value
		}
	}
	return nil
}

func (f *Facts) interfaces() error {
	interfaces := make([]Interface, 0)
	ls, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range ls {
		ipaddress := make([]string, 0)
		iface := Interface{
			Index:        i.Index,
			Name:         i.Name,
			HardwareAddr: i.HardwareAddr.String(),
		}
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		for _, ip := range addrs {
			ipaddress = append(ipaddress, ip.String())
		}
		iface.IPAddress = ipaddress
		interfaces = append(interfaces, iface)
	}
	f.Interfaces = interfaces
	return nil
}

func (f *Facts) hostnamectl() error {
	var out bytes.Buffer
	cmd := exec.Command("/usr/bin/hostnamectl")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
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
			f.MachineId = value
		case "Boot ID":
			f.BootId = value
		case "Virtualization":
			f.Virtualization = value
		case "Kernel":
			f.Kernel = value
		}
	}
	return nil
}
