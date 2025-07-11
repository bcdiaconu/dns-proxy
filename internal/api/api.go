package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type SetTxtRequest struct {
	Domain string `json:"domain"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type TxtRecordSetter interface {
	CreateTxtRecord(domain, key, value string) error
}

func SetTxtHandler(apiKey string, setter TxtRecordSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expected := "Bearer " + apiKey
		if authHeader != expected {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req SetTxtRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Domain == "" || req.Key == "" || req.Value == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = setter.CreateTxtRecord(req.Domain, req.Key, req.Value)
		if err != nil {
			log.Println("cPanel error:", err)
			http.Error(w, "Failed to set TXT record", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TXT record set"))
	}
}
