package server

import (
	"blockchain/block"
	"blockchain/types"
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// HandleConn handles the TCP connection
func HandleConn(conn net.Conn) {
	defer conn.Close()

	serverChan := make(chan []types.Block)

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
			newBlock, err := block.Generate(block.Blockchain[len(block.Blockchain)-1], nin)
			if err != nil {
				log.Println(err)
				continue
			}
			if block.IsValid(newBlock, block.Blockchain[len(block.Blockchain)-1]) {
				newBlockchain := append(block.Blockchain, newBlock)
				block.ReplaceChain(newBlockchain)
			}

			serverChan <- block.Blockchain
			io.WriteString(conn, "\nEnter a new National Identification Number:")
		}
	}()

	// sleep some time to stimulate a real world syncing step
	// after 30s all clients get a copy of the full Blockchain
	// whether they created a Block or not
	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(block.Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range serverChan {
		spew.Dump(block.Blockchain)
	}
}
