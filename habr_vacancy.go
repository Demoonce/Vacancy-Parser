package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// career.habr.com api json response struct
type HabrVacanciesPage struct {
	List []struct {
		Href   string `json:"href"`
		Title  string `json:"title"`
		Salary struct {
			Formatted string `json:"formatted"`
		} `json:"salary"`
		Locations []struct {
			Title string `json:"title"`
		} `json:"locations"`

		Skills []struct {
			Title string `json:"title"`
		} `json:"skills"`
	} `json:"list"`
	Meta struct {
		TotalResults int `json:"totalResults"`
		TotalPages   int `json:"totalPages"`
	} `json:"meta"`
}

// career.habr.com skills api structure
type SkillsSuggestions struct {
	List []struct {
		Value int `json:"value"`
	} `json:"list"`
}

// Returns the career.habr value of the request specified by --query flag
func GetSkill() (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://career.habr.com/api/frontend/suggestions/skills?term=%s", *Query))
	if err != nil {
		log.Fatalln("Can't make a request to skills suggestions service", err)
	}
	suggestion := new(SkillsSuggestions)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Can't read data from response body", err)
	}
	err = json.Unmarshal(data, suggestion)
	if err != nil {
		log.Fatalln("Can't unmarshall json got by skills service", err)
	}
	if suggestion.List[0].Value == 0 {
		return 0, errors.New("Can't find this skill")
	}
	return suggestion.List[0].Value, nil
}

// Gets programmer level id for using in response
func GetLevelId() (int, error) {
	level := strings.ToLower(*Level)
	switch level {
	case "intern":
		return 1, nil
	case "junior":
		return 3, nil
	case "middle":
		return 4, nil
	case "senior":
		return 5, nil
	case "lead":
		return 6, nil
	default:
		return 0, errors.New("Invalid value was specified as a programmer level")
	}
}

// Laucnhes goroutines to get data from each page
func GetAllPages() []*Vacancy {
	data := make([]*Vacancy, 0, 100)
	page_chan := make(chan *Vacancy)

	first_page := GetPageData(Skill, 1)
	vacancy_list := new(HabrVacanciesPage)
	err := json.Unmarshal(first_page, vacancy_list)
	if err != nil {
		log.Printf("Error when unmarshalling body of page %d\n", 1)
	}
	for page := 0; page < vacancy_list.Meta.TotalPages; page++ {
		go GetPage(Skill, page, page_chan)
	}
	for a := 0; a < vacancy_list.Meta.TotalResults; a++ {
		data = append(data, <-page_chan)
	}
	return data
}

// Gets body data from response, then returns it
func GetPageData(skill, page int) []byte {
	level, err := GetLevelId()
	if err != nil {
		log.Fatalln(err)
	}
	base_request := "https://career.habr.com/api/frontend/vacancies?qid=%d&sort=relevance&type=all&skills[]=%d&page=%d"

	resp, err := http.Get(fmt.Sprintf(base_request, level, skill, page))
	if err != nil {
		log.Printf("Error when requesting page %d\n", page)
		return nil
	}
	resp_data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error when reading from body of page %d\n", page)
		return nil
	}
	return resp_data
}

// Parses vacancies from api, then pushes them to channel
func GetPage(skill, page int, result_chan chan<- *Vacancy) {
	resp_data := GetPageData(skill, page)
	if resp_data == nil {
		return
	}
	vacancy_list := new(HabrVacanciesPage)
	err := json.Unmarshal(resp_data, vacancy_list)
	if err != nil {
		log.Printf("Error when unmarshalling body of page %d\n", page)
	}
	for _, vacancy := range vacancy_list.List {
		location_list := make([]string, 0, len(vacancy.Locations))
		skills_list := make([]string, 0, len(vacancy.Skills))
		for _, a := range vacancy.Locations {
			location_list = append(location_list, a.Title)
		}
		for _, a := range vacancy.Skills {
			skills_list = append(skills_list, a.Title)
		}
		new_vac := NewVacancy(vacancy.Title, vacancy.Salary.Formatted, strings.Join(skills_list, ","), strings.Join(location_list, ","), vacancy.Href)
		result_chan <- new_vac
	}
}
