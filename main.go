package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// ConnPort Default listen port
const ConnPort = "20000"
const header = "Hello from TestBoard Suite"
const version = "0.1.0"
const hwType = "TestBoard"

func sendReply(conn net.Conn, code int, msg string) {
	c := strconv.Itoa(code)
	s := ""
	if code < 0 {
		s += c + " Error in executing command\r\n"
	} else {
		s += c + " Command ok\r\n"
	}

	// Add message terminator, to allow the slsc class
	// to detect the end of message.
	if msg != "" {
		s += msg + "\r\n"
	}
	s += "\r\n"

	conn.Write([]byte(s))
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	clientAddr := conn.RemoteAddr()

	fmt.Println("Receive Connection from:" + clientAddr.String())
	conn.Write([]byte(header + " " + version + ".\n"))
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	for {
		// Remove prompt because slsc class, could give parse error
		//conn.Write([]byte(">> "))

		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		//fmt.Printf("Len[%v]: %s", reqLen, buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		s, code := ParseCmd(buf, reqLen)
		sendReply(conn, code, s)
		if s == "quit" {
			break
		}
	}
	fmt.Println("Close connection form:" + clientAddr.String())
	conn.Close()
}

func heartbeatLoop(beat net.Conn) {

	ip := beat.LocalAddr()
	ipAddr := strings.Split(ip.String(), ":")[0]

	macStr := ""
	ifs, err := net.Interfaces()
	if err != nil {
		macStr = ""
	} else {
		macStr = fmt.Sprintf("%s", ifs[1].HardwareAddr)
	}

	fmt.Println("Start server on: ", ipAddr, macStr)
	heartbeatStr := fmt.Sprintf("%s %s %s %s 0 2", ipAddr, macStr, hwType, version)

	for {
		_, err := beat.Write([]byte(heartbeatStr))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {

	// Init commands table
	InitCmd()

	ServerAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:23001")
	if err != nil {
		fmt.Printf("ResolveUDPAddr error %v", err)
	}
	beat, err := net.DialUDP("udp", nil, ServerAddr)
	if err != nil {
		fmt.Printf("DialUDP error %v", err)
		panic(err)
	}
	defer beat.Close()
	go heartbeatLoop(beat)

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
