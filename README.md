# Terminus

Get facts about a Linux system. Parallel execution, structured output, remote API.

## Install

```
wget https://github.com/kelseyhightower/terminus/releases/download/v0.0.1/terminus
chmod +x terminus
```

## Usage

Terminus ships with a default set of facts that represent info about the system. Terminus also supports [custom facts](docs/custom-facts.md) and a [HTTP API](docs/api.md).

### Print a single fact

```
terminus --format '{{.System.BootID}}'
```

```
029b978a8d0b4ac48c5ca9c92956eeb6
```

### Print all facts

```
terminus
```
```
{
   "System": {
     "Architecture": "x86_64",
     "BootID": "e9e1695b-4c87-40fe-b701-d19dc262fc25",
     "Date": {
       "Unix": 1428781730,
       "UTC": "2015-04-11 19:48:50.291769492 +0000 UTC"
     },
     "Domainname": "(none)",
     "Hostname": "ubuntu",
     "Network": {
       "Interfaces": {
         "docker0": {
           "Name": "docker0",
           "Index": 3,
           "HardwareAddr": "56:84:7a:fe:97:99",
           "IpAddresses": [
             "172.17.42.1/16"
           ]
         },
         "eth0": {
           "Name": "eth0",
           "Index": 2,
           "HardwareAddr": "00:0c:29:ca:67:76",
           "IpAddresses": [
             "192.168.12.139/16",
             "fe80::20c:29ff:feca:6776/64"
           ]
         },
         "lo": {
           "Name": "lo",
           "Index": 1,
           "HardwareAddr": "",
           "IpAddresses": [
             "127.0.0.1/8",
             "::1/128"
           ]
         }
       }
     },
     "Kernel": {
       "Name": "Linux",
       "Release": "3.19.0-031900rc6-generic",
       "Version": "#201501261152 SMP Mon Jan 26 16:53:27 UTC 2015"
     },
     "MachineID": "3ca6d0646855f7cc6480630a54ac4a20",
     "Memory": {
       "Total": 1024004096,
       "Free": 727904256,
       "Shared": 684032,
       "Buffered": 24641536
     },
     "OSRelease": {
       "Name": "Ubuntu",
       "ID": "ubuntu",
       "PrettyName": "Ubuntu 14.10",
       "Version": "14.10 (Utopic Unicorn)",
       "VersionID": "14.10"
     },
     "Swap": {
       "Total": 4294963200,
       "Free": 4294963200
     },
     "Uptime": 12673
   }
}
```
