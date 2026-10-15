package main

import (
	"log"
	"net/http"

	"github.com/Vlad777-bit/gateway/internal/api"
	"github.com/Vlad777-bit/ledger"
)

func main() {
	l := ledger.NewLedger()

	server := api.NewServer(l)

	log.Println("Gateway started on http://localhost:8080")

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
