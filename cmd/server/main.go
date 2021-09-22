package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/webmachinedev/api/pkg/graphql"
)

func main() {
	http.HandleFunc("/", graphql.Handler)
	port := os.Getenv("PORT")
	fmt.Println("listening on "+port)
	http.ListenAndServe(":"+port, nil)
}
