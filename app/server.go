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
}

func NewSever() *RedisServer {
	return &RedisServer{
		quitch: make(chan struct{}),
		msgch:  make(chan []byte),
	}
}

func (rs *RedisServer) Start() error {
	ln, err := net.Listen("tcp", config.Configs.ListenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()

	fmt.Println("redis server starting on port", config.Configs.ListenAddr)
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

		if !strings.Contains(string(msg), "sync_db") {
			response = response + "\r\n"
		}

		conn.Write([]byte(response))
	}
}

func (rs *RedisServer) parseMessage(message []byte) string {
	msg := strings.ToLower(strings.TrimSpace(string(message)))
	splitMsg := strings.Split(msg, "\r\n")

	//fmt.Printf("%#v\n", splitMsg)

	command := splitMsg[2]
	fmt.Println(command)

	if command == "sync_db" {
		return string(commands.SyncRdb())
	}
	if command == "psync" {
		return commands.PsyncCommand(splitMsg)
	}
	if command == "replconf" {
		return commands.ReplconfCommand(splitMsg)
	}
	if command == "ping" {
		return commands.PingCommand()
	}
	if command == "echo" {
		if len(splitMsg) >= 4 {
			return commands.EchoCommand(splitMsg)
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

func (rs *RedisServer) sendHandshake() error {
	bytes := make([]byte, 1024)

	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte("*1\r\n$4\r\nping\r\n"))
	if err != nil {
		return err
	}
	_, err = conn.Read(bytes)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n6380\r\n"))
	if err != nil {
		return err
	}
	_, err = conn.Read(bytes)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
	if err != nil {
		return err
	}
	_, err = conn.Read(bytes)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))
	if err != nil {
		return err
	}
	_, err = conn.Read(bytes)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte("*1\r\n$4\r\nsync_db\r\n"))
	if err != nil {
		return err
	}
	return nil

}

func main() {
	commands.InitCommands()
	config.SetConfig()
	redisSever := NewSever()

	if config.Configs.ClusterType == "slave" {
		err := redisSever.sendHandshake()
		if err != nil {
			panic("error")
		}
	}

	log.Fatal(redisSever.Start())
}
