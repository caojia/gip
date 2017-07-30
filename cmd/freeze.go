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

	"github.com/spf13/cobra"
	"strings"
	"path/filepath"
	"os"
	"os/exec"
	"log"
	"go/build"
	"github.com/caojia/gip/helper"
)

// freezeCmd represents the freeze command
var freezeCmd = &cobra.Command{
	Use:   "freeze",
	Short: "Output installed packages in requirements format.",
	Long: `Output installed packages in requirements format.
Only packages installed and depenced directly or indirectly by current package will be output.
packages are listed in a case-insensitive sorted order.

Edit gip.yml if the git repo or version of a package need customized.

Example of gip.yml:

	imports:
		- package: golang.org/x/net
		  repo: https://github.com/golang/net
		- package: xxx/x/x
		  repo: an internal url
		  version: master
		  global: true`,

	Run: func(cmd *cobra.Command, args []string) {
		helper.Freeze()
	},
}

func init() {
	RootCmd.AddCommand(freezeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// freezeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// freezeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
