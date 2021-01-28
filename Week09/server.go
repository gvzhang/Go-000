package main

import (
	"fmt"
	"log"
	"net"
)

type Handler interface {
	OnConnect(c net.Conn)
	OnMessage(c net.Conn, bytes []byte)
	OnClose(c net.Conn, err error)
}

func NewTcpServer(addr *net.TCPAddr, handler Handler) *TcpServer {
	ts := new(TcpServer)
	ts.addr = addr
	ts.handler = handler
	return ts
}

type TcpServer struct {
	handler  Handler
	addr     *net.TCPAddr
	listener *net.TCPListener
}

func (ts *TcpServer) Run() error {
	var err error
	ts.listener, err = net.ListenTCP("tcp", ts.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			log.Printf("accept fail, err: %v\n", err)
			continue
		}
		go func() {
			var consumerErr error
			defer func() {
				if r := recover(); r != nil {
					consumerErr = fmt.Errorf("panic %v", r)
				}
				ts.handler.OnClose(conn, consumerErr)
				err := conn.Close()
				if err != nil {
					log.Printf("conn.Close, err: %v\n", err)
				}
			}()
			consumerErr = ts.consumer(conn)
		}()
	}
}

func (ts *TcpServer) consumer(conn net.Conn) error {
	ts.handler.OnConnect(conn)

	err, data := ts.read(conn)
	if err != nil {
		return err
	}

	ts.handler.OnMessage(conn, data)
	return nil
}

func (ts *TcpServer) read(conn net.Conn) (error, []byte) {
	var readData []byte
	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			return err, readData
		}
		readData = append(readData, buf[:n]...)
	}
}

func (ts *TcpServer) Stop() error {
	var err error
	if ts.listener != nil {
		err = ts.listener.Close()
	}
	return err
}
