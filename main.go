package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	ls, err := net.Listen("tcp", ":8100")
	if err != nil {
		panic(err)
	}
	defer ls.Close()

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fullbuf := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}

		fullbuf = append(fullbuf, buf[:n]...)

		leftBuf, incomplete := decord(fullbuf)
		if incomplete {
			continue
		}

		fullbuf = leftBuf
	}
}

//解码数据包
func decord(pack []byte) ([]byte, bool) {
	//不完整包
	if len(pack) < 6 {
		return nil, true
	}
	packLen := int(binary.BigEndian.Uint32(pack[:4]))
	if packLen > len(pack) {
		return nil, true
	}

	//解析完整包
	fmt.Printf("packetLen:%v\n", packLen)

	headerLen := binary.BigEndian.Uint16(pack[4:6])
	fmt.Printf("headerLen:%v\n", headerLen)

	version := binary.BigEndian.Uint16(pack[6:8])
	fmt.Printf("version:%v\n", version)

	operation := binary.BigEndian.Uint32(pack[8:12])
	fmt.Printf("operation:%v\n", operation)

	sequence := binary.BigEndian.Uint32(pack[12:16])
	fmt.Printf("sequence:%v\n", sequence)

	body := string(pack[16:packLen])
	fmt.Printf("body:%v\n", body)

	//返回剩余数据
	return pack[packLen:], false
}
