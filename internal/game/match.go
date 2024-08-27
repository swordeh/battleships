package game

import "log"

// Match represents a game between two connected players
type Match struct {
	Players []*Player
	Status  string //WAITING, STARTED
}

func NewMatch() *Match {
	m := Match{Players: []*Player{}, Status: "WAITING"}
	return &m
}

func (m *Match) AddPlayer(player *Player) {
	m.Players = append(m.Players, player)
}

// RemovePlayer will remove a player from the list of Players for a match.
// If the match is not in the WAITING state (it has started) then the match will be cancelled.
func (m *Match) RemovePlayer(player *Player) {
	for i, p := range m.Players {
		if p == player {
			m.Players = append(m.Players[:i], m.Players[i+1:]...)
			break
		}
	}
	if m.Status != "WAITING" {
		log.Println("match not in waiting state, cancel")
		m.Cancel()
	}
}

func (m *Match) Start() {
	m.Status = "STARTED"
	for _, player := range m.Players {
		player.Board.PlaceShipsRandomly()
		player.Connection.SendBytes(player.Board.Draw())
	}
}

func (m *Match) Cancel() {
	m.Status = "CANCELLED"
	for _, player := range m.Players {
		player.Connection.Send("Match cancelled\n")
		player.SetInGameStatus(false)
	}
}

func (m *Match) End() {
	m.Status = "FINISHED"
}

func (m *Match) Opponent(p *Player) *Player {
	for i, player := range m.Players {
		if player != p {
			return m.Players[i]
		}
	}
	return nil
}
