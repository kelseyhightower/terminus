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
$ terminus --format '{{.System.BootID}}'
029b978a8d0b4ac48c5ca9c92956eeb6
```

or

```
$ terminus System.Network.Interfaces.eth0.Ip6Addresses.0.Ip
fe80::f816:3eff:fead:8549
```

### Print all facts

```
terminus
```
```
{
   "System": {
     "Architecture": "x86_64",
     "BootID": "87c81966-9e09-4627-b949-6320ae09ecfa",
     "Date": {
       "Unix": 1430413811,
       "UTC": "2015-04-30 17:10:11.736735078 +0000 UTC"
     },
     "Domainname": "(none)",
     "Hostname": "etcd",
     "Network": {
       "Interfaces": {
         "eno16777736": {
           "Name": "eno16777736",
           "Index": 2,
           "HardwareAddr": "00:0c:29:d6:9c:9a",
           "IpAddresses": [
             "192.168.12.10/16",
             "fe80::20c:29ff:fed6:9c9a/64"
           ],
           "Ip4Addresses": [
             {
               "CIDR": "192.168.12.10/16",
               "Ip": "192.168.12.10",
               "Netmask": "255.255.0.0"
             }
           ],
           "Ip6Addresses": [
             {
               "CIDR": "fe80::20c:29ff:fed6:9c9a/64",
               "Ip": "fe80::20c:29ff:fed6:9c9a",
               "Prefix": 64
             }
           ]
         },
         "lo": {
           "Name": "lo",
           "Index": 1,
           "HardwareAddr": "",
           "IpAddresses": [
             "127.0.0.1/8",
             "::1/128"
           ],
           "Ip4Addresses": [
             {
               "CIDR": "127.0.0.1/8",
               "Ip": "127.0.0.1",
               "Netmask": "255.0.0.0"
             }
           ],
           "Ip6Addresses": [
             {
               "CIDR": "::1/128",
               "Ip": "::1",
               "Prefix": 128
             }
           ]
         }
       }
     },
     "Kernel": {
       "Name": "Linux",
       "Release": "4.0.0",
       "Version": "#2 SMP Wed Apr 22 23:43:22 UTC 2015"
     },
     "MachineID": "677f2a9b43c343aa993ef4a282ba2f05",
     "Memory": {
       "Total": 1029615616,
       "Free": 683864064,
       "Shared": 184950784,
       "Buffered": 23953408
     },
     "OSRelease": {
       "Name": "CoreOS",
       "ID": "coreos",
       "PrettyName": "CoreOS 660.0.0",
       "Version": "660.0.0",
       "VersionID": "660.0.0"
     },
     "Swap": {
       "Total": 0,
       "Free": 0
     },
     "Uptime": 4927
   }
 }
```
