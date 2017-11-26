# Servers

Run local http servers

## Config (MacOS)

Network -> Advanced -> Proxies -> Automatic Proxy Configuration -> "http://localhost:8080/proxy.pac"

## Usage

```
go build cmd/servers.go
./servers sample-config.json
```
