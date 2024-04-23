package utils

import (
	"sync"
	"zk/models"
)

var dataStorage *DataStorage

type DataStorage struct {
	persons []models.Person
}

var dataStorageOnce sync.Once

func GetDataStorage() *DataStorage {
	dataStorageOnce.Do(func() {
		dataStorage = &DataStorage{
			persons: make([]models.Person, 0),
		}
	})
	return dataStorage
}

func (d *DataStorage) GetPersons() []models.Person {
	return d.persons
}

func (d *DataStorage) SetPersons(persons []models.Person) {
	d.persons = persons
}

func (d *DataStorage) AddPerson(person models.Person) {
	d.persons = append(d.persons, person)
}
