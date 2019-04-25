package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// PromptRes represents a question for the user.
// The message is the question.
// It may uses a default value for the user to accept (by just pressing enter)
type PromptRes struct {
	Message      string
	DefaultValue string
}

type DockerfileData struct {
	Name        string
	BaseImage   string
	BaseVersion string
}
type AliasData struct {
	Name string
}

func main() {

	flag.Parse()
	projectName := flag.Arg(0)

	if projectName == "" {
		fmt.Println("Please provide the name of the application.")
		return
	}

	if _, err := os.Stat(projectName); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(".", projectName), 0770)
	} else {
		fmt.Printf("Folder \"%s\" already exist.\n Aborting...\n", projectName)
		return
	}

	fmt.Printf("Generating a new dockerized application : %s\n", projectName)
	// Dockerfile
	dockerfileData := DockerfileData{
		Name:        projectName,
		BaseImage:   askQuestion("base-image"),
		BaseVersion: askQuestion("base-version"),
	}

	tDockerfile := template.Must(template.New("tmplDockerfile").Parse(tmplDockerfile))
	dockerfile, err := os.Create(filepath.Join(projectName, "Dockerfile"))
	if err != nil {
		fmt.Println("Error creating the template :", err)
		return
	}

	err = tDockerfile.Execute(dockerfile, dockerfileData)
	if err != nil {
		fmt.Println("Error creating the template :", err)
		return
	}
	// alias.sh
	aliasData := AliasData{
		Name: projectName,
	}

	tAlias := template.Must(template.New("tmplAlias").Parse(tmplAlias))
	alias, err := os.Create(filepath.Join(projectName, "alias.sh"))
	if err != nil {
		fmt.Println("Error creating the template :", err)
		return
	}
	err = tAlias.Execute(alias, aliasData)
	if err != nil {
		fmt.Println("Error creating the template :", err)
		return
	}
}

func askQuestion(question string) string {
	var input string

	questions := map[string]PromptRes{
		"base-image":   {Message: "Base image : ", DefaultValue: "alpine"},
		"base-version": {Message: "Base image version : ", DefaultValue: "latest"},
	}

	fmt.Println(questions[question].Message)
	fmt.Scanln(&input)
	if input == "" {
		return questions[question].DefaultValue
	}
	return input
}
