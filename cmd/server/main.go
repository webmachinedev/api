package main

import (
	"net/http"
	"os"

	"github.com/webmachinedev/api/pkg/graphql"
)

func main() {
	http.HandleFunc("/", graphql.Handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
