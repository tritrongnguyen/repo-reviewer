package main

import (
	"log"
	"net/http"

	"github.com/tritrongnguyen/repo-reviewer.git/internal/router"
)

func main() {
	r := router.New()

	log.Println("Server is running at http://localhost:8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
