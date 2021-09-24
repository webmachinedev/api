package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/webmachinedev/api/pkg/serve"
)

func main() {
	http.HandleFunc("/", serve.Dir("data"))
	
	port := os.Getenv("PORT")
	fmt.Println("listening on "+port)
	http.ListenAndServe(":"+port, nil)
}
