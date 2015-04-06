# Terminus

Inspired by Puppet Labs Facter, but written in Go.

## Usage

```
terminus
```
```
{
   "system": {
     "architecture": "x86_64",
     "boot_id": "7e4bed2a4cb04f47b344c45c2aabc6f2",
     "hostname": "ubuntu",
     "interfaces": {
       "docker0": {
         "name": "docker0",
         "index": 3,
         "hardware_address": "56:84:7a:fe:97:99",
         "ip_addresses": [
           "172.17.42.1/16"
         ]
       },
       "eth0": {
         "name": "eth0",
         "index": 2,
         "hardware_address": "00:0c:29:ca:67:76",
         "ip_addresses": [
           "192.168.12.139/16",
           "fe80::20c:29ff:feca:6776/64"
         ]
       },
       "lo": {
         "name": "lo",
         "index": 1,
         "hardware_address": "",
         "ip_addresses": [
           "127.0.0.1/8",
           "::1/128"
         ]
       }
     },
     "kernel": "Linux 3.19.0-031900rc6-generic",
     "machine_id": "3ca6d0646855f7cc6480630a54ac4a20",
     "os_release": {
       "name": "Ubuntu",
       "id": "ubuntu",
       "pretty_name": "Ubuntu 14.10",
       "version": "14.10 (Utopic Unicorn)",
       "version_id": "14.10"
     },
     "virtualization": ""
   }
 }
```
