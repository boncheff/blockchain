package types

// Block represents the Blockchain block
type Block struct {
	Difficulty int
	Hash       string
	Index      int
	Iterations string
	NIN        int
	PrevHash   string
	Timestamp  string
}
