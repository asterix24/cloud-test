package main

import (
	"fmt"
	"net"
	"os"
)

// ConnPort Default listen port
const ConnPort = "20000"
const header = "Hello from TestBoard Suite"
const version = "0.1.0"

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	clientAddr := conn.RemoteAddr()

	fmt.Println("Receive Connection from:" + clientAddr.String())
	conn.Write([]byte(header + " " + version + ".\n"))
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	for {
		// Send a response back to person contacting us.
		conn.Write([]byte(">> "))

		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		//fmt.Printf("Len[%v]: %s", reqLen, buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		s, err := ParseCmd(buf, reqLen)
		if err != nil {
			conn.Write([]byte(err.Error()))
		}

		if len(s) != 0 {
			if s == "quit" {
				conn.Write([]byte("Bye bye!\n"))
				break
			}
			conn.Write([]byte(s + "\n"))
		}
	}

	fmt.Println("Close connection form:" + clientAddr.String())
	conn.Close()
}

func main() {

	// Init commands table
	InitCmd()

	// Listen for incoming connections.
	l, err := net.Listen("tcp", "0.0.0.0"+":"+ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println(header + " " + version)
	fmt.Println("Listening on " + ConnPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}
