package game

// Player represents a user who has connected to the system
type Player struct {
	Board      Board
	InGame     bool
	Connection *Connection
}

func NewPlayer(conn *Connection) *Player {
	p := Player{
		Board:      NewEmptyBoard(),
		InGame:     false,
		Connection: conn,
	}
	return &p
}

func (p *Player) SetInGameStatus(status bool) {
	p.InGame = status
}
