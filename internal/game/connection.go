package game

import (
	"bufio"
	"net"
	"strings"
)

// Connection represents a net.Conn connection between the gameserver and a client.
type Connection struct {
	Conn   *net.Conn
	Status int
}

func (c *Connection) Listen(ch chan string, disconnectChannel chan *Connection) {
	connection := *c.Conn
	reader := bufio.NewReader(connection)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			connection.Close()
			disconnectChannel <- c
			return
		}
		message = strings.TrimSpace(message)
		ch <- message
	}
}

func (c *Connection) SendBytes(msg []byte) {
	connection := *c.Conn
	connection.Write(msg)
}
func (c *Connection) Send(msg string) {
	connection := *c.Conn
	connection.Write([]byte(msg))
}

func (c *Connection) Close() {
	connection := *c.Conn
	connection.Close()
}
