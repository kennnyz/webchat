package main

import (
	"fmt"
	"net/http"
)

func main() {
	routes := routes()

	fmt.Println("Server is running on port 8080")
	_ = http.ListenAndServe(":8080", routes)
}
