package main

import (
	"net"
	"os"
)

func main() {
	strEcho := "Halo\n"
	serAddr := "127.0.0.1:50000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(strEcho))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	println("write to server = ", strEcho)

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	println("reply from server=", string(reply))

	err = conn.Close()
	println("close err", err)
}
