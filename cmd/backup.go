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
	"path/filepath"
	"time"

	"github.com/U-SITE-SAS-BIC/utasker/branding"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup tasks.json to ~/.tasker/backups/",
	Run: func(cmd *cobra.Command, args []string) {
		doBackup()
		fmt.Printf("  %s Backup created\n", branding.GreenS("✓"))
	},
}

func doBackup() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	backupDir := filepath.Join(home, ".tasker", "backups")
	os.MkdirAll(backupDir, 0755)

	src := filepath.Join(home, ".tasker", "tasks.json")
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return
	}

	ts := time.Now().Format("2006-01-02_150405")
	dst := filepath.Join(backupDir, fmt.Sprintf("tasks-%s.json", ts))

	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Backup error:", err)
		return
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "Backup error:", err)
	}
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
