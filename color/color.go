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

package color

import (
	"fmt"
	"os"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Gray   = "\033[90m"
)

var useColor = true

func init() {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		useColor = false
	}
}

func Sprintf(code, format string, args ...interface{}) string {
	if !useColor {
		return fmt.Sprintf(format, args...)
	}
	return fmt.Sprintf("%s"+format+"%s", code, fmt.Sprintf(format, args...), Reset)
}

func Sprint(code, s string) string {
	if !useColor {
		return s
	}
	return code + s + Reset
}

func BoldS(s string) string  { return Sprint(Bold, s) }
func RedS(s string) string   { return Sprint(Red, s) }
func GreenS(s string) string { return Sprint(Green, s) }
func YellowS(s string) string { return Sprint(Yellow, s) }
func BlueS(s string) string  { return Sprint(Blue, s) }
func PurpleS(s string) string { return Sprint(Purple, s) }
func CyanS(s string) string  { return Sprint(Cyan, s) }
func GrayS(s string) string  { return Sprint(Gray, s) }
func DimS(s string) string   { return Sprint(Dim, s) }

func StatusIcon(status string) string {
	switch status {
	case "done":
		return GreenS("✓")
	case "cancelled":
		return RedS("✕")
	default:
		return YellowS("○")
	}
}

func PriorityLabel(p int) string {
	if p <= 0 {
		return ""
	}
	code := Cyan
	switch {
	case p >= 5:
		code = Red
	case p >= 4:
		code = Yellow
	case p >= 3:
		code = Purple
	case p >= 2:
		code = Green
	}
	return Sprint(code, fmt.Sprintf("P%d", p))
}

func IDLabel(id int) string {
	return CyanS(fmt.Sprintf("#%d", id))
}

func ProjectLabel(project string) string {
	if project == "" {
		return ""
	}
	return BlueS(fmt.Sprintf("[%s]", project))
}

func TagLabel(tag string) string {
	return PurpleS(tag)
}
