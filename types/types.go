package types

// Block represents the Blockchain block
type Block struct {
	Index     int
	Timestamp string
	NIN       int
	Hash      string
	PrevHash  string
}
