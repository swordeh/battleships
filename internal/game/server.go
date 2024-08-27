package game

import (
	"fmt"
	"log"
	"net"
	"os"
)

// Server represents the system
type Server struct {
	Connections     []*Connection
	Players         []*Player
	PlayerQueue     chan *Player
	ConnectionEvent chan *Connection
	Matches         []*Match
}

func NewServer() *Server {
	return &Server{
		Connections:     make([]*Connection, 0),
		PlayerQueue:     make(chan *Player, 2),
		ConnectionEvent: make(chan *Connection),
		Matches:         make([]*Match, 0),
	}
}

func (s *Server) HandleConnection(conn *net.Conn) {
	// Create a Connection object which stores a pointer reference to the net.conn made
	c := &Connection{
		Conn:   conn,
		Status: 1,
	}
	log.Println("New connection from", (*conn).RemoteAddr())
	// Create a Player
	p := NewPlayer(c)
	// Add the Player to the Player pool
	s.AddPlayer(p)
	// Send the Message of the day
	s.SendMotd(p)
	// Handle commands
	go s.HandleCommands(p)

}

// AddConnection appends a connection to the server's list of connections.
// The connection is added to the server's Connections slice.
// TODO: Add lock
func (s *Server) AddConnection(conn *Connection) {
	s.Connections = append(s.Connections, conn)
}

func (s *Server) AddPlayer(p *Player) {
	s.Players = append(s.Players, p)
}

// JoinQueue adds a player to the player queue by sending them to the PlayerQueue channel
func (s *Server) JoinQueue(p *Player) {
	s.PlayerQueue <- p
}

// MatchMaker continuously pairs players from the PlayerQueue channel and pairs them for a game.
// It retrieves two players from the channel, calls PairPlayers to set them as in-game and notify them,
// and then repeats the process.
func (s *Server) MatchMaker() {

	// This function runs as a goroutine. It should run forever, and should handle players wanting to play a game
	// and put them in a match with another player.
	// To do this, the code should search for any eligible matches for a player to join. If there is, add them to the
	// match and start.
	// If there is no match to join, it should start a new one.
	log.Println("MATCHMAKER: process started")

	for {
		// Get a list of currently open matches, which are in the state WAITING
		// In theory this should never be more than 1
		var waitingMatches []*Match
		for _, match := range s.Matches {
			if match.Status == "WAITING" {
				waitingMatches = append(waitingMatches, match)
			}
		}

		// If there are no waiting matches, create one
		if len(waitingMatches) == 0 {
			log.Println("MATCHMAKER: no waiting matches, creating one")
			match := NewMatch()
			s.Matches = append(s.Matches, match)
			waitingMatches = append(waitingMatches, match)
		}

		select {
		case p := <-s.PlayerQueue:
			match := waitingMatches[0]

			// make sure that player has not tried to join again
			if !p.InGame {
				match.AddPlayer(p)
				p.SetInGameStatus(true)
			}

			if len(match.Players) == 2 {
				match.Start()
			}
		case <-s.ConnectionEvent:
			log.Println("MATCHMAKER: player dc")
		}
	}
}

func (s *Server) HandleCommands(p *Player) {
	// Set up channels for either a command or connection status update (disconnect)
	cmdChan := make(chan string)
	dcChan := make(chan *Connection)

	// Start a goroutine to listen to a connection
	go p.Connection.Listen(cmdChan, dcChan)
	for {
		select {
		case cmd := <-cmdChan:
			switch cmd {
			case "quit":
				p.Connection.Close()
				//os.Exit(0)
			case "play":
				s.JoinQueue(p)
			default:
				log.Println("Unknown command:", cmd)
			}
		case <-dcChan:
			// Disconnection
			log.Println("Server handing disconnection")
			s.HandleDisconnection(p)
		}
	}
}

func (s *Server) SendMotd(p *Player) {
	// open public/motd.txt
	motd, err := os.ReadFile("public/motd.txt")
	if err != nil {
		fmt.Println("Failed to read MOTD:", err)
	}
	p.Connection.Send(string(motd))
}

// HandleDisconnection will run when a player disconnects from the server.
// If the player is currently in game, HandleDisconnection will first cancel their match.
// It will then remove the player from the server's list of players, then send the player's *Connection to
// the server's ConnectionEvent channel.
func (s *Server) HandleDisconnection(p *Player) {
	if p.InGame == true {
		// Search for a player in the match and cancel
		for _, match := range s.Matches {
			for _, player := range match.Players {
				if player == p {
					match.RemovePlayer(p)
				}
			}
		}
	}
	s.RemovePlayer(p)
	s.ConnectionEvent <- p.Connection
}

func (s *Server) RemovePlayer(p *Player) {
	for i, player := range s.Players {
		log.Println("player removed from server")
		if p == player {
			s.Players = append(s.Players[:i], s.Players[i+1:]...)
			break

		}
	}
}
