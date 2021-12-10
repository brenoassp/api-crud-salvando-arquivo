package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/brenoassp/api-crud-salvando-arquivo/domain"
	"github.com/brenoassp/api-crud-salvando-arquivo/domain/person"
)

func main() {
	personService, err := person.NewService("person.json")
	if err != nil {
		fmt.Printf("Error trying to creating personService: %s\n", err.Error())
	}

	http.HandleFunc("/person/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			path := strings.TrimPrefix(r.URL.Path, "/person/")
			if path == "" {
				// list all people
			} else {
				_, err := strconv.Atoi(path)
				if err != nil {
					fmt.Println("Invalid id given. person ID must be an integer")
					return
				}

			}
		}
		if r.Method == "POST" {
			var person domain.Person
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s\n", err.Error())
				http.Error(w, "Error trying to create person", http.StatusBadRequest)
			}

			err = personService.Create(person)
			if err != nil {
				fmt.Printf("Error trying to create person: %s\n", err.Error())
				http.Error(w, "Error trying to create person", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
