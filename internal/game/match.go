package game

// Match represents a game between two connected players
type Match struct {
	Players []*Player
	Status  string //WAITING, STARTED, FINISHED or CANCELLED
}

func NewMatch() *Match {
	m := Match{Players: []*Player{}, Status: "WAITING"}
	return &m
}

func (m *Match) AddPlayer(player *Player) {
	m.Players = append(m.Players, player)
}

func (m *Match) RemovePlayer(player *Player) {
	for i, p := range m.Players {
		if p == player {
			m.Players = append(m.Players[:i], m.Players[i+1:]...)
			break
		}
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
}

func (m *Match) End() {
	m.Status = "FINISHED"
}
