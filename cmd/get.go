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
	"os/exec"
	"github.com/caojia/gip/helper"
	"github.com/pkg/errors"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Extension of go get, support a customized url and version for given package",
	Long: `Extension of go get, there are 2 use cases:
1. use --, all parameters after -- will pass to go get
2. use the format url#version,pkg, it will download the pkg from url with specified version.

e.g.
gip get -- -u golang.org/x/net
gip get github.com/golang/net#master,golang.org/x/net`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dash := cmd.ArgsLenAtDash()
		if dash >= 0 {
			c := exec.Command("go", append([]string{"get"}, args[dash:]...)...)
			o, err := c.CombinedOutput()
			if err != nil {
				fmt.Println(string(o))
			}
			return nil
		} else {
			if len(args) <= 0 {
				return errors.New("args is empty")
			}
			pkg, err := helper.LoadPackage(args[0])
			if err != nil {
				return err
			}
			return helper.Get(pkg)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
