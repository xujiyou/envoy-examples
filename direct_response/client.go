package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:82")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = fmt.Fprintln(conn)
	result, _ := bufio.NewReader(conn).ReadString('\n')
	log.Println(result)
}
