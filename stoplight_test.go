package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

var yamlExample = []byte(`
server:
  url: http://localhost:8080/cc.xml
  username: username
  password: password
  #only_projects:
  projects:
    - ^project
    - f([a-z]+)_
    - .*::.*::.*
    - ^test
    - mail
  ignored_projects:
    - ^project2
    - ^foo_
    - smail
    - f([a-z]+)_
`)

func mockConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadConfig(bytes.NewBuffer(yamlExample)); err != nil {
		log.Fatalf("Error reading config, %s", err)
	}
}

func TestXMLParser(t *testing.T) {
	t.Skip("skipping test for now.")
}
func TestLastBuildStatus(t *testing.T) {
	t.Skip("should convert XML attribute")
}
func TestCurrentStatus(t *testing.T) {
	t.Skip("should convert XML attribute")
}

func TestExcludeFilter(t *testing.T) {
	mockConfig()
	type testpair struct {
		arg      Project
		expected bool
	}

	tests := []testpair{
		{
			Project{Name: "foo_bar"},
			false,
		},
		{
			Project{Name: "smailer"},
			false,
		},
		{
			Project{Name: "project1"},
			true,
		},
		{
			Project{Name: "project2"},
			false,
		},
		{
			Project{Name: "project3"},
			true,
		},
	}

	for _, pair := range tests {
		result := FilterExcludedProjects(pair.arg)
		if result != pair.expected {
			t.Error(
				"For", pair.arg,
				"Expected", pair.expected,
				"Got", result,
			)
		}
	}
}

func TestIncludeFilter(t *testing.T) {
	mockConfig()
	type testpair struct {
		arg      Project
		expected bool
	}

	tests := []testpair{
		{
			Project{Name: "project1"},
			true,
		},
		{
			Project{Name: "project2"},
			true,
		},
		{
			Project{Name: "project3"},
			true,
		},
		{
			Project{Name: "nope"},
			false,
		},
	}

	for _, pair := range tests {
		result := FilterIncludeProjects(pair.arg)
		if result != pair.expected {
			t.Error(
				"For", pair.arg,
				"Expected", pair.expected,
				"Got", result,
			)
		}
	}
}
