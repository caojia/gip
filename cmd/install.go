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
	"errors"

	"github.com/caojia/gip/helper"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install requirements.txt",
	Short: "Install the dependencies from requirements file",
	Long:  `Install the dependencies from requirements file，gip will look up [requirements.txt] and install dependencies recursively.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("please specify a requirements file.")
		}
		pkgs, err := helper.LoadPackagesFromFile(args[0])
		if err != nil {
			return err
		}
		helper.GetRecursively(pkgs, args[0])
		return nil
	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	installCmd.PersistentFlags().Bool("ignore-global", false, "if it is set to true, none of the packages will be installed in global GOPATH.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
