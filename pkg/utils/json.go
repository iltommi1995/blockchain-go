package utils

import "encoding/json"

// Funzione per restituire lo stato di una request in json
func JsonStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}
