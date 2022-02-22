package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain/block"
	blockchain_transaction "github.com/iltommi1995/blockchain-go/pkg/blockchain/transaction"
	"github.com/iltommi1995/blockchain-go/pkg/utils"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "COINBASE TRANSACTION"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20
)

// Struct della blockchain
type Blockchain struct {
	transactionPool   []*blockchain_transaction.Transaction
	chain             []*block.Block
	blockchainAddress string
	port              uint16
	mux               sync.Mutex
}

// Funzione per creare una nuova Blockchain
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	// Crea la blockchain passando un blocco vuoto
	b := &block.Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

func (bc *Blockchain) TransactionPool() []*blockchain_transaction.Transaction {
	return bc.transactionPool
}

// Metodo per restituire in json la block
func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*block.Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

// Metodo di Blockchain, utilizzato per cerare un nuovo blocco, ritorna un puntatore a Block
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *block.Block {
	// Viene passato il transaction pool per le transazioni
	b := block.NewBlock(nonce, previousHash, bc.transactionPool)
	// Si appende il blocco alla catena di blocchi
	bc.chain = append(bc.chain, b)
	// Si svuota il transaction pool
	bc.transactionPool = []*blockchain_transaction.Transaction{}
	return b
}

// Metodo per ritornare l'ultimo blocco della blockchain
func (bc *Blockchain) LastBlock() *block.Block {
	return bc.chain[len(bc.chain)-1]
}

// Metodo per stampare i dati della blockchain
func (bc *Blockchain) Print() {
	fmt.Printf("\n%s BLOCKCHAIN WITH %x BLOCKS %s\n\n", strings.Repeat("*", 25), len(bc.chain), strings.Repeat("*", 25))
	for i, block := range bc.chain {
		fmt.Printf("\n%s Block %d %s \n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		fmt.Print("||\n")
		block.Print()
		fmt.Print("||\n")
		fmt.Printf("%s\n", strings.Repeat("=", 57))
	}
}

// Metodo della blockchain per creare una transazione
// ritorna un bool per verificare che AddTransaction sia andato a ubon fine
func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	// TODO
	// Sync
	return isTransacted
}

// Metodo per aggiungere una transazione al transactionPool
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := blockchain_transaction.NewTransaction(sender, recipient, value)

	// Se il sender è il miner, non va confermata la transazione
	if sender == MINING_SENDER {
		bc.transactionPool = append([]*blockchain_transaction.Transaction{t}, bc.transactionPool...)
		return true
	}

	// Se la firma della transazione viene verificata
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// Controllo che il sender abbia i soldi che invia
		/*
			if bc.CalculateTotalAmount(sender) < value {
				// In caso negativo do errore
				log.Println("ERROR: transaction rejected because sender doasn't have enough balance in wallet")
				return false
			}
		*/
		// Aggiungi la transazione al transaction Pool
		bc.transactionPool = append(bc.transactionPool, t)
		return true

		// In caso contrario errore
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false
}

// Metodo per aggiungere una transazione al transactionPool
func (bc *Blockchain) AddCoinbaseTransaction(sender string, recipient string, value float32) {
	t := blockchain_transaction.NewTransaction(sender, recipient, value)
	bc.transactionPool = append([]*blockchain_transaction.Transaction{t}, bc.transactionPool...)
}

// Metodo per copiatre il transaction pool
func (bc *Blockchain) CopyTransactionPool() []*blockchain_transaction.Transaction {
	// Creo array di transazioni vuoto
	transactions := make([]*blockchain_transaction.Transaction, 0)
	// Per ogni transazione nel transactionPool della blockchain
	for _, t := range bc.transactionPool {
		// Aggiungo all'array transactions la transazione
		transactions = append(transactions,
			blockchain_transaction.NewTransaction(
				t.SenderBlockchainAddress,
				t.RecipientBlockchainAddress,
				t.Value))
	}
	// Ritorno l'array popolato
	return transactions
}

// Metodo di *Blockchain per
// Prende come parametri i dati di un blocco (nonce, previousHash, transactions[]) più la difficoltà
// Ritorna true o false
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*blockchain_transaction.Transaction, difficulty int) bool {
	// In base alla difficoltà viene scelto il numero di zeri che deve avere l'hash del blocco'
	zeros := strings.Repeat("0", difficulty)
	// Blocco da indovinare
	guessBlock := block.Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
	// Hash del blocco
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	if guessHashStr[:difficulty] == zeros {

		fmt.Println(guessHashStr)
	}
	// Si controlla se gli n caratteri (dove n è la difficoltà) all'inizio dell'hash sono uguale agli zeros definiti prima
	return guessHashStr[:difficulty] == zeros
}

// Metodo per *Blockchain, ritorna un int, il nonce
func (bc *Blockchain) ProofOfWork() int {
	// Si crea la copia delle transazioni del transaction pool
	transactions := bc.CopyTransactionPool()
	// Si recupera l'hash del blocco precedente
	previousHash := bc.LastBlock().Hash()
	fmt.Println(previousHash)
	// Si parte da nonce = 0
	nonce := 0
	// Si calcola l'hash del nuovo blocco richiamando il metodo ValidProof, se non ritorna
	// true si aumenta di 1 il nonce e si riprova, fin quando il target non viene raggiunto
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

// Metodo di Blockchain per il mining
func (bc *Blockchain) Mining() bool {
	// Mutex permette di bloccare una parte di codice per far sì che
	// venga eseguita da una sola goroutine alla volta
	bc.mux.Lock()
	// Con defer, sì fa si che l'esecuzione di uno statement sia
	// differita fino al return della funzione
	// In sostanza rendo sincrono ogni cosa che avviene
	// dentro alla funzione di mining
	defer bc.mux.Unlock()

	// Se il transaction pool è vuoto, non mino il blocco
	if len(bc.transactionPool) == 0 {
		fmt.Println("**** Mining skipped because there are no transactions ****")
		return false
	}

	// Creo transazione coinbasem passando i dati
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	// Creo il nonce
	nonce := bc.ProofOfWork()
	// Creo il nuovo blocco
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
}

// Metodo per calcolare il bilancio di un account
// prende in input l'indirizzo dell'account di cui bisogna calcolare il bilancio
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	// Si setta inizialmente a zero
	var totalAmount float32 = 0.0
	// Per ogni blocco della blockchain
	for _, b := range bc.chain {
		// Per ogni transazione
		for _, t := range b.Transactions {
			// Prendo valore transazione
			value := t.Value
			// Se l'indirizzo di cui devo calcolare il bilancio è uguale a quello del recipient
			// sommo value al suo bilancio
			if blockchainAddress == t.RecipientBlockchainAddress {
				totalAmount += value
			}
			// Se invece è uguale al sender, sottraggo value al suo bilancio
			if blockchainAddress == t.SenderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

// Metodo per verificare la signature di una transazione
// Prende 3 parametri:
// 1- Public Key del sender della transazione
// 2- Signature generata dal sender attraverso la sua chiave privata
// 3- La transazione firmata
// Ritorna un bool
func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature,
	t *blockchain_transaction.Transaction) bool {
	// Faccio json della transazione
	m, _ := json.Marshal(t)
	// Calcolo l'hash del json della transazione
	h := sha256.Sum256([]byte(m))
	// Uso funzione della libreria ecdsa
	// Devo passare:
	// - chiave pubblica sender
	// - hash transazione
	// - R (da signature)
	// - S (da signature)
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}
