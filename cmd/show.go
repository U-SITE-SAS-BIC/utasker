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

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show task details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error: invalid task ID"))
			os.Exit(1)
		}

		task, err := db.GetTask(id)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Printf("  %s %s\n", color.StatusIcon(task.Status), color.BoldS(task.Title))
		fmt.Printf("  %s %s\n", color.DimS("ID:"), color.IDLabel(task.ID))
		fmt.Printf("  %s %s\n", color.DimS("Status:"), color.StatusIcon(task.Status)+" "+task.Status)
		if task.Project != "" {
			fmt.Printf("  %s %s\n", color.DimS("Project:"), color.ProjectLabel(task.Project))
		}
		if task.Priority > 0 {
			fmt.Printf("  %s %s\n", color.DimS("Priority:"), color.PriorityLabel(task.Priority))
		}
		if task.Description != "" {
			fmt.Printf("  %s %s\n", color.DimS("Desc:"), task.Description)
		}
		if len(task.Tags) > 0 {
			tagStrs := make([]string, len(task.Tags))
			for i, t := range task.Tags {
				tagStrs[i] = color.TagLabel(t)
			}
			fmt.Printf("  %s %s\n", color.DimS("Tags:"), strings.Join(tagStrs, ", "))
		}
		if task.DueDate != "" {
			fmt.Printf("  %s %s\n", color.DimS("Due:"), task.DueDate)
		}
		fmt.Printf("  %s %s\n", color.DimS("Created:"), task.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("  %s %s\n", color.DimS("Updated:"), task.UpdatedAt.Format("2006-01-02 15:04"))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
