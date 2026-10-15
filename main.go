package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/GhandyP/42-sales-scripts-calls/internal/web"
)

func main() {
	server := web.NewServer(http.NewServeMux())
	log.Printf("Sales script app listening at http://%s", server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
