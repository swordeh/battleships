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
	c := Connection{
		Conn:   conn,
		Status: 1,
	}

	log.Println("New connection from", (*conn).RemoteAddr())

	// Create a Player
	p := NewPlayer(&c)

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

// RemovePlayer removes a player from the server's list of players.
// The player is removed from the server's Players slice.
func (s *Server) RemovePlayerByConnection(c *Connection) {
	for i, player := range s.Players {
		if player.Connection.Conn == c.Conn {
			player.Connection.Status = 0
			s.Players = append(s.Players[:i], s.Players[i+1:]...)
			break
		}
	}
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
			log.Println("MATCHMAKER player already in game, ignore")
			match := waitingMatches[0]

			// make sure that player has not tried to join again
			if !p.InGame {
				match.AddPlayer(p)
				p.InGame = true
			}

			if len(match.Players) == 2 {
				match.Start()
			}
		case c := <-s.ConnectionEvent:
			// This event will trigger when a connection is disconnected, but it's not the concern of the matchmaker
			// to handle anything other than the match making process. If they are in a match, this code will not
			// remove players from those matches, only those that are in the WAITING state
			log.Println("MATCHMAKER: player dc")
			for _, match := range waitingMatches {
				for _, player := range match.Players {
					if player.Connection == c {
						log.Println("MATCHMAKER: dc'd player was in a match and has been removed ")
						match.RemovePlayer(player)
					}
				}
			}
		}
	}
}

//func (s *Server) PairPlayers(p1, p2 *Player) error {
//
//	log.Println("pairing players")
//	p1.InGame = true
//	p2.InGame = true
//	m := NewMatch(p1, p2)
//	s.Matches = append(s.Matches, m)
//	s.StartGame(m)
//	return nil
//}

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
		case c := <-dcChan:
			// Disconnection
			log.Println("Server handing disconnection")
			s.RemovePlayerByConnection(c)
			s.ConnectionEvent <- c
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

func (s *Server) CheckConnection(m *Match) {
	for _, player := range m.Players {
		if player.Connection.Status == 0 {
			// Player disconnected
			s.RemovePlayerByConnection(player.Connection)
			m.Cancel()
		}
	}
}

func (s *Server) GetPlayerByConnection(c *Connection) *Player {
	for _, player := range s.Players {
		if player.Connection == c {
			return player
		}
	}
	return nil
}
