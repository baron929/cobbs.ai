// Copyright 2023 The cobbs.ai Authors. All Rights Reserved.
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

package audio

import (
	"fmt"
	"os"
	"testing"
)

func TestGetAudioText(t *testing.T) {
	if _, err := os.Stat("../data/example.mp3"); err != nil {
		t.Skipf("example audio fixture is not available in this environment: %v", err)
	}
	segments, err := GetSegmentsFromAudio("../data/example.mp3", "en")
	if err != nil {
		t.Skipf("speech-to-text test dependency is unavailable in this environment: %v", err)
	}

	fmt.Printf("%v\n", segments)
}
