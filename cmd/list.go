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
	"github.com/U-SITE-SAS-BIC/utasker/models"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `List tasks filtered by project and status.

By default shows pending tasks for the current project.
Use -a to show all projects, -A to show all statuses.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		project := getProject()
		if allFlag {
			project = ""
		}

		tasks, err := db.ListTasks(project, statusFlag, allFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println(color.DimS("No tasks found."))
			return
		}

		fmt.Println(color.BoldS(fmt.Sprintf("\nTasks (%d):", len(tasks))))
		for _, t := range tasks {
			icon := color.StatusIcon(t.Status)
			id := color.IDLabel(t.ID)
			title := t.Title

			parts := []string{icon, id}

			if allFlag && t.Project != "" {
				parts = append(parts, color.ProjectLabel(t.Project))
			}

			prio := color.PriorityLabel(t.Priority)
			if prio != "" {
				parts = append(parts, prio)
			}

			if len(t.Tags) > 0 {
				tagStrs := make([]string, len(t.Tags))
				for i, tag := range t.Tags {
					tagStrs[i] = color.TagLabel(tag)
				}
				parts = append(parts, "("+strings.Join(tagStrs, ", ")+")")
			}

			if t.DueDate != "" {
				parts = append(parts, color.DimS("due:"+t.DueDate))
			}

			if t.Status == models.StatusDone {
				title = color.DimS(title)
			}

			parts = append(parts, title)
			fmt.Println("  " + strings.Join(parts, " "))
		}
		fmt.Println()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Show tasks from all projects")
	listCmd.Flags().StringVarP(&statusFlag, "status", "s", "", "Filter by status (pending, done, cancelled)")
	rootCmd.AddCommand(listCmd)
}
