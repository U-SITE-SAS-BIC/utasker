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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/U-SITE-SAS-BIC/utasker/branding"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/U-SITE-SAS-BIC/utasker/models"
	"github.com/spf13/cobra"
)

type remoteConfig struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token,omitempty"`
}

type remoteStore struct {
	Remotes []remoteConfig `json:"remotes"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".tasker")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func loadRemotes() (remoteStore, error) {
	var s remoteStore
	path, err := getConfigPath()
	if err != nil {
		return s, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return s, err
	}
	json.Unmarshal(data, &s)
	return s, nil
}

func saveRemotes(s remoteStore) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote API connections",
	Long: `Configure and sync tasks with a remote API.

Examples:
  utasker remote add prod https://api.example.com/tasks mytoken
  utasker remote list
  utasker push prod
  utasker pull prod`,
}

var remoteAddCmd = &cobra.Command{
	Use:   "add <name> <url> [token]",
	Short: "Add a remote endpoint",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		rc := remoteConfig{
			Name: args[0],
			URL:  args[1],
		}
		if len(args) > 2 {
			rc.Token = args[2]
		}

		s, err := loadRemotes()
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		for i, r := range s.Remotes {
			if r.Name == rc.Name {
				s.Remotes[i] = rc
				goto save
			}
		}
		s.Remotes = append(s.Remotes, rc)

	save:
		if err := saveRemotes(s); err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}
		fmt.Printf("  %s Remote %s added\n", branding.GreenS("✓"), branding.CyanS(rc.Name))
	},
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured remotes",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := loadRemotes()
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}
		if len(s.Remotes) == 0 {
			fmt.Println(branding.YellowS("No remotes configured."))
			return
		}
		for _, r := range s.Remotes {
			token := ""
			if r.Token != "" {
				token = " (token set)"
			}
			fmt.Printf("  %s %s%s\n", branding.CyanS(r.Name), r.URL, branding.GrayS(token))
		}
	},
}

var remoteRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a remote",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, err := loadRemotes()
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}
		for i, r := range s.Remotes {
			if r.Name == args[0] {
				s.Remotes = append(s.Remotes[:i], s.Remotes[i+1:]...)
				if err := saveRemotes(s); err != nil {
					fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
					os.Exit(1)
				}
				fmt.Printf("  %s Remote %s removed\n", branding.RedS("✕"), branding.CyanS(args[0]))
				return
			}
		}
		fmt.Fprintln(os.Stderr, branding.RedS("Remote not found:"), args[0])
		os.Exit(1)
	},
}

var pushCmd = &cobra.Command{
	Use:   "push [remote-name]",
	Short: "Push tasks to remote API",
	Long: `Sends all tasks as JSON to the configured remote API.
If no name is given, uses the first configured remote.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rc := getRemote(args)
		tasks := loadAllTasks()

		body, _ := json.Marshal(tasks)
		req, _ := http.NewRequest("POST", rc.URL, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if rc.Token != "" {
			req.Header.Set("Authorization", "Bearer "+rc.Token)
		}

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		fmt.Printf("  %s Pushed %d tasks to %s\n", branding.GreenS("✓"), len(tasks), branding.CyanS(rc.Name))
		if resp.StatusCode >= 300 {
			respBody, _ := io.ReadAll(resp.Body)
			fmt.Fprintln(os.Stderr, branding.RedS("API error:"), string(respBody))
		}
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull [remote-name]",
	Short: "Pull tasks from remote API",
	Long: `Fetches tasks from the remote API and merges them locally.
If no name is given, uses the first configured remote.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rc := getRemote(args)

		client := &http.Client{Timeout: 30 * time.Second}
		req, _ := http.NewRequest("GET", rc.URL, nil)
		if rc.Token != "" {
			req.Header.Set("Authorization", "Bearer "+rc.Token)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var remoteTasks []models.Task
		if err := json.NewDecoder(resp.Body).Decode(&remoteTasks); err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error decoding response:"), err)
			os.Exit(1)
		}

		// Backup before merging
		doBackup()

		// Simple merge: add tasks that don't exist locally
		local, _ := db.ListTasks(db.ListOpts{All: true})
		localMap := make(map[int]bool)
		for _, t := range local {
			localMap[t.ID] = true
		}

		added := 0
		for _, t := range remoteTasks {
			if !localMap[t.ID] {
				db.AddTask(t.Title, t.Project, t.Description, t.Priority, t.Tags, t.DueDate)
				added++
			}
		}

		fmt.Printf("  %s Pulled from %s: %d new tasks\n", branding.GreenS("✓"), branding.CyanS(rc.Name), added)
	},
}

func getRemote(args []string) remoteConfig {
	s, err := loadRemotes()
	if err != nil {
		fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
		os.Exit(1)
	}
	if len(s.Remotes) == 0 {
		fmt.Fprintln(os.Stderr, branding.YellowS("No remotes configured. Use: utasker remote add <name> <url>"))
		os.Exit(1)
	}
	if len(args) > 0 {
		for _, r := range s.Remotes {
			if r.Name == args[0] {
				return r
			}
		}
		fmt.Fprintln(os.Stderr, branding.RedS("Remote not found:"), args[0])
		os.Exit(1)
	}
	return s.Remotes[0]
}

func loadAllTasks() []models.Task {
	tasks, err := db.ListTasks(db.ListOpts{All: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
		os.Exit(1)
	}
	return tasks
}

func init() {
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteListCmd)
	remoteCmd.AddCommand(remoteRemoveCmd)
	rootCmd.AddCommand(remoteCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
}
