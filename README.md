# Facter

Puppet Labs facter written in Go for CoreOS.

## Usage

```
facter

{
  "coreos": {
    "architecture": "x86-64",
    "boot_id": "cdb46530199c4096aab1878c43780d7d",
    "id": "coreos",
    "interfaces": [
      {
        "index": 1,
        "name": "lo",
        "hardware_address": "",
        "ip_address": [
          "127.0.0.1/8",
          "::1/128"
        ]
      },
      {
        "index": 2,
        "name": "ens33",
        "hardware_address": "00:0c:29:83:a9:55",
        "ip_address": [
          "192.168.12.11/24",
          "fe80::20c:29ff:fe83:a955/64"
        ]
      },
      {
        "index": 3,
        "name": "docker0",
        "hardware_address": "56:84:7a:fe:97:99",
        "ip_address": [
          "172.17.42.1/16"
        ]
      }
    ],
    "hostname": "core1.example.com",
    "kernel": "Linux 3.15.8+",
    "machine_id": "62c44d9b990343049317a85f42592a5b",
    "name": "CoreOS",
    "pretty_name": "CoreOS 410.0.0",
    "version": "410.0.0",
    "version_id": "410.0.0",
    "virtualization": "vmware"
  }
}
```
