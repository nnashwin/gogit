package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"io"
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

func createCredDirIfNotExist(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.FileMode(0522))
	}
	return
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

				homeDir, err := homedir.Dir()

				err = createCredDirIfNotExist(homeDir + "/.gogit")
				if err != nil {
					fmt.Printf("creating a directory encountered an error")
					return nil
				}

				h := sha1.New()
				io.WriteString(h, "Sha1")
				io.WriteString(h, "Sha2")

				fmt.Printf("% x", h.Sum(nil))

				fmt.Printf("username: %s; password: %s", answers.Username, answers.Password)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
