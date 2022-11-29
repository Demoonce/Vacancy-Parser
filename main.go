package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		return
	}
	var err error

	Skill, err = GetSkill()
	if err != nil {
		log.Fatalln(err)
	}

	vacancies := GetAllPages()
	WriteToFile(vacancies, "vacs.csv")
}
