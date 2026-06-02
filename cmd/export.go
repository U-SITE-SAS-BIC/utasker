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

var (
	exportFile   string
	exportAll    bool
	exportStatus string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export tasks to plain text",
	Long: `Export tasks to TXT format (stdout or file).

Examples:
  task export
  task export --file tasks.txt
  task export -a -s done --file done.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		project := getProject()
		if exportAll {
			project = ""
		}

		tasks, err := db.ListTasks(project, exportStatus, exportAll)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}

		var b strings.Builder

		fmt.Fprintf(&b, "Tasks Report\n")
		fmt.Fprintf(&b, "════════════\n")
		fmt.Fprintf(&b, "Generated: %s\n", "2026-06-01 18:24")
		fmt.Fprintf(&b, "Filter: project=%q status=%q\n\n", project, exportStatus)

		if len(tasks) == 0 {
			fmt.Fprintln(&b, "No tasks found.")
		} else {
			for _, t := range tasks {
				status := "○"
				if t.Status == models.StatusDone {
					status = "✓"
				} else if t.Status == models.StatusCancelled {
					status = "✕"
				}

				fmt.Fprintf(&b, "%s #%d", status, t.ID)

				if t.Project != "" {
					fmt.Fprintf(&b, " [%s]", t.Project)
				}

				if t.Priority > 0 {
					fmt.Fprintf(&b, " P%d", t.Priority)
				}

				if len(t.Tags) > 0 {
					fmt.Fprintf(&b, " (%s)", strings.Join(t.Tags, ", "))
				}

				if t.DueDate != "" {
					fmt.Fprintf(&b, " due:%s", t.DueDate)
				}

				fmt.Fprintf(&b, " %s", t.Title)

				if t.Description != "" {
					fmt.Fprintf(&b, "\n   ↳ %s", t.Description)
				}

				fmt.Fprintln(&b)
			}
		}

		output := b.String()

		if exportFile != "" {
			if err := os.WriteFile(exportFile, []byte(output), 0644); err != nil {
				fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
				os.Exit(1)
			}
			fmt.Printf("  %s Exported to %s\n", color.GreenS("✓"), color.CyanS(exportFile))
		} else {
			fmt.Print(output)
		}
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "Output file (default: stdout)")
	exportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export from all projects")
	exportCmd.Flags().StringVarP(&exportStatus, "status", "s", "", "Filter by status (pending, done, cancelled)")
	rootCmd.AddCommand(exportCmd)
}
