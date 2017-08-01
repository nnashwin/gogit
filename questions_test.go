package main

import (
	"gopkg.in/AlecAivazis/survey.v1"
	"testing"
)

func TestGetUserInput(*testing.T) {
	answers := Answer{"tyler", "cool"}
	GetUserInput([]*survey.Question{})

}
