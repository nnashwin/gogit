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

func GetCredPathString(basePath string) string {
	return basePath + "/.gogit/creds.json"
}

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		fmt.Println("Welcome to Gogit!  Add a git credential (gogit au) to begin!!")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "getMain",
			Aliases: []string{"gm"},
			Usage:   "get the current Main Profile in the creds file",
			Action: func(c *cli.Context) error {
				homeDir, err := homedir.Dir()
				checkErr(err)

				if doesFileExist(GetCredPathString(homeDir)) == false {
					return errors.New("You currently do not have a cred file.  Run the addUser (au) command to configure a cred file")
				}

				creds := readFile(GetCredPathString(homeDir))

				err = json.Unmarshal(creds, &Creds)
				checkErr(err)

				if Creds.MainProfile == (Profile{}) {
					return errors.New("You currently have an empty Main Profile.  Run the addUser (au) command to create one.")
				}

				fmt.Printf("Your Main Profile:\n nick: %s\n username: %s\n name: %s", Creds.MainProfile.Nick, Creds.MainProfile.Nick, Creds.MainProfile.Name)

				return nil
			},
		},

		{
			Name:    "changeMain",
			Aliases: []string{"cm"},
			Usage:   "change the current Main Profile to a different Profile in the creds file",
			Action: func(c *cli.Context) error {
				homeDir, err := homedir.Dir()
				checkErr(err)

				if doesFileExist(GetCredPathString(homeDir)) == false {
					return errors.New("You currently do not have a cred file.  Run the addUser (au) command to configure a cred file")
				}

				creds := readFile(GetCredPathString(homeDir))

				err = json.Unmarshal(creds, &Creds)
				checkErr(err)

				if Creds.MainProfile == (Profile{}) {
					return errors.New("You currently have an empty Main Profile.  Run the addUser (au) command to create one.")
				}

				if Creds.Profiles[c.Args().First()] == (Profile{}) {
					return errors.New(fmt.Sprintf("There is no profile stored under the nickname %s, please either create a new user or use a different nickname", c.Args().First()))
				}

				Creds.MainProfile = Creds.Profiles[c.Args().First()]

				b, err := json.Marshal(Creds)
				checkErr(err)

				ioutil.WriteFile(GetCredPathString(homeDir), b, os.ModePerm)

				fmt.Printf("Main Profile changed to %s", Creds.MainProfile.Nick)

				return nil
			},
		},
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
					err = os.Mkdir(dirString, os.ModePerm)
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

				ioutil.WriteFile(fileString, b, os.ModePerm)

				fmt.Printf("Profile \"%s\" added to the creds file", answers.Nick)

				return nil
			},
		},

		{
			Name:    "createDir",
			Aliases: []string{"cd"},
			Usage:   "create a new git repo with your current Main Profile",
			Action: func(c *cli.Context) error {
				// get current path
				wd, err := os.Getwd()
				checkErr(err)

				// create path from the current path + arguments
				dirPath := wd + "/" + c.Args().First()

				doesExist := doesFileExist(dirPath)
				checkErr(err)

				if doesExist == true {
					return errors.New("That directory already exists.  Please either delete the directory or try again")
				}

				err = os.Mkdir(dirPath, os.ModePerm)
				checkErr(err)

				err = os.Chdir(dirPath)
				checkErr(err)

				cmd := exec.Command("git", "init")

				_, err = cmd.Output()
				checkErr(err)

				// get creds from the creds file
				homeDir, err := homedir.Dir()
				checkErr(err)

				if doesFileExist(GetCredPathString(homeDir)) == false {
					return errors.New("You currently do not have a cred file.  Run the addUser (au) command to configure a cred file")
				}

				creds := readFile(GetCredPathString(homeDir))

				err = json.Unmarshal(creds, &Creds)
				checkErr(err)

				if doesFileExist(dirPath+"/"+gitInfoPath) == false {
					return errors.New("The config file in the .git folder does not exist, can not change the account attached")
				}
				content := readFile(dirPath + "/" + gitInfoPath)

				ui := strings.Index(string(content), "[user]")

				// if the options for user in the git config exist, delete them
				if ui != -1 {
					content = content[:ui]
				}

				if Creds.MainProfile == (Profile{}) {
					return errors.New("There was a problem finding the Main Profile used to write the git config file.  Run the addUser (au) command and try again.")
				}

				content = append([]byte(content), createConfigString(Creds.MainProfile.Name, Creds.MainProfile.Username, Creds.MainProfile.Password)...)

				ioutil.WriteFile(dirPath+"/"+gitInfoPath, content, os.ModePerm)

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

				credPath := GetCredPathString(homeDir)

				if doesFileExist(credPath) == false {
					return errors.New("The cred file is empty.  Run the createDir command to add a main account")
				}

				creds := readFile(credPath)
				err = json.Unmarshal(creds, &Creds)

				wd, err := os.Getwd()
				checkErr(err)

				if doesFileExist(wd+"/"+gitInfoPath) == false {
					return errors.New("The config file in the .git folder does not exist, can not change the account attached")
				}
				content := readFile(wd + "/" + gitInfoPath)
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

				// add an empty check after the assignment in case either the Creds.MainProfile or the specific value of the map are empty
				if profile == (Profile{}) {
					return errors.New("There was a problem finding the profile used to write the git config file.  Either change your Nick or Run the addUser (au) command and try again.")
				}

				sc = append([]byte(sc), createConfigString(profile.Name, profile.Username, profile.Password)...)

				ioutil.WriteFile(wd+"/"+gitInfoPath, sc, os.ModePerm)

				return nil
			},
		},
	}

	app.Run(os.Args)
}
