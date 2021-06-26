package main

import (
	"fmt"
	"strings"

	"github.com/cucumber/gherkin-go"
)

func main() {
	feature := `
	Feature: whatevs

	Scenario Outline: test1
	Given A <what>
	And B
	Then C

	Examples:
		| what |
		| foo  |
	`
	doc, err := gherkin.ParseGherkinDocument(strings.NewReader(feature))
	if err != nil {
		panic(err)
	}
	for _, p := range doc.Pickles() {
		for _, s := range p.Steps {
			fmt.Printf("s: %v\n", s.Text)
			fmt.Printf("s.Arguments: %v\n", s.Arguments)
		}
	}
}
