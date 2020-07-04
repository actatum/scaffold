/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Language string
	Service string
	serverFiles = []string{"main.go"}
	k8sFiles = []string{"service.yml", "deploy.yml"}
	apiFiles = []string{"http.go", "routes.go"}
	serviceFiles = []string{"logic.go", "model.go", "repository.go", "service.go"}
	circleCIFiles = []string{"config.yml"}
)

type TopLevelFolder struct {
	Name string
	SubFolders []SubFolder
	Files []string
}

type SubFolder struct {
	Name string
	Files []string
}

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "This command will scaffold out a rest api",
	Long: `This command will scaffold out a rest api. Supported Languages: (Go, Python). `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rest called")
		if err := scaffoldRest(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restCmd)

	// Here you will define your flags and configuration settings.
	restCmd.Flags().StringVarP(&Language, "language", "l", "", "programming language of choice")
	restCmd.Flags().StringVarP(&Service, "service", "s", "", "name of the service you are creating")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scaffoldRest() error {
	if Service == "" {
		return errors.New("--service or -s flag is required")
	}
	folders := getTopLevelFolders()
	folders = fillFolders(folders)
	if err := createProjectFolder(); err != nil {
		return err
	}
	if err := create(folders); err != nil {
		return err
	}

	return nil
}

func getTopLevelFolders() []TopLevelFolder {
	cmd := TopLevelFolder{
		Name:      	"cmd",
		SubFolders: nil,
		Files:      nil,
	}
	k8s := TopLevelFolder{
		Name:      "k8s",
		SubFolders: nil,
		Files:      k8sFiles,
	}
	pkg := TopLevelFolder{
		Name:       "pkg",
		SubFolders: nil,
		Files:      nil,
	}
	circleCI := TopLevelFolder{
		Name:       ".circleci",
		SubFolders: nil,
		Files:      circleCIFiles,
	}

	var folders []TopLevelFolder
	folders = append(folders, cmd, k8s, pkg, circleCI)

	return folders
}

func fillFolders(folders []TopLevelFolder) []TopLevelFolder {
	var filled []TopLevelFolder
	for _, top := range folders {
		fmt.Println(top.Name)
		switch top.Name {
		case "cmd":
			server := SubFolder{
				Name:  "server",
				Files: serverFiles,
			}
			top.SubFolders = append(top.SubFolders, server)

		case "pkg":
			api := SubFolder{
				Name:  "api",
				Files: apiFiles,
			}
			service := SubFolder{
				Name:  "service",
				Files: serviceFiles,
			}

			top.SubFolders = append(top.SubFolders, api, service)

		}
		filled = append(filled, top)

	}

	return filled
}

func create(structure []TopLevelFolder) error {
	for _, dir := range structure {
		if err := os.Mkdir(Service + "/" + dir.Name, 0755); err != nil {
			return err
		}
		if err := writeTop(dir); err != nil {
			return err
		}
	}

	return nil
}

func createProjectFolder() error {
	if err := os.Mkdir(Service, 0755); err != nil {
		return err
	}

	return nil
}

func StripComment(code string) string {
	code = strings.Trim(code, "// ")
	return code
}

func writeTop(folder TopLevelFolder) error {
	if folder.SubFolders != nil {
		for _, sub := range folder.SubFolders {
			if err := writeSub(folder, sub); err != nil {
				return err
			}
		}
	}

	if err := writeTopTemplate(folder); err != nil {
		return err
	}

	return nil
}

func writeSub(top TopLevelFolder, sub SubFolder) error {
	if err := os.MkdirAll(Service + "/" + top.Name + "/" + sub.Name, 0755); err != nil {
		return err
	}

	if err := writeSubTemplate(top, sub); err != nil {
		return err
	}

	return nil
}

func writeTopTemplate(top TopLevelFolder) error {
	for _, fileName := range top.Files {
		file, err := os.Create(Service + "/" + top.Name + "/" + fileName)
		if err != nil {
			return err
		}
		contents, err := ioutil.ReadFile("templates/Go/" + top.Name + "/" + fileName)
		if err != nil {
			return err
		}

		text := StripComment(string(contents))

		_, err = file.Write([]byte(text))
		if err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func writeSubTemplate(top TopLevelFolder, sub SubFolder) error {
	for _, fileName := range sub.Files {
		file, err := os.Create(Service + "/" + top.Name + "/" + sub.Name + "/" + fileName)
		if err != nil {
			return err
		}
		contents, err := ioutil.ReadFile("templates/Go/" + top.Name + "/" + sub.Name + "/" + fileName)
		if err != nil {
			return err
		}

		text := StripComment(string(contents))

		_, err = file.Write([]byte(text))
		if err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}