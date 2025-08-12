# Golang Server Proxy

This project is a TCP proxy server written in Go. It forwards connections from a local address to a remote server, allowing only one active connection per client IP at a time.

## Features
- Configurable local and remote addresses via `config.json`
- Only one connection per client IP is allowed simultaneously
- Bidirectional proxying between client and remote server

## Setup Instructions

### 1. Prerequisites
- Go 1.18 or newer installed: https://golang.org/dl/

### 2. Clone the Repository
```
git clone https://github.com/yrludev/Golang-Server-Proxy.git
cd Golang-Server-Proxy
```

### 3. Configuration
Edit the `assets/config.json` file to set your local and remote addresses and ports:
```json
{
  "local": {
    "host": "127.0.0.1",
    "port": 8080
  },
  "remote": {
    "host": "example.com",
    "port": 80
  }
}
```

- `local.host` and `local.port`: The address and port the proxy will listen on.
- `remote.host` and `remote.port`: The address and port the proxy will forward connections to.

### 4. Build and Run

```
go run main.go
```

Or build a binary:

```
go build -o proxy main.go
./proxy
```

### 5. Usage
- Connect to the proxy using the local address and port you specified in `config.json`.
- Only one connection per client IP is allowed at a time.

## Notes
- Make sure the `assets/config.json` file exists and is properly formatted.
- Adjust the maximum number of connections in the code if needed.

## License
YRLU.dev
