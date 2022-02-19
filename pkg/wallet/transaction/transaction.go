package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"github.com/iltommi1995/blockchain-go/pkg/utils"
)

// Transazione lato wallet
type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBloackchainAddress   string
	recipientBlockchainAddress string
	value                      float32
}

// Funzione per creare nuova transaction
func NewTransaction(senderPrivateKey *ecdsa.PrivateKey,
	senderPublicKey *ecdsa.PublicKey,
	senderBloackchainAddress string,
	recipientBlockchainAddress string,
	value float32) *Transaction {
	return &Transaction{senderPrivateKey: senderPrivateKey,
		senderPublicKey:            senderPublicKey,
		senderBloackchainAddress:   senderBloackchainAddress,
		recipientBlockchainAddress: recipientBlockchainAddress,
		value:                      value}
}

// Metodo per generare la signature
func (t *Transaction) GenerateSignature() *utils.Signature {
	// Vogliamo computare l'hash della transazione
	m, _ := json.Marshal(t)
	// Calcoliamo l'hash della transazione
	h := sha256.Sum256([]byte(m))
	// Generiamo la signature a partire dalla private key
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{r, s}
}

// Json della transazione
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBloackchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}
