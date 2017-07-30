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
	"github.com/caojia/gip/log"
	"runtime"
)

var direnvVersion = "v2.12.2"

func setupBash() {
	exec.Command("echo", "'eval $(direnv hook bash)'", ">>", "~/.bashrc")
}

func setupZsh() {
	exec.Command("echo", "'eval $(direnv hook zsh)'", ">>", "~/.zshrc")
}

func setupDarwin() {
	_, err := exec.Command("command", "-v", "brew").Output()
	// if brew is installed, use brew, otherwise, use default
	if err == nil {
		output, err := exec.Command("brew", "install", "direnv").CombinedOutput()
		if err != nil {
			log.Error(string(output))
			return
		}
	} else {
		setupDefault()
	}
}

func setupDefault() {
	//https://github.com/direnv/direnv/releases/download/v2.12.2/direnv.darwin-386
	url := fmt.Sprintf("https://github.com/direnv/direnv/releases/download/%s/direnv.%s-%s",
		direnvVersion, runtime.GOOS, runtime.GOARCH)
	exec.Command("mkdir", "-p", "~/bin")
	exec.Command("wget", "-O", "~/bin/direnv", url)
	exec.Command("chmod", "+x", "~/bin/direnv")
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install and setup the required commands.",
	Long: `Install and setup the required commands. Including:

- install direnv
- setup direnv for bash/zsh`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := exec.Command("command", "-v", "direnv").Output()
		if err == nil {
			log.Info(`direnv is already installed, please make sure it is set up properly.

For more information, check out: https://github.com/direnv/direnv
			`)
			return
		}
		switch runtime.GOOS {
		case "darwin":
			setupDarwin()
		case "windows":
			panic("Sorry, we don't support windows yet.")
		default:
			setupDefault()
		}

		setupBash()
		setupZsh()
		log.Info(`If you're not using bash or zsh, please follow https://github.com/direnv/direnv to setup direnv for your shell.`)
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
