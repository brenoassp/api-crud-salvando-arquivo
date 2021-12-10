package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/brenoassp/api-crud-salvando-arquivo/domain"
)

type Service struct {
	dbFilePath string
	people     domain.People
}

func NewService(dbFilePath string) (Service, error) {
	_, err := os.Stat(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = createEmptyFile(dbFilePath)
			if err != nil {
				return Service{}, err
			}
			return Service{
				dbFilePath: dbFilePath,
				people:     domain.People{},
			}, nil
		} else {
			return Service{}, err
		}
	}

	jsonFile, err := os.Open(dbFilePath)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to open file that contains all people: %s", err.Error())
	}

	jsonFileContentByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to read people file: %s", err.Error())
	}

	var allPeople domain.People
	json.Unmarshal(jsonFileContentByte, &allPeople)

	return Service{
		dbFilePath: dbFilePath,
		people:     allPeople,
	}, nil
}

func (s *Service) addPerson(person domain.Person) error {
	s.people.People = append(s.people.People, person)
	allPeopleJSON, err := json.Marshal(s.people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as JSON: %s", err.Error())
	}
	return ioutil.WriteFile(s.dbFilePath, allPeopleJSON, 0755)
}

func (s *Service) Create(person domain.Person) error {
	if s.exists(person) {
		return fmt.Errorf("There is already a person with this ID registered")
	}

	err := s.addPerson(person)
	if err != nil {
		return fmt.Errorf("Error trying to add Person to file: %s", err.Error())
	}

	return nil
}

func (s Service) exists(person domain.Person) bool {
	for _, currentPerson := range s.people.People {
		if currentPerson.ID == person.ID {
			return true
		}
	}
	return false
}

func (s Service) List() domain.People {
	return s.people
}

func (s Service) GetByID(personID int) (domain.Person, error) {
	for _, currentPerson := range s.people.People {
		if currentPerson.ID == personID {
			return currentPerson, nil
		}
	}
	return domain.Person{}, fmt.Errorf("Person not found")
}

func createEmptyFile(dbFilePath string) error {
	var people domain.People = domain.People{
		People: []domain.Person{},
	}
	peopleJSON, err := json.Marshal(people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as JSON: %s", err.Error())
	}

	err = ioutil.WriteFile(dbFilePath, peopleJSON, 0755)
	if err != nil {
		return fmt.Errorf("Error trying to writing people file: %s", err.Error())
	}

	return nil
}
