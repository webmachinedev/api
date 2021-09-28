package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/webmachinedev/go-clients/github"
	"github.com/webmachinedev/models"
)

var rootIndex = []string{"types", "functions"}
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			HandleRead(w, r)
		case "POST":
			HandleWrite(w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	port := os.Getenv("PORT")
	fmt.Println("listening on "+port)
	http.ListenAndServe(":"+port, nil)
}

type ListItem struct {
	Name string `json:"name"`
	Doc string `json:"doc"`
	URL string `json:"url"`
}

func HandleRead(w http.ResponseWriter, r *http.Request) {
	resource := getResource(r)
	id := getID(r)

	switch resource {
	case "":
		WriteRootIndexResponse(w)
	case "types":
		if id == "" {
			json.NewEncoder(w).Encode(types)
		} else {
			json.NewEncoder(w).Encode(types[id])
		}
	case "functions":
		if id == "" {
			json.NewEncoder(w).Encode(functions)
		} else {
			json.NewEncoder(w).Encode(functions[id])
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func WriteRootIndexResponse(w http.ResponseWriter) {
	var index []ListItem
	for _, resource := range rootIndex {
		index = append(index, ListItem{resource, "", "/"+resource})
	}
	json.NewEncoder(w).Encode(index)
}

func HandleWrite(w http.ResponseWriter, r *http.Request) {
	resource := getResource(r)
	id := getID(r)

	switch resource {
	case "types":
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			var t models.Type
			err = json.Unmarshal(bytes, &t)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				if id != string(t.ID) {
					w.WriteHeader(http.StatusBadRequest)
				} else {
					githubkey := r.Header.Get("githubkey")
					err = WriteType(t, githubkey)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}
			}
		}
	case "functions":
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			var f models.Function
			err = json.Unmarshal(bytes, &f)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				if id != string(f.ID) {
					w.WriteHeader(http.StatusBadRequest)
				} else {
					githubkey := r.Header.Get("githubkey")
					err = WriteFunction(f, githubkey)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}
			}
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func WriteType(t models.Type, githubkey string) error {
	path := "data/types/"+string(t.ID)
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	err = github.WriteFile(
		"webmachinedev",
		"api",
		"main",
		path,
		string(bytes),
		"Update type "+string(t.ID),
		githubkey,
	)
	if err != nil {
		return err
	}

	types[string(t.ID)] = t
	return nil
}

func WriteFunction(f models.Function, githubkey string) error {
	path := "data/functions/"+string(f.ID)
	bytes, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	err = github.WriteFile(
		"webmachinedev",
		"api",
		"main",
		path,
		string(bytes),
		"Update type "+string(f.ID),
		githubkey,
	)
	if err != nil {
		return err
	}

	functions[string(f.ID)] = f
	return nil
}

func getResource(r *http.Request) string {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		return ""
	} else {
		return path[1]
	}
}

func getID(r *http.Request) string {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 3 {
		return ""
	} else {
		return path[2]
	}
}
