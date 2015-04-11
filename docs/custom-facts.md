# Custom Facts

There are two types of custom facts: executable and static. Custom facts are stored in the external facts directory, which defaults to `/etc/terminus/facts.d`. Use the `-external-facts-dir` flag to specify a different location.

## Executable Facts

Executable facts can be written in any language and reside under the external facts directory with the executable bit set.

### Example

```
sudo vim /etc/terminus/facts.d/date
```

```
#!/bin/bash

echo "{\"Now\": \"$(date)\"}"
exit 0
```

```
sudo chmod +x /etc/terminus/facts.d/date
```

```
terminus -format '{{.date.Now}}'
```

```
Sat Apr 11 13:38:26 PDT 2015
```

## Static Facts

Static facts must be in the JSON format and reside under the external facts directory with a `.json` file extension.

### Example

```
sudo vim /etc/terminus/facts.d/docker.json
```

```
{
  "ClientAPIVersion": "1.16",
  "ClientOSArch": "linux/amd64",
  "ClientVersion": "1.5.0",
  "ServerAPIVersion": "1.16",
  "ServerOSArch": "linux/amd64",
  "ServerVersion": "1.5.0"
}
```

```
terminus -format '{{.docker.ServerAPIVersion}}'
```

```
1.16
```
