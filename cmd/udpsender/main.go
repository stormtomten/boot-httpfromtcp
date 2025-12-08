package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const localAddr = "localhost:42069"

func main() {
	addr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		log.Fatalf("Failed to open network on %s: %s\n", localAddr, err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Faild to establish a connection to %s: %s", localAddr, err.Error())
	}
	defer conn.Close()

	buf := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		str, err := buf.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read line: %s\n", err.Error())
			continue
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			fmt.Printf("Failed to send message %s: %s\n", str, err.Error())
		}

	}
}
