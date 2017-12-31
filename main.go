package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
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

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func fileDoesExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func createCredDirIfNotExist(path string) (err error) {
	doesExist, err := fileDoesExist(path)
	checkErr(err)

	if doesExist == false {
		err = os.Mkdir(path, 0766)
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
			Usage: "add a new github user account",
			Action: func(c *cli.Context) error {
				answers := struct {
					Username string
					Password string
				}{}

				err := survey.Ask(qs, &answers)
				checkErr(err)

				homeDir, err := homedir.Dir()
				checkErr(err)

				dirString := homeDir + "/.gogit"
				fileString := dirString + "/creds.json"

				err = createCredDirIfNotExist(dirString)
				checkErr(err)

				b, err := json.Marshal(answers)
				checkErr(err)

				ioutil.WriteFile(fileString, b, 0766)

				checkErr(err)

				return nil
			},
		},
	}

	app.Run(os.Args)
}
