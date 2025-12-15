package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("The Council Backend is starting...")
	
	// Basic health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := ":8080"
	fmt.Printf("Server listening on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
