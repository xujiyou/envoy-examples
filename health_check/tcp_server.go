package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	log.Println("Listening on 127.0.0.1:1234")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}

		log.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

		go doServerStuff(conn)
	}
}

func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			log.Println("END reading", err.Error())
			return //终止程序
		}
		log.Printf("Received data: %v", string(buf[:len]))

		if string(buf) == "A" {
			_, _ = conn.Write([]byte("B"))
		} else {
			_, _ = conn.Write(buf)
		}
	}
}
