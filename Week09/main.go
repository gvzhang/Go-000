package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errChan := make(chan error, 0)
	tcpAddr := new(net.TCPAddr)
	tcpAddr.IP = net.IPv4(127, 0, 0, 1)
	tcpAddr.Port = 50000
	server := NewTcpServer(tcpAddr)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic %v", r)
			}
		}()
		errChan <- server.Run()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-errChan:
		if err != nil {
			log.Println("server.Run error", err)
		}
	case q := <-quit:
		log.Printf("catch exit signal %s\n", q.String())
	}
	err := server.Stop()
	if err != nil {
		log.Println("server.Stop error", err)
	}
	log.Println("server exit...")
}

func NewTcpServer(addr *net.TCPAddr) *TcpServer {
	ts := new(TcpServer)
	ts.addr = addr
	return ts
}

type TcpServer struct {
	addr     *net.TCPAddr
	listener *net.TCPListener
}

func (ts *TcpServer) Run() error {
	var err error
	ts.listener, err = net.ListenTCP("tcp", ts.addr)
	if err != nil {
		return err
	}

	readData := make(chan string, 0)
	go func() {
		for {
			conn, err := ts.listener.Accept()
			if err != nil {
				fmt.Printf("accept fail, err: %v\n", err)
				continue
			}
			err, data := ts.read(conn)
			readData <- data
		}
	}()

	go func() {
		for {
			conn := <-connChan
			ts.read(conn)
		}
	}()
	return nil
}

func (ts *TcpServer) read(conn net.Conn) (error, string) {
	var readData string
	defer conn.Close()
	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			return err, readData
		}
		readData += string(buf[:n])
	}
}

func (ts *TcpServer) Stop() error {
	var err error
	if ts.listener != nil {
		err = ts.listener.Close()
	}
	return err
}
