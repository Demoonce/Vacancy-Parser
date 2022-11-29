package main

import (
	"encoding/csv"
	"log"
	"os"
)

// Unified structure of vacancy
type Vacancy struct {
	Title    string
	Salary   string // formatted so string
	Skills   string
	Location string
	Link     string
}

func NewVacancy(title, salary, skills, location, link string) *Vacancy {
	return &Vacancy{
		title,
		salary,
		skills,
		location,
		BaseURL + link,
	}
}

func (vac *Vacancy) GetData() []string {
	return []string{vac.Title, vac.Salary, vac.Skills, vac.Location, vac.Link}
}

// Writes all parsed data into .csv file
func WriteToFile(data []*Vacancy, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln("Can't create file to write results")
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, vac := range data {
		err := writer.Write(vac.GetData())
		if err != nil {
			log.Println(err)
		}
	}
}
