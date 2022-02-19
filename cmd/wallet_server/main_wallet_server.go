package main

import (
	"flag"
	"log"

	"github.com/iltommi1995/blockchain-go/pkg/wallet_server"
)

func init() {
	log.SetPrefix("Wallet Server: ")
}

// Main del wallet server
func main() {
	// Parametro per il comando (-port <portnumber>)
	port := flag.Uint("port", 8000, "TCP Port Number for Wallet Server")
	// Parametro per la gateway che permette al wallet di comunicare con il blockchain server
	gateway := flag.String("gateway", "http://localhost:5000", "Blockchain Gateway")
	flag.Parse()

	// Creo il wallet server
	app := wallet_server.NewWalletServer(uint16(*port), *gateway)
	// Starto il wallet server
	log.Println("Wallet Server running")
	app.Run()
}
