// Copyright 2025 The cobbs.ai Authors. All Rights Reserved.
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

package object

import (
	"strings"
	"testing"
)

func TestApprovalStateTransitions(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{name: "pending to approved", from: ApprovalStatePending, to: ApprovalStateApproved, want: true},
		{name: "pending to rejected", from: ApprovalStatePending, to: ApprovalStateRejected, want: true},
		{name: "approved to consumed", from: ApprovalStateApproved, to: ApprovalStateConsumed, want: true},
		{name: "approved to revoked", from: ApprovalStateApproved, to: ApprovalStateRevoked, want: true},
		{name: "pending to consumed", from: ApprovalStatePending, to: ApprovalStateConsumed, want: false},
		{name: "rejected to approved", from: ApprovalStateRejected, to: ApprovalStateApproved, want: false},
		{name: "consumed to approved", from: ApprovalStateConsumed, to: ApprovalStateApproved, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validApprovalTransition(test.from, test.to); got != test.want {
				t.Fatalf("validApprovalTransition(%q, %q) = %v, want %v", test.from, test.to, got, test.want)
			}
		})
	}
}

func TestApprovalTokenHashDoesNotContainRawToken(t *testing.T) {
	token := "approval-token-for-test-only"
	hash := hashApprovalToken(token)
	if hash == "" || hash == token || strings.Contains(hash, token) {
		t.Fatalf("token hash exposes raw token: %q", hash)
	}
}
