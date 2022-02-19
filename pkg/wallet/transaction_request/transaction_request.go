package transaction_request

// Transaction request, che si fa lato wallet
type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBloackchainAddress   *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      *string `json:"value"`
}

// Metodo per validare TransactionRequest
func (tr *TransactionRequest) Validate() bool {
	// Controllo se tutti i dati non sono null e restituisco un bool
	if tr.SenderPrivateKey == nil ||
		tr.SenderBloackchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
