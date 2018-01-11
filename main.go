package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var gitInfoPath = ".git/config"

var qs = []*survey.Question{
	{
		Name:   "name",
		Prompt: &survey.Input{Message: "Please enter the name associated with the account"},
	},
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

func createDirIfNotExist(path string) (doesExist bool, err error) {
	doesExist, err = fileDoesExist(path)

	if doesExist == true {
		return
	}

	err = os.Mkdir(path, 0766)
	checkErr(err)

	return
}

func createCredDirIfNotExist(path string) (err error) {
	doesExist, err := fileDoesExist(path)
	checkErr(err)

	if doesExist == false {
		err = os.Mkdir(path, 0766)
	}
	return
}

func readFileIfExist(path string) (content []byte) {
	// error handles later
	doesExist, _ := fileDoesExist(path)

	if doesExist == true {
		content, err := ioutil.ReadFile(path)
		checkErr(err)

		return content
	}

	return
}

var Creds = struct {
	MainProfile Profile                `json: "mainProfile"`
	Profiles    map[string]interface{} `json: "profiles"`
}{}

type Profile struct {
	Name     string `json: name`
	Username string `json: username`
	Password string `json: password`
}

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello friend!")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "addUser",
			Aliases: []string{"au"},
			Usage:   "add a new github user account",
			Action: func(c *cli.Context) error {
				var nilProfile = Profile{}
				answers := Profile{}

				err := survey.Ask(qs, &answers)
				checkErr(err)

				homeDir, err := homedir.Dir()
				checkErr(err)

				dirString := homeDir + "/.gogit"
				fileString := dirString + "/creds.json"

				err = createCredDirIfNotExist(dirString)
				checkErr(err)

				content := readFileIfExist(fileString)

				err = json.Unmarshal(content, &Creds)
				checkErr(err)

				// If MainProfile doesn't exist, make the profile the MainProfile
				if Creds.MainProfile == nilProfile {
					Creds.MainProfile = answers
				}

				if Creds.Profiles == nil {
					Creds.Profiles = make(map[string]interface{})
				}

				Creds.Profiles[answers.Username] = answers

				b, err := json.Marshal(Creds)
				checkErr(err)

				ioutil.WriteFile(fileString, b, 0766)

				return nil
			},
		},

		{
			Name:    "createDir",
			Aliases: []string{"cd"},
			Usage:   "create a new git repo with your current stored git profile",
			Action: func(c *cli.Context) error {
				// get current path
				ex, err := os.Executable()
				checkErr(err)
				exPath := path.Dir(ex)

				// create path from the current path + arguments
				dirPath := exPath + "/" + c.Args().First()

				doesExist, err := createDirIfNotExist(dirPath)
				if doesExist == true {
					fmt.Println("That directory already exists.  Please either delete the directory or try again")
					return nil
				}

				err = os.Chdir(dirPath)
				checkErr(err)

				cmd := exec.Command("git", "init")

				_, err = cmd.Output()
				checkErr(err)

				return nil
			},
		},

		{
			Name:    "changeAcct",
			Aliases: []string{"ca"},
			Usage:   "change the account tied to the Git repo",
			Action: func(c *cli.Context) error {
				ex, err := os.Executable()
				checkErr(err)
				exPath := path.Dir(ex)

				content := readFileIfExist(exPath + "/" + gitInfoPath)
				sc := string(content)
				ui := strings.Index(sc, "[user]")

				// if the options for user in the git config exist, delete them
				if ui != -1 {
					sc = sc[:ui]
					fmt.Println(sc)
				}

				// append the new user options to the git config

				return nil
			},
		},
	}

	app.Run(os.Args)
}
