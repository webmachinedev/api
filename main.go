package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/webmachinedev/go-clients/github"
	"github.com/webmachinedev/models"
)

var functions = make(map[string]models.Function)
var types = make(map[string]models.Type)

func init() {
	functionsDir, err := os.ReadDir("data/functions")
	if err != nil {
		panic(err)
	}
	for _, functionFile := range functionsDir {
		filename := functionFile.Name()
		path := "data/functions/"+filename
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		decoder := json.NewDecoder(file)
		var function models.Function
		decoder.Decode(&function)
		functions[string(function.ID)] = function
	}

	typesDir, err := os.ReadDir("data/types")
	if err != nil {
		panic(err)
	}
	for _, typeFile := range typesDir {
		filename := typeFile.Name()
		path := "data/types/"+filename
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		decoder := json.NewDecoder(file)
		var t models.Type
		decoder.Decode(&t)
		types[string(t.ID)] = t
	}
}

func main() {
	http.HandleFunc("/functions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(functions)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.HandleFunc("/functions/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			id := strings.TrimPrefix(r.URL.Path, "/functions/")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(functions[id])
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(types)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.HandleFunc("/types/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			id := strings.TrimPrefix(r.URL.Path, "/types/")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(types[id])
		case "POST":
			id := strings.TrimPrefix(r.URL.Path, "/types/")
			var t models.Type
    		err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				err = setType(id, t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(types[id])
				}
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			index := Index{
				{"functions", "", "/functions"},
				{"types", "", "/types"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(index)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	port := os.Getenv("PORT")
	fmt.Println("listening on "+port)
	http.ListenAndServe(":"+port, nil)
}

type Index []struct {
	Name string `json:"name"`
	Doc string `json:"doc"`
	URL string `json:"url"`
}

func setType(id string, t models.Type) error {
	owner := "webmachinedev"
	repo := "api"
	branch := "main"
	filename := "data/types/"+string(t.ID)+".json"
	bytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	file := string(bytes)
	commitmessage := "Update "+t.Name+" type"
	githubkey := os.Getenv("GITHUB_TOKEN")
	return github.WriteFile(owner, repo, branch, filename, file, commitmessage, githubkey)
}
