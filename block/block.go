package block

import (
	"blockchain/types"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Blockchain represents a Blockchain
var Blockchain []types.Block

// Generate creates a new Block
func Generate(oldBlock types.Block, NIN int) (types.Block, error) {
	var newBlock types.Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.NIN = NIN
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)

	return newBlock, nil
}

// CalculateHash generates a new hex encoded string
func CalculateHash(block types.Block) string {
	record := string(block.Index) + block.Timestamp + string(block.NIN) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// IsValid checks if a given Block is valid
func IsValid(newBlock, oldBlock types.Block) bool {
	switch {
	case oldBlock.Index+1 != newBlock.Index:
		return false
	case oldBlock.Hash != newBlock.PrevHash:
		return false
	case CalculateHash(newBlock) != newBlock.Hash:
		return false
	default:
		return true
	}
}

// ReplaceChain replaces the chain - the longer chain will have the latest blocks
func ReplaceChain(newBlocks []types.Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// CreateGenesis creates a genesis(initial) block
func CreateGenesis() types.Block {
	t := time.Now()
	return types.Block{
		Index:     0,
		Timestamp: t.String(),
		NIN:       0,
		Hash:      "",
		PrevHash:  "",
	}
}
