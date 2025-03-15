# dnsupd - DynDNS update daemon
[![Go Report Card](https://goreportcard.com/badge/github.com/x-way/dnsupd)](https://goreportcard.com/report/github.com/x-way/dnsupd)

dnsupd - a small and simple DynDNS server

dnsupd listens for HTTP requests conforming to the DynDNS format and sends out DNS updates to an authoritative DNS server.
The DNS zone, server and TSIG credentials as well as the useraccounts and hostnames are provided via a config file.

## Installation
Either install the go package
```
# go install github.com/x-way/dnsupd@latest
```
Or alternatively install the docker image
```
# docker pull ghcr.io/x-way/dnsupd:latest
```

## Usage
Run the go binary from your local path
```
# dnsupd -f dnsupd.json
```
Or run the docker image while passing the config file as volume
```
# docker run -v $(pwd)/dnsupd.json:/etc/dnsupd.json ghcr.io/x-way/dnsupd:latest
```

## Configuration

dnsupd reads its configuration from the config file at `/etc/dnsupd.json` (default location, can be changed with the `-f` flag).

Sample config:
```
{
  "Zone": "dyn-zone.example.com",
  "Server": "ns.example.com",
  "TsigName": "my.tsig.name.example.com.",
  "TsigSecret": "Base64encodedsecret==",
  "TsigAlgorithm": "hmac-sha512.",
  "Hosts": [
    {
      "Hostname": "host1",
      "User": "user1",
      "Password": "password1"
    },
    {
      "Hostname": "host2",
      "User": "user2",
      "Password": "password2"
    }
  ]
}
```
