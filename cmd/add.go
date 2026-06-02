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
	"strings"

	"github.com/U-SITE-SAS-BIC/utasker/color"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/spf13/cobra"
)

var (
	addPriority  int
	addDesc      string
	addTags      string
	addDueDate   string
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new task",
	Long: `Add a new task to the current project.

Examples:
  task add "Fix login bug"
  task add "Deploy v2" -p 3 -d "Deploy to production" -t "urgent,deploy" --due 2025-01-15`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		project := getProject()
		title := args[0]

		var tags []string
		if addTags != "" {
			tags = strings.Split(addTags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		task, err := db.AddTask(title, project, addDesc, addPriority, tags, addDueDate)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}

		fmt.Printf("  %s %s created\n", color.GreenS("✓"), color.IDLabel(task.ID))
	},
}

func init() {
	addCmd.Flags().IntVarP(&addPriority, "priority", "p", 0, "Priority (1-5)")
	addCmd.Flags().StringVarP(&addDesc, "desc", "d", "", "Description")
	addCmd.Flags().StringVarP(&addTags, "tags", "t", "", "Comma-separated tags")
	addCmd.Flags().StringVarP(&addDueDate, "due", "", "", "Due date (YYYY-MM-DD)")
	rootCmd.AddCommand(addCmd)
}
