package main

import (
	"fmt"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain/blockchain"
	"github.com/iltommi1995/blockchain-go/pkg/wallet/transaction"
	"github.com/iltommi1995/blockchain-go/pkg/wallet/wallet"
)

// Main file per eseguire la blockchain in console

/*
func main() {
	myBlockchainAddress := "my_blockchain_address"
	blockchain := blockchain.NewBlockchain(myBlockchainAddress)
	blockchain.Print()

	blockchain.AddTransaction("A", "B", 1.0)
	blockchain.Mining()
	blockchain.Print()

	blockchain.AddTransaction("C", "D", 2.0)
	blockchain.AddTransaction("D", "A", 1.0)
	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("A %.1f\n", blockchain.CalculateTotalAmount("A"))
	fmt.Printf("B %.1f\n", blockchain.CalculateTotalAmount("B"))
	fmt.Printf("C %.1f\n", blockchain.CalculateTotalAmount("C"))
	fmt.Printf("D %.1f\n", blockchain.CalculateTotalAmount("D"))
}
*/

/*
func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.BlockchainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "B", 1.0)
	fmt.Printf("signature %s \n", t.GenerateSignature())
}
*/

func main() {
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	walletC := wallet.NewWallet()

	blockchain := blockchain.NewBlockchain(walletM.BlockchainAddress(), 20)

	// Primo blocco

	// Wallet transaction
	t := transaction.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)

	//  Blockchain node side

	isAdded := blockchain.AddTransaction(walletA.BlockchainAddress(),
		walletB.BlockchainAddress(), 1.0,
		walletA.PublicKey(),
		t.GenerateSignature())
	fmt.Println("Added? ", isAdded)

	blockchain.Mining()

	// Secondo blocco

	t2 := transaction.NewTransaction(walletC.PrivateKey(), walletC.PublicKey(), walletC.BlockchainAddress(), walletA.BlockchainAddress(), 2.0)

	isAdded = blockchain.AddTransaction(walletC.BlockchainAddress(),
		walletA.BlockchainAddress(), 2.0,
		walletC.PublicKey(),
		t2.GenerateSignature())
	fmt.Println("Added? ", isAdded)

	blockchain.Mining()

	blockchain.Print()

	fmt.Printf("A %.1f\n", blockchain.CalculateTotalAmount(walletA.BlockchainAddress()))
	fmt.Printf("B %.1f\n", blockchain.CalculateTotalAmount(walletB.BlockchainAddress()))
	fmt.Printf("C %.1f\n", blockchain.CalculateTotalAmount(walletC.BlockchainAddress()))
	fmt.Printf("M %.1f\n", blockchain.CalculateTotalAmount(walletM.BlockchainAddress()))

}
