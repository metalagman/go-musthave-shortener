package main

import (
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	addr := "localhost:8080"
	shortener := app.NewMemoryShortenerService(addr)
	http.HandleFunc("/", app.RouterHandler(shortener))
	fmt.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
