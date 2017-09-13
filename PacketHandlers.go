package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	//NMServer "./NMServer"
)

type PacketHandler func(net.Conn, *PacketReader)

var PacketHandlers map[uint16]PacketHandler

func GetHandler(header uint16) PacketHandler {
	handler, ok := PacketHandlers[header]
	if !ok {
		handler = HandleNothing
	}

	return handler
}

func InitializePacketHandlers() {
	PacketHandlers = make(map[uint16]PacketHandler)
	PacketHandlers[51] = HandleLogin
	PacketHandlers[45] = HandleLogin2
	PacketHandlers[53] = HandleLogin3
}

func HandleNothing(conn net.Conn, p *PacketReader) {
	fmt.Printf("!!!! Did NOT handle this packet. Opcode: %d", p.opcode)
	fmt.Println()
}

func mb2sb(inp string) string {
	var ret string = ""
	for i := 0; i < len(inp); i = i + 2 {
		ret = ret + string(inp[i])
	}
	return ret
}

func HandleLogin(conn net.Conn, p *PacketReader) {
	p.ReadInt32()
	username := mb2sb(p.ReadString(int(p.ReadInt16()), false))
	password := mb2sb(p.ReadString(int(p.ReadInt16()), false))

	fmt.Println("Username:", username)
	//fmt.Println("Password:", password)

	var token string = "TOKEN"

	if *pserver {
		token = username
	} else {
		resp, err := http.PostForm(*server,
			url.Values{"userID": {username}, "password": {password}})

		if err != nil {
			fmt.Println("ERROR:", err)
		} else {

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				fmt.Println("ERROR:", err)
			} else {
				fmt.Println("Body:", body)
				for _, element := range resp.Cookies() {
					if element.Name == "NPPv2" {
						token = element.Value
					} else {
						fmt.Println("Dunno about this element:", element.Name)
					}
				}
			}
		}
	}

	fmt.Println("Token:", token)

	op := NewPacketWriter(52)
	op.WriteInt32(8)
	op.WriteInt32(0)
	op.WriteInt16(0)

	op.WriteString(token, false)

	op.WriteInt32(0)
	op.WriteInt32(0)
	op.WriteInt32(87)
	op.WriteInt16(0)

	SendPacket(conn, op)

}

func HandleLogin2(conn net.Conn, p *PacketReader) {

	/*token := */
	p.ReadString(int(p.ReadInt16()), false)

	op := NewPacketWriter(46)
	op.WriteInt32(0)
	op.WriteInt16(0)

	op.WriteString("RANDOM STUFF", false)
	op.WriteString("RANDOM NAME", false)
	op.WriteInt32(1234) // Account ID
	op.WriteInt16(0x16)
	op.WriteInt32(1)
	op.WriteByte(0)
	op.WriteInt32(0) // Token part 1
	op.WriteInt32(0) // Token part 2
	op.WriteInt32(0) // Token part 3
	op.WriteInt32(0) // Token part 4

	SendPacket(conn, op)
	//fmt.Println("DEBUG:", op)
}

func HandleLogin3(conn net.Conn, p *PacketReader) {
	p.ReadInt16() // 2 ??
	/*token := */ p.ReadString(int(p.ReadInt16()), false)
	p.ReadInt32() // Something

	op := NewPacketWriter(54)
	op.WriteInt32(2)
	op.WriteInt32(0)
	op.WriteInt16(0)

	SendPacket(conn, op)

}
