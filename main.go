package main

import (
	"blockchain/types"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// server handles incoming concurrent Blocks
var server chan []types.Block

// Blockchain represents a Blockchain
var Blockchain []types.Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	server = make(chan []types.Block)

	// create Genesis Block
	t := time.Now()
	genesisBlock := types.Block{
		Index:     0,
		Timestamp: t.String(),
		NIN:       0,
		Hash:      "",
		PrevHash:  "",
	}
	Blockchain = append(Blockchain, genesisBlock)
	spew.Dump(genesisBlock)

	// create TCP server (port 9000 - from env)
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a new National Identification Number:")

	scanner := bufio.NewScanner(conn)

	// take in NIN from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			nin, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], nin)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			server <- Blockchain
			io.WriteString(conn, "\nEnter a new National Identification Number:")
		}
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range server {
		spew.Dump(Blockchain)
	}
}

func calculateHash(block types.Block) string {
	record := string(block.Index) + block.Timestamp + string(block.NIN) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock types.Block, NIN int) (types.Block, error) {
	var newBlock types.Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.NIN = NIN
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock types.Block) bool {
	switch {
	case oldBlock.Index+1 != newBlock.Index:
		return false
	case oldBlock.Hash != newBlock.PrevHash:
		return false
	case calculateHash(newBlock) != newBlock.Hash:
		return false
	default:
		return true
	}
}

// The longer chain will have the latest blocks
func replaceChain(newBlocks []types.Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
