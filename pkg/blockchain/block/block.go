package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain/transaction"
)

// Struct dei singoli blocchi
type Block struct {
	Timestamp    int64
	Nonce        int
	PreviousHash [32]byte
	Transactions []*transaction.Transaction
}

// Funzione per creare un nuovo blocco, ha come parametri:
// 		- Nonce
//		- previousHash
//		- transactions (che è un array di puntatori a Transaction)
// Ha come tipo di ritorno un puntatore a Block
func NewBlock(timestamp int64, nonce int, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	b := new(Block)
	b.Timestamp = timestamp
	b.Nonce = nonce
	b.PreviousHash = previousHash
	b.Transactions = transactions
	return b
}

func FirstBlock(nonce int, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	b := new(Block)
	b.Timestamp = 1645635740778068200
	b.Nonce = nonce
	b.PreviousHash = previousHash
	b.Transactions = transactions
	return b
}

// Metodo per stampare i dati del blocco
func (b *Block) Print() {
	fmt.Printf("|| block_hash     %x\n", b.Hash())
	fmt.Printf("|| timestamp     %d\n", b.Timestamp)
	fmt.Printf("|| nonce     %d\n", b.Nonce)
	fmt.Printf("|| previous_hash     %x\n", b.PreviousHash)
	fmt.Printf("|| transactions:\n")
	for _, t := range b.Transactions {
		t.Print()
	}
}

// Metodo per creare l'hash del blocco
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

// Funzione per formattare il json
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64                      `json:"timestamp"`
		Nonce        int                        `json:"nonce"`
		PreviousHash string                     `json:"previous_hash"`
		Transactions []*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    b.Timestamp,
		Nonce:        b.Nonce,
		PreviousHash: fmt.Sprintf("%x", b.PreviousHash),
		Transactions: b.Transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	v := &struct {
		Timestamp    *int64                      `json:"timestamp"`
		Nonce        *int                        `json:"nonce"`
		PreviousHash *string                     `json:"previous_hash"`
		Transactions *[]*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    &b.Timestamp,
		Nonce:        &b.Nonce,
		PreviousHash: &previousHash,
		Transactions: &b.Transactions,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	ph, _ := hex.DecodeString(*v.PreviousHash)
	copy(b.PreviousHash[:], ph[:32])
	return nil
}
