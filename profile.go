package main

type Profile struct {
	Name     string `json: name`
	Username string `json: username`
	Password string `json: password`
	Nick     string `json: nick`
}
