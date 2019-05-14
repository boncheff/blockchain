package main

import (
	"blockchain/block"
	"blockchain/server"
	"log"
	"net"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// create initial Block
	genesisBlock := block.CreateGenesis()
	block.Blockchain = append(block.Blockchain, genesisBlock)
	spew.Dump(genesisBlock)

	// create TCP server
	srvr, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer srvr.Close()

	for {
		conn, err := srvr.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.HandleConn(conn)
	}
}
