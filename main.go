package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net"
    "os"
    "sync"
)

type Config struct {
    Local struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"local"`
    Remote struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"remote"`
}

var (
    activeConnections = make(map[string]bool)
    mu                sync.Mutex
)

func main() {
    // Load config
    file, err := os.Open("config.json")
    if err != nil {
        fmt.Println("Error opening config.json:", err)
        return
    }
    defer file.Close()

    var config Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        fmt.Println("Error parsing config.json:", err)
        return
    }

    localAddr := fmt.Sprintf("%s:%d", config.Local.Host, config.Local.Port)
    remoteAddr := fmt.Sprintf("%s:%d", config.Remote.Host, config.Remote.Port)

    ln, err := net.Listen("tcp", localAddr)
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer ln.Close()

    fmt.Printf("Proxy server listening on %s...\n", localAddr)

    for {
        localConn, err := ln.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }

        clientIP, _, _ := net.SplitHostPort(localConn.RemoteAddr().String())

        mu.Lock()
        if activeConnections[clientIP] {
            mu.Unlock()
            fmt.Printf("Connection from %s rejected (already connected)\n", clientIP)
            localConn.Close()
            continue
        }
        activeConnections[clientIP] = true
        mu.Unlock()

        fmt.Printf("Incoming connection from: %s\n", localConn.RemoteAddr())

        go handleConnection(localConn, remoteAddr, clientIP)
    }
}

func handleConnection(localConn net.Conn, remoteAddr, clientIP string) {
    defer func() {
        localConn.Close()
        mu.Lock()
        delete(activeConnections, clientIP)
        mu.Unlock()
        fmt.Printf("Connection from %s ended\n", clientIP)
    }()

    remoteConn, err := net.Dial("tcp", remoteAddr)
    if err != nil {
        fmt.Printf("Error connecting to remote server: %v\n", err)
        return
    }
    defer remoteConn.Close()

    fmt.Printf("Connected to remote server: %s\n", remoteAddr)

    // Bidirectional copy
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        io.Copy(remoteConn, localConn)
        remoteConn.(*net.TCPConn).CloseWrite()
    }()
    go func() {
        defer wg.Done()
        io.Copy(localConn, remoteConn)
        localConn.(*net.TCPConn).CloseWrite()
    }()

    wg.Wait()
}