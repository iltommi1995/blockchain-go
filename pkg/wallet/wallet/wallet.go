package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/iltommi1995/blockchain-go/pkg/utils"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	w := new(Wallet)

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey

	// Creare address
	// 2. SHA-256 della chiave pubblica (32 bytes)
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// 3. Hash function RIPEMD-160 sul risultato di SHA-256 (20 byres)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// 4. Aggiungere version byte davanti all'hash RIPEMD-160 (0x00 per mainnet)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	// 5 - Si fa SHA-256 del risultato ottenuto nel punto 4
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// 6. Si fa SHA-256 del risultato del punto 5
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// 7. Prendere i primi 4 bytes del risultato del punto 5
	// come checksum
	checksum := digest6[:4]

	// 8. Agiungere i 4 bytes di checksum del punto 7 alla fine
	// del output del punto 4 (), si arriva quindi a 25 byes
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], checksum[:])

	// 9. Convertire il risultato da byte string a base58
	address := base58.Encode(dc8)

	w.blockchainAddress = address
	return w
}

// è una specie di getter, alternativa a scrivere in
// maiuscolo la prima lettera delle proprietà di
// uno struct
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

// TRANSAZIONI

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBloackchainAddress   string
	recipientBlockchainAddress string
	value                      float32
}

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

func (t *Transaction) GenerateSignature() *utils.Signature {
	// Vogliamo computare l'hash della transazione
	m, _ := json.Marshal(t)
	// Calcoliamo l'hash della transazione
	h := sha256.Sum256([]byte(m))
	// Generiamo la signature a partire dalla private key
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{r, s}
}

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
