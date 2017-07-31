// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"os/exec"
	"github.com/caojia/gip/log"
)

var envrc = ".envrc"
var template = `mkdir -p %[1]s
export GOPATH=%[1]s:$GOPATH
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init name",
	Short: "Initialize a project with [name]",
	Long: `Initialize a project with [name]; it will create a .envrc file in current workspace, and set GOPATH.

e.g.
gip init gip

This command will generate the following .envrc file:

mkdir -p ~/.gip/gip
export GOPATH=~/.gip/gip:$GOPATH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("not enough arguments")
		}
		path := filepath.Join("~/.gip", args[0])
		content := fmt.Sprintf(template, path)
		ioutil.WriteFile(envrc, []byte(content), os.ModePerm)
		output, err := exec.Command("direnv", "allow", ".").CombinedOutput()
		if err != nil {
			log.Error(string(output))
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
