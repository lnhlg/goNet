package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8100")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("connect success")

	go func(conn net.Conn) {
		for i := 0; i < 2; i++ {
			conn.Write(encoder("test*****"))
		}
		fmt.Println("send over")
	}(conn)

	for {
		time.Sleep(1 * 1e9)
	}
}

func encoder(body string) []byte {
	headerLen := 16
	packLen := len(body) + headerLen

	ret := make([]byte, packLen)
	binary.BigEndian.PutUint32(ret[:4], uint32(packLen))
	binary.BigEndian.PutUint16(ret[4:6], uint16(headerLen))

	version := 5
	binary.BigEndian.PutUint16(ret[6:8], uint16(version))

	operation := 6
	binary.BigEndian.PutUint32(ret[8:12], uint32(operation))

	sequence := 7
	binary.BigEndian.PutUint32(ret[12:16], uint32(sequence))

	byteBody := []byte(body)
	copy(ret[16:], byteBody)

	return ret
}
