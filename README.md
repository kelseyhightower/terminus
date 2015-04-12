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
   "system": {
     "architecture": "x86_64",
     "bootID": "e9e1695b-4c87-40fe-b701-d19dc262fc25",
     "date": {
       "unix": 1428781730,
       "UTC": "2015-04-11 19:48:50.291769492 +0000 UTC"
     },
     "domainName": "(none)",
     "hostName": "ubuntu",
     "network": {
       "interfaces": {
         "docker0": {
           "name": "docker0",
           "index": 3,
           "hardwareAddr": "56:84:7a:fe:97:99",
           "ipAddresses": [
             "172.17.42.1/16"
           ]
         },
         "eth0": {
           "name": "eth0",
           "index": 2,
           "hardwareAddr": "00:0c:29:ca:67:76",
           "ipAddresses": [
             "192.168.12.139/16",
             "fe80::20c:29ff:feca:6776/64"
           ]
         },
         "lo": {
           "name": "lo",
           "index": 1,
           "hardwareAddr": "",
           "ipAddresses": [
             "127.0.0.1/8",
             "::1/128"
           ]
         }
       }
     },
     "kernel": {
       "name": "Linux",
       "release": "3.19.0-031900rc6-generic",
       "version": "#201501261152 SMP Mon Jan 26 16:53:27 UTC 2015"
     },
     "machineID": "3ca6d0646855f7cc6480630a54ac4a20",
     "memory": {
       "total": 1024004096,
       "free": 727904256,
       "shared": 684032,
       "buffered": 24641536
     },
     "osRelease": {
       "name": "Ubuntu",
       "id": "ubuntu",
       "prettyName": "Ubuntu 14.10",
       "version": "14.10 (Utopic Unicorn)",
       "versionID": "14.10"
     },
     "swap": {
       "total": 4294963200,
       "free": 4294963200
     },
     "uptime": 12673
   }
}
```
