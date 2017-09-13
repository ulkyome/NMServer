package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

type connectionHandler func(net.Conn)

func PrintBytes(buffer []byte) {
	for i := 0; i < len(buffer); i++ {
		fmt.Printf("%02X ", buffer[i])
	}
	fmt.Println()
}

func StartServer(port int, handler connectionHandler) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		// handle error
		fmt.Println("Unable to bind. Error", err)
		os.Exit(1)
		return
	}

	fmt.Println("Listening on", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Unable to accept. Error", err)
			continue
		}
		go handler(conn)
	}
}

func ReadBytes(conn net.Conn, amount int) ([]byte, error) {
	//var read int
	var err error
	buffer := make([]byte, amount)

	if _, err = conn.Read(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func ReadPackets(conn net.Conn) {
	for StartReadPacket(conn) {
	}
	conn.Close()
}

func StartReadPacket(conn net.Conn) bool {
	packet, err := ReadPacket(conn)
	if err != nil {
		return false
	}

	GetHandler(uint16(packet.opcode))(conn, packet)

	return true
}

var pserver = flag.Bool("pserver", true, "a bool")
var server = flag.String("server", "http://127.0.0.1/api/login", "a string")
var port = flag.Int("port", 47611, "an int")

func main() {
	flag.Parse()

	InitializePacketHandlers()
	fmt.Println("Pserver =", *pserver)
	fmt.Println(*server)

	StartServer(*port, func(conn net.Conn) {
		fmt.Println("Got LOGIN connection! ", conn)
		go ReadPackets(conn)
	})

	fmt.Println("Done?")
}
