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
	"time"

	"github.com/U-SITE-SAS-BIC/utasker/color"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/U-SITE-SAS-BIC/utasker/models"
	"github.com/spf13/cobra"
)

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Show full panorama of all tasks",
	Long: `Shows tasks grouped by status with counts and overdue highlights.

Examples:
  task board
  task board -a`,
	Aliases: []string{"b", "panorama"},
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.ListTasks("", "", true)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println(color.DimS("No tasks yet. Create one with 'task add \"something\"'"))
			return
		}

		var pending, doneP, cancelled []models.Task
		var overdue []models.Task
		today := time.Now().Format("2006-01-02")

		for _, t := range tasks {
			switch t.Status {
			case models.StatusPending:
				pending = append(pending, t)
				if t.DueDate != "" && t.DueDate < today {
					overdue = append(overdue, t)
				}
			case models.StatusDone:
				doneP = append(doneP, t)
			case models.StatusCancelled:
				cancelled = append(cancelled, t)
			}
		}

		total := len(pending) + len(doneP) + len(cancelled)

		printSection("Pending", pending, total)
		printSection("Done", doneP, total)
		printSection("Cancelled", cancelled, total)

		if len(overdue) > 0 {
			fmt.Printf("\n  %s %s\n", color.RedS("!"), color.BoldS("Overdue"))
			for _, t := range overdue {
				fmt.Printf("    %s %s %s %s\n",
					color.RedS("!"),
					color.IDLabel(t.ID),
					color.RedS(t.DueDate),
					t.Title)
			}
		}

		fmt.Printf("\n  %s %s\n", color.DimS("───"), color.BoldS("Summary"))
		fmt.Printf("  %s %s\n", color.DimS(fmt.Sprintf("Total: %d", total)), "")
		if len(overdue) > 0 {
			fmt.Printf("  %s\n", color.RedS(fmt.Sprintf("Overdue: %d", len(overdue))))
		}
		fmt.Println()
	},
}

func printSection(label string, tasks []models.Task, total int) {
	if len(tasks) == 0 {
		return
	}

	icon := "○"
	labelColor := color.YellowS
	if label == "Done" {
		icon = "✓"
		labelColor = color.GreenS
	} else if label == "Cancelled" {
		icon = "✕"
		labelColor = color.RedS
	}

	fmt.Printf("\n  %s %s (%d/%d)\n", labelColor(icon), color.BoldS(label), len(tasks), total)

	for _, t := range tasks {
		parts := []string{color.IDLabel(t.ID)}

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

		title := t.Title
		if label == "Done" {
			title = color.DimS(title)
		}

		parts = append(parts, title)
		fmt.Println("    " + strings.Join(parts, " "))
	}
}

func init() {
	boardCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Show projects column in board")
	rootCmd.AddCommand(boardCmd)
}
