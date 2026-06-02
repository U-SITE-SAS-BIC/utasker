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
	"strconv"
	"strings"

	"github.com/U-SITE-SAS-BIC/utasker/color"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/spf13/cobra"
)

var (
	editTitle string
	editDesc  string
	editPri   int
	editTags  string
	editDue   string
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a task",
	Long: `Edit fields of an existing task.

All flags are optional; only provided fields are updated.

Examples:
  task edit 5 --title "New title"
  task edit 5 --desc "Updated description" --priority 4`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: invalid task ID")
			os.Exit(1)
		}

		var tags []string
		if editTags != "" {
			tags = strings.Split(editTags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		task, err := db.UpdateTask(id, editTitle, editDesc, editPri, tags, editDue)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}
		fmt.Printf("  %s %s updated: %s\n", color.CyanS("↻"), color.IDLabel(task.ID), task.Title)
	},
}

func init() {
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "New title")
	editCmd.Flags().StringVarP(&editDesc, "desc", "d", "", "New description")
	editCmd.Flags().IntVarP(&editPri, "priority", "p", 0, "New priority (1-5)")
	editCmd.Flags().StringVarP(&editTags, "tags", "", "", "New tags (comma-separated)")
	editCmd.Flags().StringVarP(&editDue, "due", "", "", "New due date (YYYY-MM-DD)")
	rootCmd.AddCommand(editCmd)
}
