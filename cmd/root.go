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

	"github.com/U-SITE-SAS-BIC/utasker/branding"
	"github.com/U-SITE-SAS-BIC/utasker/db"
	"github.com/spf13/cobra"
)

var (
	projectFlag string
	allFlag     bool
	statusFlag  string
	Version     = "dev"
	Commit      = "none"
	Date        = "unknown"
)

func getProject() string {
	if projectFlag != "" {
		return projectFlag
	}
	p, err := db.GetProjectFromDir()
	if err == nil && p != "" {
		return p
	}
	return ""
}

func helpFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Print(branding.TagLine())
	}
	cmd.Parent().Help()
}

var rootCmd = &cobra.Command{
	Use:   "utasker",
	Short: "Offline task manager · by U-SITE",
	Long: `utasker is an offline-first CLI task manager that connects tasks
to your project directories. No cloud, no database — just JSON.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func helpTemplate() string {
	tpl := rootCmd.HelpTemplate()
	brandBlock := branding.SmallLogo() + "\n"
	tpl = brandBlock + tpl
	tpl = strings.ReplaceAll(tpl, "{{.Short}}", branding.CyanS("utasker")+" · offline task manager · by "+branding.BlueS("U-SITE"))
	tpl = strings.ReplaceAll(tpl, "{{.UseLine}}", "{{.UseLine}}")
	return tpl
}

func Execute() {
	rootCmd.Version = Version
	rootCmd.SetHelpTemplate(helpTemplate())
	rootCmd.SetVersionTemplate(branding.FullBanner(Version, Commit, Date))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&projectFlag, "project", "", "Project name (overrides .task-project file)")
}
