package main

import (
	"github.com/swordeh/battleships/internal/game"
	"log"
	"net"
	"time"
)

func main() {
	// Create a new Game Server
	gs := game.NewServer()

	// Run a matchmaker goroutine
	go gs.MatchMaker()

	// Listen for TCP connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	log.Println("Listening on", listener.Addr().(*net.TCPAddr).Port)

	for {
		// accept the connection made to the server
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Because net.Conn is an interface, use type assertion to use the TCPConn specific functions
		// to manage keepalive. net.Listen it set up work using TCP, but we check the concrete type
		// that implements net.Conn anyway
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(1 * time.Second)
		}
		// Pass the handling of the connection over to a function
		go gs.HandleConnection(&conn)
	}
}
