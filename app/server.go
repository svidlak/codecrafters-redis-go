package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type RedisServer struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan []byte
}

var Set map[string]string

func NewSever(listenAddr string) *RedisServer {
	return &RedisServer{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan []byte),
	}
}

func (rs *RedisServer) Start() error {
	ln, err := net.Listen("tcp", rs.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	Set = make(map[string]string)

	fmt.Println("redis server starting on port", rs.listenAddr)
	rs.ln = ln
	go rs.acceptConnectionsLoop()

	<-rs.quitch
	close(rs.msgch)

	return nil
}

func (rs *RedisServer) acceptConnectionsLoop() {
	for {
		conn, err := rs.ln.Accept()
		if err != nil {
			fmt.Print("accept connectiob error:", err)
			continue
		}
		fmt.Println("new client connected")

		go rs.readConnectionMessages(conn)
	}
}

func (rs *RedisServer) readConnectionMessages(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 2048)

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read connection error:", err)
			return
		}

		msg := buf[:n]
		response := rs.parseMessage(msg)
		response = response + "\r\n"

		conn.Write([]byte(response))
	}
}

func (rs *RedisServer) parseMessage(message []byte) string {
	msg := strings.ToLower(strings.TrimSpace(string(message)))
	splitMsg := strings.Split(msg, "\r\n")

	fmt.Printf("%#v\n", splitMsg)

	command := splitMsg[2]

	if command == "ping" {
		return "+PONG"
	}
	if command == "echo" {
		if len(splitMsg) >= 4 {
			return "+" + splitMsg[4]
		}
	}
	if command == "set" {
		if len(splitMsg) >= 6 {
			key := splitMsg[4]
			val := splitMsg[6]

			if len(splitMsg) >= 10 {
				expirationTime, err := strconv.Atoi(splitMsg[10])
				if err == nil {
					go time.AfterFunc(time.Duration(expirationTime)*time.Millisecond, func() {
						delete(Set, key)
					})
				}
			}
			Set[key] = val

			return "+OK"
		}
	}
	if command == "get" {
		if len(splitMsg) >= 4 {
			key := splitMsg[4]
			val, ok := Set[key]
			if ok {
				return "+" + val
			}
			return "$-1"
		}
	}

	return "-ERROR"

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	redisSever := NewSever(":6379")
	log.Fatal(redisSever.Start())
}
