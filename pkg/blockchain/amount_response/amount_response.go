package amount_response

import "encoding/json"

// Bilancio di un account, risposta in json
type AmountResponse struct {
	Amount float32 `json:"amount"`
}

// Json di AmountResponse
func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: ar.Amount,
	})
}
