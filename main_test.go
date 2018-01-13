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
