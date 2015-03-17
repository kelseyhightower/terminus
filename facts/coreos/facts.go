package coreos

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Facts map[string]string

func Run() Facts {
	facts := make(Facts)
	facts.osRelease()
	facts.interfaces()
	facts.hostnamectl()
	return facts
}

func (f Facts) osRelease() error {
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
			f["/coreos/name"] = value
		case "ID":
			f["/coreos/id"] = value
		case "PRETTY_NAME":
			f["/coreos/pretty_name"] = value
		case "VERSION":
			f["/coreos/version"] = value
		case "VERSION_ID":
			f["/coreos/version_id"] = value
		}
	}
	return nil
}

func (f Facts) interfaces() error {
	ls, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range ls {
		ipaddress := make([]string, 0)
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		for _, ip := range addrs {
			ipaddress = append(ipaddress, ip.String())
		}
		f[path.Join("/coreos/interfaces", i.Name, "index")] = strconv.Itoa(i.Index)
		f[path.Join("/coreos/interfaces", i.Name, "hardware_addr")] = i.HardwareAddr.String()
		f[path.Join("/coreos/interfaces", i.Name, "ip_addresses")] = strings.Join(ipaddress, ",")
	}
	return nil
}

func (f Facts) hostnamectl() error {
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
			f["/coreos/architecture"] = value
		case "Static hostname":
			f["/coreos/hostname"] = value
		case "Machine ID":
			f["/coreos/machine_id"] = value
		case "Boot ID":
			f["/coreos/boot_id"] = value
		case "Virtualization":
			f["/coreos/virtualization"] = value
		case "Kernel":
			f["/coreos/kernel"] = value
		}
	}
	return nil
}
