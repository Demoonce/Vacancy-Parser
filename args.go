package main

import (
	"flag"
)

var (
	Query = flag.String("query", "golang", "Specify work query to the services that are being parsed")
	Level = flag.String("level", "senior", "Specify a job level which you want (junior, middle, senior)")
	BaseURL = "https://career.habr.com"
	Skill int
)
