package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
)

var qs = []*survey.Question{
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "Please enter your github username or email"},
		Validate: survey.Required,
	},
	{
		Name:   "password",
		Prompt: &survey.Password{Message: "Please enter your github password or api token"},
	},
}

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello friend!")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "addUser",
			Usage: "add a github user account",
			Action: func(c *cli.Context) error {
				answers := struct {
					Username string
					Password string
				}{}

				err := survey.Ask(qs, &answers)
				if err != nil {
					fmt.Printf("obtaining user account info encountered an error")
					fmt.Printf("%+v", err)
					return nil
				}

				fmt.Printf("username: %s; password: %s", answers.Username, answers.Password)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
