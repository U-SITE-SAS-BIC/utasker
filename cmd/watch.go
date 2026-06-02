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
	"sort"
	"strings"
	"time"

	"github.com/U-SITE-SAS-BIC/utasker/branding"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/U-SITE-SAS-BIC/utasker/models"
	"github.com/spf13/cobra"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Live dashboard — auto-refreshes every N seconds",
	Long: `Shows a live-updating task board organized by project and status.
Like "top" but for your tasks.

Examples:
  utasker watch
  utasker watch --interval 5`,
	Aliases: []string{"w", "live"},
	Run: func(cmd *cobra.Command, args []string) {
		project := ""

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
				renderWatchByProject(project)
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

func renderWatchByProject(project string) {
	clearScreen()
	fmt.Print(branding.TagLine())
	fmt.Println()

	tasks, err := db.ListTasks(db.ListOpts{All: true})
	if err != nil {
		fmt.Println("  Error:", err)
		return
	}

	byProject := make(map[string][]models.Task)
	projectNames := []string{}

	for _, t := range tasks {
		if project != "" && t.Project != project {
			continue
		}
		p := t.Project
		if p == "" {
			p = "(no project)"
		}
		if _, ok := byProject[p]; !ok {
			projectNames = append(projectNames, p)
		}
		byProject[p] = append(byProject[p], t)
	}

	sort.Strings(projectNames)

	totalTasks := len(tasks)
	var totalPending, totalDone int
	now := time.Now().Format("2006-01-02")
	var overdue int

	for _, p := range projectNames {
		ts := byProject[p]
		var pending, doneList []models.Task

		for _, t := range ts {
			switch t.Status {
			case models.StatusPending:
				pending = append(pending, t)
				totalPending++
				if t.DueDate != "" && t.DueDate < now {
					overdue++
				}
			case models.StatusDone:
				doneList = append(doneList, t)
				totalDone++
			}
		}

		projectLabel := p
		color := branding.CyanS
		if p == "(no project)" {
			color = branding.GrayS
		}

		fmt.Printf("  %s %s\n", color(branding.BoldS("["+projectLabel+"]")), branding.GrayS(fmt.Sprintf("(%d)", len(ts))))

		if len(pending) > 0 {
			fmt.Printf("    %s\n", branding.YellowS("○ Pending"))
			for _, t := range pending {
				parts := []string{branding.CyanS(fmt.Sprintf("#%d", t.ID))}
				if t.Priority > 0 {
					parts = append(parts, branding.PriorityLabel(t.Priority))
				}
				if len(t.Tags) > 0 {
					parts = append(parts, "("+strings.Join(t.Tags, ", ")+")")
				}
				if t.DueDate != "" {
					due := t.DueDate
					if due < now {
						due = branding.RedS(due)
					}
					parts = append(parts, branding.GrayS("due:"+due))
				}
				parts = append(parts, t.Title)
				fmt.Println("      " + strings.Join(parts, " "))
			}
		}

		if len(doneList) > 0 {
			fmt.Printf("    %s\n", branding.GreenS("✓ Done"))
			for _, t := range doneList {
				fmt.Printf("      %s %s\n", branding.CyanS(fmt.Sprintf("#%d", t.ID)), branding.GrayS(t.Title))
			}
		}
		fmt.Println()
	}

	if totalTasks == 0 {
		fmt.Println("  No tasks yet.")
		fmt.Println()
		return
	}

	fmt.Printf("  %s\n", branding.BoldS("Summary"))
	fmt.Printf("  %s Total: %d  |  %s  |  %s  |  %s\n",
		branding.GrayS("─"),
		totalTasks,
		branding.YellowS(fmt.Sprintf("Pending: %d", totalPending)),
		branding.GreenS(fmt.Sprintf("Done: %d", totalDone)),
		branding.RedS(fmt.Sprintf("Overdue: %d", overdue)))

	fmt.Printf("  %s Refreshing every %ds · Ctrl+C to quit\n",
		branding.GrayS("─"), watchInterval)
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2, "Refresh interval in seconds")
	rootCmd.AddCommand(watchCmd)
}
