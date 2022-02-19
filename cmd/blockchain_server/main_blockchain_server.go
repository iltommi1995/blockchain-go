package main

import (
	"flag"
	"log"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain_server"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

// Main per avviare il server della blockchain (un nodo)
// Si può passare il parametro con "- port <portnumber>"
func main() {
	// Flag serve a parsare i comandi da command line
	// https://pkg.go.dev/flag
	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Service")
	flag.Parse()

	// Creo il blockchain server
	app := blockchain_server.NewBlockchainServer(uint16(*port))
	// Starto il server
	app.Run()
}
