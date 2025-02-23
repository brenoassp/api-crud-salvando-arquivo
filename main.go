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
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(personService.List())
				if err != nil {
					http.Error(w, "Error trying to list people", http.StatusInternalServerError)
					return
				}
			} else {
				personID, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "Invalid id given. person ID must be an integer", http.StatusBadRequest)
					return
				}
				person, err := personService.GetByID(personID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(person)
				if err != nil {
					http.Error(w, "Error trying to get person", http.StatusInternalServerError)
					return
				}
			}
			return
		}
		if r.Method == "POST" {
			var person domain.Person
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s\n", err.Error())
				http.Error(w, "Error trying to create person", http.StatusBadRequest)
				return
			}
			if person.ID <= 0 {
				http.Error(w, "person ID should be a positive integer", http.StatusBadRequest)
				return
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
		if r.Method == "DELETE" {
			path := strings.TrimPrefix(r.URL.Path, "/person/")
			if path == "" {
				http.Error(w, "ID is required to delete a person", http.StatusBadRequest)
				return
			} else {
				personID, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "Invalid id given. person ID must be an integer", http.StatusBadRequest)
					return
				}
				err = personService.DeleteByID(personID)
				if err != nil {
					fmt.Printf("Error trying to delete person: %s\n", err.Error())
					http.Error(w, "Error trying to delete person", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}
			return
		}
		if r.Method == "PUT" {
			var person domain.Person
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s\n", err.Error())
				http.Error(w, "Error trying to update person", http.StatusBadRequest)
				return
			}
			if person.ID <= 0 {
				http.Error(w, "person ID should be a positive integer", http.StatusBadRequest)
				return
			}

			err = personService.Update(person)
			if err != nil {
				fmt.Printf("Error trying to update person: %s\n", err.Error())
				http.Error(w, "Error trying to update person", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	http.ListenAndServe(":80", nil)
}
