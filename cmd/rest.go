/*
Copyright © 2020 actatum

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
	"time"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	Language         string
	Service          string
	ScaffoldPath     string
	rootFiles        = []string{"protoc-gen.sh", "Dockerfile", ".gitignore"}
	cmdServerFiles   = []string{"cmd/server/main.go"}
	pkgApiGRPCFiles  = []string{"pkg/api/grpc/grpc.go"}
	pkgApiHttpFiles  = []string{"pkg/api/http/http.go", "pkg/api/http/routes.go"}
	pkgApiProtoFiles = []string{"pkg/api/proto/service.proto"}
	pkgServiceFiles  = []string{"pkg/service/logic.go", "pkg/service/repository.go"}
	circleCIFiles    = []string{".circleci/config.yml"}
)

type Folder struct {
	Name       string
	SubFolders []Folder
	Files      []string
}

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:     "rest",
	Aliases: []string{"r"},
	Short:   "This command will scaffold out a rest api",
	Long:    `This command will scaffold out a rest api. Supported Languages: (Go). `,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		if err := scaffoldRest(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		t := time.Now()
		fmt.Printf("Scaffolded %s in: %v\n", Service, t.Sub(start))
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
	ScaffoldPath = viper.GetString("root")
	ok, err := hasTemplates()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("template files not found")
	}
	root := setupFolders()
	if err := create(root); err != nil {
		return err
	}

	return nil
}

func setupFolders() Folder {
	sub := topLevelFolders()
	root := Folder{
		Name:       Service,
		SubFolders: sub,
		Files:      rootFiles,
	}

	return root
}

func topLevelFolders() []Folder {
	var folders []Folder
	pkgSub := pkgFolders()
	cmdSub := cmdFolders()
	cmd := Folder{
		Name:       Service + "/cmd",
		SubFolders: cmdSub,
		Files:      nil,
	}
	pkg := Folder{
		Name:       Service + "/pkg",
		SubFolders: pkgSub,
		Files:      nil,
	}
	ci := Folder{
		Name:       Service + "/.circleci",
		SubFolders: nil,
		Files:      circleCIFiles,
	}

	folders = append(folders, cmd, pkg, ci)
	return folders
}

func cmdFolders() []Folder {
	var folders []Folder
	cmdServer := Folder{
		Name:       Service + "/cmd/server",
		SubFolders: nil,
		Files:      cmdServerFiles,
	}

	folders = append(folders, cmdServer)
	return folders
}

func pkgFolders() []Folder {
	var folders []Folder
	apiSub := apiFolders()
	api := Folder{
		Name:       Service + "/pkg/api",
		SubFolders: apiSub,
		Files:      nil,
	}
	service := Folder{
		Name:       Service + "/pkg/service",
		SubFolders: nil,
		Files:      pkgServiceFiles,
	}

	folders = append(folders, api, service)
	return folders
}

func apiFolders() []Folder {
	var folders []Folder
	grpc := Folder{
		Name:       Service + "/pkg/api/grpc",
		SubFolders: nil,
		Files:      pkgApiGRPCFiles,
	}
	http := Folder{
		Name:       Service + "/pkg/api/http",
		SubFolders: nil,
		Files:      pkgApiHttpFiles,
	}
	proto := Folder{
		Name:       Service + "/pkg/api/proto",
		SubFolders: nil,
		Files:      pkgApiProtoFiles,
	}

	folders = append(folders, grpc, http, proto)

	return folders
}

func getRootFolder() Folder {
	service := Folder{
		Name:       Service,
		SubFolders: nil,
		Files:      rootFiles,
	}

	return service
}

func getTopLevelFolders() []Folder {
	cmd := Folder{
		Name:       "cmd",
		SubFolders: nil,
		Files:      nil,
	}
	pkg := Folder{
		Name:       "pkg",
		SubFolders: nil,
		Files:      nil,
	}
	circleCI := Folder{
		Name:       ".circleci",
		SubFolders: nil,
		Files:      circleCIFiles,
	}

	var folders []Folder
	folders = append(folders, cmd, pkg, circleCI)

	return folders
}

func stripComment(code string) string {
	code = strings.Trim(code, "// ")
	return code
}

func hasTemplates() (bool, error) {
	if _, err := os.Stat(ScaffoldPath + "/templates"); !os.IsNotExist(err) {
		return true, err
	}

	return false, nil
}

func create(root Folder) error {
	// Create root directory
	if err := os.Mkdir(root.Name, 0755); err != nil {
		return err
	}
	// Base Case
	if root.SubFolders == nil {
		if err := writeFiles(root); err != nil {
			return err
		}
		return nil
	}

	for _, dir := range root.SubFolders {
		if err := abc(dir); err != nil {
			return err
		}
	}

	if err := writeFiles(root); err != nil {
		return err
	}

	return nil
}

func writeFiles(f Folder) error {
	for _, fileName := range f.Files {
		file, err := os.Create(Service + "/" + fileName)
		if err != nil {
			return err
		}

		contents, err := ioutil.ReadFile(ScaffoldPath + "/templates/Go/" + fileName)
		if err != nil {
			return err
		}

		text := stripComment(string(contents))

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
