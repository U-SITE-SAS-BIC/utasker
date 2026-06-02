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

	"github.com/U-SITE-SAS-BIC/utasker/color"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/U-SITE-SAS-BIC/utasker/models"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark a task as done",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: invalid task ID")
			os.Exit(1)
		}

		task, err := db.UpdateTaskStatus(id, models.StatusDone)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}
		fmt.Printf("  %s %s marked as done: %s\n", color.GreenS("✓"), color.IDLabel(task.ID), task.Title)
	},
}

var undoCmd = &cobra.Command{
	Use:   "undo <id>",
	Short: "Reopen a completed task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: invalid task ID")
			os.Exit(1)
		}

		task, err := db.UpdateTaskStatus(id, models.StatusPending)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedS("Error:"), err)
			os.Exit(1)
		}
		fmt.Printf("  %s %s reopened: %s\n", color.YellowS("○"), color.IDLabel(task.ID), task.Title)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(undoCmd)
}
