package main

import (
	"flag"
	"log"

	"github.com/iltommi1995/blockchain-go/pkg/wallet_server"
)

func init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	port := flag.Uint("port", 8000, "TCP Port Number for Wallet Server")
	gateway := flag.String("gateway", "http://localhost:5000", "Blockchain Gateway")
	flag.Parse()

	app := wallet_server.NewWalletServer(uint16(*port), *gateway)

	app.Run()
}
