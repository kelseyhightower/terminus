# Terminus API

To enable the terminus HTTP API use the `-http` flag.

## Usage

### Server

```
terminus -http=":8080"
```

### Client

#### Get all facts

```
curl http://$SERVER_IP:8080/facts
```

#### Get a single fact

```
curl http://$SERVER_IP:8080/facts -d '{{.System.MachineID}}'
```
