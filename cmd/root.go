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

	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/spf13/cobra"
)

var (
	projectFlag string
	allFlag     bool
	statusFlag  string
	Version     = "dev"
	Commit      = "none"
	Date        = "unknown"
)

func getProject() string {
	if projectFlag != "" {
		return projectFlag
	}
	p, err := db.GetProjectFromDir()
	if err == nil && p != "" {
		return p
	}
	return ""
}

var rootCmd = &cobra.Command{
	Use:     "task",
	Version: Version,
	Short:   "Task manager - simple, offline, project-aware",
	Long: `task is a CLI task manager that works offline and connects tasks to your projects.

Create a .task-project file in your project directory to automatically scope tasks:
  task init myproject

Then run commands without specifying the project each time.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&projectFlag, "project", "", "Project name (overrides .task-project file)")
}
