package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain_server"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	fmt.Println("ciao")
	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Service")
	flag.Parse()

	app := blockchain_server.NewBlockchainServer(uint16(*port))
	app.Run()
}
