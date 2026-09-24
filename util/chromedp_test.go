// Copyright 2026 The cobbs.ai Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !skipCi
// +build !skipCi

package util

import (
	"context"
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

const defaultChromePath = `C:\Program Files\Google\Chrome\Application\chrome.exe`

func chromeExecutablePath(t *testing.T) string {
	t.Helper()

	execPath := os.Getenv("CHROME_PATH")
	if execPath != "" {
		validateChromeExecutable(t, execPath)
		return execPath
	}

	if runtime.GOOS == "windows" {
		execPath = defaultChromePath
		validateChromeExecutable(t, execPath)
		return execPath
	}

	for _, candidate := range []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser"} {
		if execPath, err := exec.LookPath(candidate); err == nil {
			return execPath
		}
	}

	t.Skip("Chrome is unavailable; set CHROME_PATH to run this browser integration test")
	return ""
}

func validateChromeExecutable(t *testing.T, execPath string) {
	t.Helper()

	info, err := os.Stat(execPath)
	if err != nil {
		t.Fatalf("Chrome executable %q is unavailable: %v; set CHROME_PATH to a valid browser path", execPath, err)
	}
	if info.IsDir() {
		t.Fatalf("Chrome executable path %q is a directory; set CHROME_PATH to the browser executable", execPath)
	}
}

func TestGetGoogleHomepage(t *testing.T) {
	execPath := chromeExecutablePath(t)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(execPath),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var title string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.google.com"),
		chromedp.Title(&title),
	); err != nil {
		t.Fatalf("chromedp run failed: %v", err)
	}

	log.Println("Page title:", title)
}
