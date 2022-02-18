package main

import (
	"blockchain-go/src/wallet_server"
	"flag"
	"log"
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
