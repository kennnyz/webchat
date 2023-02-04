package main

import (
	"fmt"
	"github.com/kennnyz/webchat/handlers"
	"log"
	"net/http"
)

func main() {
	routes := routes()

	log.Println("Starting chanel listener")

	go handlers.ListenToWsChannel()

	fmt.Println("Server is running on port 8080")
	_ = http.ListenAndServe(":8080", routes)
}
