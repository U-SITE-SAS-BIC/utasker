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
	"os/exec"
	"strings"
	"time"

	"github.com/U-SITE-SAS-BIC/utasker/branding"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/U-SITE-SAS-BIC/utasker/models"
	"github.com/spf13/cobra"
)

var (
	watchInterval int
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Live dashboard — auto-refreshes every N seconds",
	Long: `Shows a live-updating task board that refreshes automatically.
Like "top" but for your tasks.

Examples:
  utasker watch
  utasker watch --interval 5
  utasker watch -a`,
	Aliases: []string{"w", "live"},
	Run: func(cmd *cobra.Command, args []string) {
		project := getProject()
		if allFlag {
			project = ""
		}

		lastMod := time.Time{}
		first := true

		for {
			mod, err := db.TasksChangedAt()
			if err != nil {
				mod = time.Now()
			}

			if mod != lastMod || first {
				lastMod = mod
				first = false
				renderWatch(project)
			}

			time.Sleep(time.Duration(watchInterval) * time.Second)
		}
	},
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func renderWatch(project string) {
	clearScreen()

	fmt.Print(branding.TagLine())
	fmt.Println()

	opts := db.ListOpts{All: true}
	tasks, err := db.ListTasks(opts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var pending, doneP, cancelled []models.Task
	today := time.Now().Format("2006-01-02")

	for _, t := range tasks {
		if project != "" && t.Project != project {
			continue
		}
		switch t.Status {
		case models.StatusPending:
			pending = append(pending, t)
		case models.StatusDone:
			doneP = append(doneP, t)
		case models.StatusCancelled:
			cancelled = append(cancelled, t)
		}
	}

	total := len(pending) + len(doneP) + len(cancelled)

	if total == 0 {
		fmt.Println("  No tasks yet.")
		fmt.Println()
		return
	}

	printLiveSection("○ Pending", pending, allFlag)
	printLiveSection("✓ Done", doneP, allFlag)
	printLiveSection("✕ Cancelled", cancelled, allFlag)

	var overdue int
	for _, t := range pending {
		if t.DueDate != "" && t.DueDate < today {
			overdue++
		}
	}

	fmt.Printf("\n  %s  Total: %d  |  Pending: %d  |  Done: %d  |  Overdue: %d\n",
		branding.GrayS("─"),
		total, len(pending), len(doneP), overdue)

	fmt.Printf("  %s  Refreshing every %ds · Ctrl+C to quit\n",
		branding.GrayS("─"), watchInterval)
}

func printLiveSection(label string, tasks []models.Task, showProject bool) {
	if len(tasks) == 0 {
		return
	}
	colorized := label
	if strings.Contains(label, "Pending") {
		colorized = branding.YellowS(label)
	} else if strings.Contains(label, "Done") {
		colorized = branding.GreenS(label)
	} else if strings.Contains(label, "Cancelled") {
		colorized = branding.RedS(label)
	}
	fmt.Printf("  %s (%d)\n", branding.BoldS(colorized), len(tasks))
	for _, t := range tasks {
		parts := []string{branding.CyanS(fmt.Sprintf("#%d", t.ID))}
		if showProject && t.Project != "" {
			parts = append(parts, branding.BlueS("["+t.Project+"]"))
		}
		if t.Priority > 0 {
			parts = append(parts, branding.PriorityLabel(t.Priority))
		}
		if len(t.Tags) > 0 {
			parts = append(parts, "("+strings.Join(t.Tags, ", ")+")")
		}
		if t.DueDate != "" {
			parts = append(parts, branding.GrayS("due:"+t.DueDate))
		}
		parts = append(parts, t.Title)
		fmt.Println("    " + strings.Join(parts, " "))
	}
}

func init() {
	watchCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Show all projects")
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2, "Refresh interval in seconds")
	rootCmd.AddCommand(watchCmd)
}
