package main

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

func GetUserInput(qs []*survey.Question) (Answer, error) {
	answers := Answer{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return Answer{}, err
	}
	return answers, nil
}
