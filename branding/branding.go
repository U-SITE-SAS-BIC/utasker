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

package branding

import (
	"fmt"
	"os"
	"strings"
)

var useColor = true

func init() {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		useColor = false
	}
}

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Cyan    = "\033[36m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	Yellow  = "\033[33m"
	Green   = "\033[32m"
	Gray    = "\033[90m"
	White   = "\033[97m"
	BgBlue  = "\033[44m"
)

func colorize(code, s string) string {
	if !useColor {
		return s
	}
	return code + s + Reset
}

func CyanS(s string) string   { return colorize(Cyan, s) }
func BlueS(s string) string   { return colorize(Blue, s) }
func PurpleS(s string) string { return colorize(Purple, s) }
func YellowS(s string) string { return colorize(Yellow, s) }
func GreenS(s string) string  { return colorize(Green, s) }
func GrayS(s string) string   { return colorize(Gray, s) }
func BoldS(s string) string   { return colorize(Bold, s) }
func WhiteS(s string) string  { return colorize(White, s) }

var logoLines = []string{
	"██╗   ██╗████████╗ █████╗ ███████╗██╗  ██╗███████╗██████╗ ",
	"██║   ██║╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗",
	"██║   ██║   ██║   ███████║███████╗█████╔╝ █████╗  ██████╔╝",
	"██║   ██║   ██║   ██╔══██║╚════██║██╔═██╗ ██╔══╝  ██╔══██╗",
	"╚██████╔╝   ██║   ██║  ██║███████║██║  ██╗███████╗██║  ██║",
	" ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝",
}

var smallLogo = ` _   _ _____ ____  _  __
| | | |_   _/ ___|| |/ / ___  __ _ ___
| | | | | | \___ \| ' / / __|/ _` + "`" + ` / __|
| |_| | | |  ___) | . \ \__ \ (_| \__ \
 \___/  |_| |____/|_|\_\|___/\__,_|___/`

func Logo() string {
	primary := Cyan
	secondary := Blue
	accent := Purple

	var b strings.Builder
	for i, line := range logoLines {
		mid := len(line) / 2
		if !useColor {
			b.WriteString("  " + line + "\n")
			continue
		}
		if i%2 == 0 {
			b.WriteString("  " + colorize(primary, line[:mid]) + colorize(secondary, line[mid:]) + "\n")
		} else {
			b.WriteString("  " + colorize(secondary, line[:mid]) + colorize(primary, line[mid:]) + "\n")
		}
	}
	b.WriteString("\n")
	b.WriteString("  " + colorize(accent, BoldS("utasker")) + colorize(Gray, " · offline task manager") + "\n")
	b.WriteString("  " + colorize(Gray, "by ") + colorize(Blue, BoldS("U-SITE")) + colorize(Gray, " · ") + colorize(Gray, "https://u-site.app") + "\n")
	return b.String()
}

func SmallLogo() string {
	primary := Cyan
	secondary := Blue

	var b strings.Builder
	lines := strings.Split(smallLogo, "\n")
	for i, line := range lines {
		mid := len(line) * 3 / 5
		if !useColor {
			b.WriteString("  " + line + "\n")
			continue
		}
		if i%2 == 0 {
			b.WriteString("  " + colorize(primary, line[:mid]) + colorize(secondary, line[mid:]) + "\n")
		} else {
			b.WriteString("  " + colorize(secondary, line[:mid]) + colorize(primary, line[mid:]) + "\n")
		}
	}
	return b.String()
}

func TagLine() string {
	logo := SmallLogo()
	msg := colorize(Cyan, BoldS("utasker"))
	by := colorize(Blue, BoldS("U-SITE"))
	return fmt.Sprintf("%s\n  %s · %s\n", logo, msg, by)
}

func FullBanner(version, commit, date string) string {
	var b strings.Builder
	b.WriteString(Logo())
	b.WriteString("\n")
	b.WriteString("  " + colorize(Gray, "Version:  ") + colorize(Cyan, version) + "\n")
	b.WriteString("  " + colorize(Gray, "Commit:   ") + colorize(Cyan, commit) + "\n")
	b.WriteString("  " + colorize(Gray, "Built:    ") + colorize(Cyan, date) + "\n")
	b.WriteString("\n")
	return b.String()
}
