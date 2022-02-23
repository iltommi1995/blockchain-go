package transaction

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Transazione, contiene solo address del sender, del recipient e il valore inviato
type Transaction struct {
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

// Funzione per creare nuova transazione, tipo di ritorno puntatore a Transaction
//		- * => è un puntatore in Go, in questo caso "*Transaction" + un puntatore a
//			   una Transaction
//		- & => operatore usato per trovare l'indirizzo della variabile, ritorna un
//			   puntatore "*Transaction"
// Passo i parametri e creo nuova transazione
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

// Questo è un metodo, perché viene specificato un receiver (t *Transaction).
// In particolare, il receiver in questo caso è un pointer a Transaction.
// Il pointer receiver serve a creare metodi che possano modificare il valore della
// variabile a cui punta il receiver
func (t *Transaction) Print() {
	fmt.Printf("|| %s\n", strings.Repeat("-", 40))
	fmt.Printf("||  sender_blockchain_address    %s\n", t.SenderBlockchainAddress)
	fmt.Printf("||  recipient_blockchain_address    %s\n", t.RecipientBlockchainAddress)
	fmt.Printf("||  value    %.1f\n", t.Value)
}

// Anche in questo caso si tratta di un metodo, serve a formattare il json
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.SenderBlockchainAddress,
		Recipient: t.RecipientBlockchainAddress,
		Value:     t.Value,
	})
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	v := &struct {
		Sender    *string  `json:"sender_blockchain_address"`
		Recipient *string  `json:"recipient_blockchain_address"`
		Value     *float32 `json:"value"`
	}{
		Sender:    &t.SenderBlockchainAddress,
		Recipient: &t.RecipientBlockchainAddress,
		Value:     &t.Value,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
