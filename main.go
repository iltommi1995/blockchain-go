package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
)

/*
Per cercare i processi:
netstat -ano | findstr :8000

taskkill /PID 26728 /F
*/

func main() {

	WALLET_SERVER_PATH := "./cmd/wallet_server/main_wallet_server.go"
	BLOCKCHAIN_SERVER_PATH := "./cmd/blockchain_server/main_blockchain_server.go"
	CONSOLE_PATH := "./cmd/console/main.go"

	MSG := "You have to insert one of accepted values:\n- blockchainserver -> to start a blockchainserver\n- walletserver -> to start a wallet server\n- console -> to show the blockchain on the console"

	program := flag.String("program", "blockchain_server", "Program to execute")

	port := flag.Uint("port", 8000, "TCP Port Number for Server")
	// Parametro per la gateway che permette al wallet di comunicare con il blockchain server
	gateway := flag.String("gateway", "http://localhost:5000", "Blockchain Gateway")
	flag.Parse()

	programStr := string(*program)

	cmd := exec.Command("cmd.exe", "/C", "start", "cmd", "/k", "go", "run")

	switch programStr {
	case "blockchainserver":
		cmd.Args = append(cmd.Args, BLOCKCHAIN_SERVER_PATH, "-port", fmt.Sprintf("%d", *port))
	case "walletserver":
		cmd.Args = append(cmd.Args, WALLET_SERVER_PATH, "-port", fmt.Sprintf("%d", *port), "-gateway", string(*gateway))
	case "console":
		cmd.Args = append(cmd.Args, CONSOLE_PATH)
	default:
		fmt.Println(MSG)
		return
	}

	fmt.Println(cmd.String())
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started in new cmd window")

	//var out bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Start()
	//err := cmd.Run()
	//fmt.Println("Process Pid: ")
	//cmd.Process.Kill()

	//if err != nil {
	//	log.Fatal(err)
	// }

	//fmt.Printf("Prova: %q\n", out.String())
}
