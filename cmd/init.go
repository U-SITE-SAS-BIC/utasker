// Copyright 2026 Lizandro
//
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
	"os"

	"github.com/U-SITE-SAS-BIC/utasker/color"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project]",
	Short: "Initialize task tracking in this directory",
	Long: `Creates a .task-project file in the current directory linking it to a project.

If no project name is given, the current directory name is used.
Run this once per project to automatically scope tasks.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		project := ""
		if len(args) > 0 {
			project = args[0]
		} else {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			project = dir
		}

		if err := db.InitProject(project); err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}
		fmt.Printf("  %s Project initialized: %s\n", color.GreenS("✓"), color.BlueS(project))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
