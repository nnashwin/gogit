package main

import (
	"strings"
	"testing"
)

func TestCreateConfigString(t *testing.T) {
	expected := true
	testString := createConfigString("cookies", "cake", "candies")
	actual := strings.Contains(testString, "cookies")

	if actual != expected {
		t.Fail()
	}

	actual = strings.Contains(testString, "cake")
	if actual != expected {
		t.Fail()
	}

	actual = strings.Contains(testString, "candies")
	if actual != expected {
		t.Fail()
	}
}

func TestDoesFileExist(t *testing.T) {
	expected := false
	actual := doesFileExist("./cookies")

	if actual != expected {
		t.Fail()
	}

	expected = true
	actual = doesFileExist("./main.go")

	if actual != expected {
		t.Fail()
	}
}

func TestGetCredPathString(t *testing.T) {
	expected := "./.gogit/creds.json"
	actual := getCredPathString("./")

	if actual != expected {
		t.Fail()
	}
}
