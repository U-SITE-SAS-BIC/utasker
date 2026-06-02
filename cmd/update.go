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
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/U-SITE-SAS-BIC/utasker/branding"
	"github.com/spf13/cobra"
)

const repo = "U-SITE-SAS-BIC/utasker"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Self-update to the latest version",
	Long: `Downloads and replaces the current binary with the latest GitHub release.
Requires internet access.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentPath, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), "cannot find current binary:", err)
			os.Exit(1)
		}

		info, err := os.Stat(currentPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		fmt.Printf("  %s Checking for updates...\n", branding.CyanS("↻"))

		release, err := fetchLatestRelease()
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		fmt.Printf("  %s Latest version: %s\n", branding.GreenS("✓"), branding.CyanS(release.Tag))

		assetName := assetName()
		downloadURL := ""
		for _, a := range release.Assets {
			if a.Name == assetName {
				downloadURL = a.DownloadURL
				break
			}
		}

		if downloadURL == "" {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), "no asset found for", assetName)
			fmt.Printf("  Available: %s\n", branding.YellowS("create a GitHub release first"))
			os.Exit(1)
		}

		fmt.Printf("  %s Downloading %s...\n", branding.CyanS("↓"), assetName)

		tmpDir, err := os.MkdirTemp("", "utasker-update")
		if err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)

		archivePath := filepath.Join(tmpDir, assetName)
		if err := downloadFile(downloadURL, archivePath); err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		fmt.Printf("  %s Extracting...\n", branding.CyanS("◉"))

		extractPath := filepath.Join(tmpDir, "utasker")
		if runtime.GOOS == "windows" {
			extractPath += ".exe"
		}

		if err := extractBinary(archivePath, extractPath); err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		if err := os.Chmod(extractPath, info.Mode()); err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		backupPath := currentPath + ".bak"
		if err := os.Rename(currentPath, backupPath); err != nil {
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		if err := copyFile(extractPath, currentPath); err != nil {
			os.Rename(backupPath, currentPath)
			fmt.Fprintln(os.Stderr, branding.RedS("Error:"), err)
			os.Exit(1)
		}

		os.Remove(backupPath)

		fmt.Printf("\n  %s Updated to %s\n", branding.GreenS("✓"), branding.CyanS(release.Tag))
		fmt.Printf("  %s Restart utasker to use the new version\n", branding.GrayS("→"))
	},
}

type ghRelease struct {
	Tag    string     `json:"tag_name"`
	Assets []ghAsset  `json:"assets"`
}

type ghAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

func fetchLatestRelease() (ghRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ghRelease{}, fmt.Errorf("cannot reach GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return ghRelease{}, fmt.Errorf("no releases found — create a tag and push it")
	}
	if resp.StatusCode != 200 {
		return ghRelease{}, fmt.Errorf("GitHub returned %d", resp.StatusCode)
	}

	var r ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return ghRelease{}, fmt.Errorf("bad response: %w", err)
	}
	return r, nil
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if strings.HasSuffix(archivePath, ".tar.gz") {
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if hdr.Name == "utasker" || hdr.Name == "utasker.exe" {
				out, err := os.Create(dest)
				if err != nil {
					return err
				}
				defer out.Close()
				_, err = io.Copy(out, tr)
				return err
			}
		}
		return fmt.Errorf("binary not found in archive")
	}
	return fmt.Errorf("unsupported archive format: %s", archivePath)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func assetName() string {
	osMap := map[string]string{
		"darwin":  "macOS",
		"linux":   "linux",
		"windows": "windows",
	}
	goos := osMap[runtime.GOOS]
	if goos == "" {
		goos = runtime.GOOS
	}

	archMap := map[string]string{
		"amd64": "x86_64",
		"arm64": "aarch64",
	}
	arch := archMap[runtime.GOARCH]
	if arch == "" {
		arch = runtime.GOARCH
	}

	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}

	return fmt.Sprintf("utasker_%s_%s_%s.%s", "latest", goos, arch, ext)
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
