package serve

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/webmachinedev/models"
)

func Dir(path string) func(w http.ResponseWriter, r *http.Request) {
	schema := readSchema(path)
	fmt.Println(schema)
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func readSchema(path string) map[string]models.Type {
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range dir {
		fmt.Println(folder.Name())
		if folder.IsDir() {

		}
	}
	return nil
}
