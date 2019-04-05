package main

import (
    "flag"
    "fmt"
    "net"
    "log"
    "regexp"
)

var reg,_ = regexp.Compile("[^a-zA-Z0-9]+")

type Client struct {
    id int
    name string
    co net.Conn
}

var clients map[int]Client
var cntID = 0

func main() {
    // init map
    clients = make(map[int]Client)
    // custom port option
    portInt := flag.Int("p",8080,"port")
    flag.Parse()
    portStr := fmt.Sprintf(":%d",*portInt)
    // listener
    ln, err := net.Listen("tcp", portStr)
    if err != nil {
        log.Fatal(err)
        return
    }
    fmt.Printf("listening to port %d\n",*portInt)
    for { // new connection
        co, err := ln.Accept()
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        go clientHandler(co)
    }
}

func clientHandler(co net.Conn) {
    buf := make([]byte, 1024)
    id := cntID
    cntID++ // race condition?
    co.Write([]byte("who are you? "))
    { // init
        n, err := co.Read(buf)
        if err != nil {
            return
        }
        name := fmt.Sprintf("%s", buf[:n-1]) // without br
        co.Write([]byte(fmt.Sprintf("Welcome %s!\n", name)))
        clients[id] = Client{id, name, co}
    }
    for { // listen
        n, err := co.Read(buf)
        if err != nil {
            delete(clients, id)
            co.Close()
            return
        }
        spread(id, buf[:n])
    }
}

func spread(id int, msg []byte) {
    sign := clients[id].name + "> "
    buf := append([]byte(sign), msg...) // variadic
    for _,c := range clients {
        if c.id == id {
            continue
        }
        c.co.Write(buf)
    }
}

