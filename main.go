package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type ConnectionRequest struct {
	Name string `json:"name"`
}

type MessageRequest struct {
	To        string `json:"to,omitempty"`
	Message   string `json:"message"`
	Broadcast bool   `json:"broadcast,omitempty"`
}

func readSock(conn net.Conn) {

	if conn == nil {
		panic("Connection is nil")
	}
	buf := make([]byte, 256)
	eof_count := 0
	for {
		for i := 0; i < 256; i++ {
			buf[i] = 0
		}

		readed_len, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				eof_count++
				time.Sleep(time.Second * 2)
				fmt.Println("EOF")
				if eof_count > 7 {

					fmt.Println("Timeout connection")
					break
				}
				continue
			}
			if strings.Index(err.Error(), "use of closed network connection") > 0 {

				fmt.Println("connection not exist or closed")
				continue
			}
			panic(err.Error())
		}
		if readed_len > 0 {
			fmt.Println(string(buf))
		}

	}
}

func readConsole(ch chan string) {
	for {
		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if len(line) > 250 {
			fmt.Println("Error: message is very lagre")
			continue
		}
		fmt.Print(">")
		out := line[:len(line)-1]

		ch <- out
	}
}

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	if conn == nil {
		panic("Connection is nil")

	}

	fmt.Println("Introduce yourself")

	name, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	request := ConnectionRequest{Name: name}

	enc := json.NewEncoder(conn)
	err := enc.Encode(request)

	if err != nil {
		fmt.Println("Write error:", err.Error())
	}

	ch := make(chan string)

	defer close(ch)

	go readConsole(ch)
	go readSock(conn)

	for {
		val, ok := <-ch
		if ok {
			out := MessageRequest{Message: val, Broadcast: true}

			enc := json.NewEncoder(conn)
			err := enc.Encode(out)
			if err != nil {
				fmt.Println("Write error:", err.Error())
				break
			}

		} else {
			time.Sleep(time.Second * 2)
		}

	}
	fmt.Println("Finished...")

	conn.Close()
}
