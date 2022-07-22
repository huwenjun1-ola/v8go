package v8go

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
)

type ClientConnection struct {
	Connection      *websocket.Conn
	SendChan        chan interface{}
	ReadChan        chan interface{}
	StopSignal      chan interface{}
	ReadStopSignal  chan interface{}
	WriteStopSignal chan interface{}
	ProcessFun      func(connection *ClientConnection, data []byte)
	PostClosed      func(connection *ClientConnection)
}

func NewClientConnection(conn *websocket.Conn, procss func(connection *ClientConnection, data []byte), postClosed func(connection *ClientConnection)) *ClientConnection {
	return &ClientConnection{
		Connection:      conn,
		SendChan:        make(chan interface{}, 1000),
		ReadChan:        make(chan interface{}, 1000),
		StopSignal:      make(chan interface{}, 1),
		ReadStopSignal:  make(chan interface{}, 1),
		WriteStopSignal: make(chan interface{}, 1),
		ProcessFun:      procss,
		PostClosed:      postClosed,
	}
}
func (c *ClientConnection) Run() {
	go c.writeLoop()
	go c.readLoop()
	<-c.StopSignal
	c.ReadStopSignal <- 1
	c.WriteStopSignal <- 1
	_ = c.Connection.Close()
	if c.PostClosed != nil {
		c.PostClosed(c)
	}
}

func (c *ClientConnection) writeLoop() {
	for {
		select {
		case <-c.WriteStopSignal:
			return
		case mes := <-c.SendChan:
			err := websocket.Message.Send(c.Connection, mes)
			if err != nil && err != io.EOF {
				c.StopSignal <- 1
				return
			}
		}
	}
}

func (c *ClientConnection) readLoop() {
	for {
		select {
		case <-c.ReadStopSignal:
			return
		default:
			data2 := make([]byte, 4096)
			err := websocket.Message.Receive(c.Connection, &data2)
			if err != nil {
				c.StopSignal <- 1
				fmt.Printf("read error %v \n", err)
				return
			}

			c.ProcessFun(c, data2)

		}
	}
}
