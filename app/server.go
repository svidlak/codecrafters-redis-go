package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	config "github.com/codecrafters-io/redis-starter-go"
	commands "github.com/codecrafters-io/redis-starter-go/internal"
)

type RedisServer struct {
	ln     net.Listener
	quitch chan struct{}
	msgch  chan []byte
	config.Config
}

func NewSever() *RedisServer {
	return &RedisServer{
		quitch: make(chan struct{}),
		msgch:  make(chan []byte),
		Config: config.Configs,
	}
}

func (rs *RedisServer) Start() error {
	ln, err := net.Listen("tcp", rs.Config.ListenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	commands.Set = make(map[string]string)

	fmt.Println("redis server starting on port", rs.Config.ListenAddr)
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
			return commands.SetCommand(splitMsg)
		}
	}
	if command == "get" {
		if len(splitMsg) >= 4 {
			return commands.GetCommand(splitMsg)
		}
	}
	if command == "info" {
		if len(splitMsg) >= 4 {
			return commands.InfoCommand(splitMsg)
		}
	}

	return "-ERROR"

}

func main() {
	config.SetConfig()
	redisSever := NewSever()
	log.Fatal(redisSever.Start())
}
