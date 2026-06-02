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

var projectCmd = &cobra.Command{
	Use:   "project [name]",
	Short: "Show or set the current project",
	Long: `With no arguments, shows the current project detected from the .task-project file.
With a project name, creates a .task-project file linking this directory to that project.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if err := db.InitProject(args[0]); err != nil {
				fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
				os.Exit(1)
			}
			fmt.Printf("  %s Project set to: %s\n", color.GreenS("✓"), color.BlueS(args[0]))
			return
		}

		project, err := db.GetProjectFromDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}
		if project == "" {
			fmt.Println(color.YellowS("No project set."), "Use 'task init <name>' or 'task project <name>'")
		} else {
			fmt.Printf("  %s %s\n", color.GreenS("→"), color.BlueS(project))
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
