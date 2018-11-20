package harmony

type MessageActivityType int

const (
	Join MessageActivityType = iota
	Spectate
	Listen
	JoinRequest
)

type MessageActivity struct {
	Type    MessageActivityType
	PartyID string
}

type MessageApplication struct {
	ID          string
	CoverImage  string
	Description string
	Icon        string
	Name        string
}
