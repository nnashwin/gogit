package main

import (
	"encoding/json"
	"errors"
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
		Name:   "nick",
		Prompt: &survey.Input{Message: "Please enter the name you want to store the account under"},
	},
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

func doesFileExist(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	return true
}

func readFile(path string) (content []byte) {
	content, err := ioutil.ReadFile(path)
	checkErr(err)

	return content
}

func createConfigString(name string, username string, password string) string {
	return fmt.Sprintf("[user]\n        name = %s\n        email = %s\n        password = %s", name, username, password)
}

var Creds = struct {
	MainProfile Profile            `json: "mainProfile,omitempty"`
	Profiles    map[string]Profile `json: "profiles,omitempty"`
}{}

type Profile struct {
	Name     string `json: name`
	Username string `json: username`
	Password string `json: password`
	Nick     string `json: nick`
}

func getCredPathString(basePath string) string {
	return basePath + "/.gogit/creds.json"
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

				// create the dir if it doesn't exist
				if doesFileExist(dirString) == false {
					err = os.Mkdir(dirString, 0766)
					checkErr(err)
				}

				fileString := dirString + "/creds.json"

				// create the file if it doesn't exists
				if doesFileExist(fileString) == false {
					_, err := os.Create(fileString)
					checkErr(err)
				}

				content := readFile(fileString)

				if len(content) > 0 {
					err = json.Unmarshal(content, &Creds)
					checkErr(err)
				}

				// If MainProfile doesn't exist, make the profile the MainProfile
				if Creds.MainProfile == nilProfile {
					Creds.MainProfile = answers
				}

				if Creds.Profiles == nil {
					Creds.Profiles = make(map[string]Profile)
				}

				Creds.Profiles[answers.Nick] = answers

				b, err := json.Marshal(Creds)
				checkErr(err)

				ioutil.WriteFile(fileString, b, 0766)

				fmt.Printf("Profile \"%s\" added to the creds file", answers.Nick)

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

				doesExist := doesFileExist(dirPath)
				checkErr(err)

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
				homeDir, err := homedir.Dir()
				checkErr(err)

				credPath := getCredPathString(homeDir)

				if doesFileExist(credPath) == false {
					return errors.New("The cred file is empty.  Run the createDir command to add a main account")
				}

				creds := readFile(credPath)
				err = json.Unmarshal(creds, &Creds)

				ex, err := os.Executable()
				checkErr(err)
				exPath := path.Dir(ex)

				if doesFileExist(exPath+"/"+gitInfoPath) == false {
					fmt.Println("The config file in the .git folder does not exist, can not change the account attached")
					return nil
				}
				content := readFile(exPath + "/" + gitInfoPath)
				sc := content
				ui := strings.Index(string(sc), "[user]")

				// if the options for user in the git config exist, delete them
				if ui != -1 {
					sc = sc[:ui]
				}

				var profile Profile

				if c.Args().First() == "" {
					profile = Creds.MainProfile
				} else {
					profile = Creds.Profiles[c.Args().First()]
				}

				sc = append([]byte(sc), createConfigString(profile.Name, profile.Username, profile.Password)...)

				ioutil.WriteFile(exPath+"/"+gitInfoPath, sc, 0766)

				return nil
			},
		},
	}

	app.Run(os.Args)
}
