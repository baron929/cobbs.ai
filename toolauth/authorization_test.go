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

package toolauth

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/baron929/cobbs.ai/guard"
)

func TestDefaultAuthorizerAllowsReadOnlyBuiltin(t *testing.T) {
	a := NewDefaultAuthorizer()
	decision, err := a.Authorize(context.Background(), Request{Subject: "user-a", Tool: "local_file_read", Category: guard.CategoryRead, Arguments: map[string]interface{}{"path": "workspace/readme.md"}})
	if err != nil || !decision.Approved || decision.Effect != guard.EffectAllow {
		t.Fatalf("read-only tool should be allowed, decision=%+v err=%v", decision, err)
	}
}

func TestUnknownToolDeniedByDefault(t *testing.T) {
	a := NewDefaultAuthorizer()
	decision, err := a.Authorize(context.Background(), Request{Subject: "user-a", Tool: "future_tool", Category: guard.CategoryUnknown, Arguments: map[string]interface{}{}})
	if err == nil || decision.Approved || decision.Effect != guard.EffectAsk {
		t.Fatalf("unknown tool should fail closed as approval-required, decision=%+v err=%v", decision, err)
	}
}

func TestHighRiskToolRequiresApproval(t *testing.T) {
	a := NewDefaultAuthorizer()
	decision, err := a.Authorize(context.Background(), Request{Subject: "user-a", Tool: "shell", Category: guard.CategoryExec, Arguments: map[string]interface{}{"command": "echo safe"}})
	if err != ErrApprovalRequired || decision.Approved || !decision.ApprovalNeeded {
		t.Fatalf("shell should require approval, decision=%+v err=%v", decision, err)
	}
}

func TestApprovedActionIsBoundToScopeAndArguments(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	a := NewDefaultAuthorizer()
	a.Now = func() time.Time { return now }
	a.ApprovalLedger = &mockApprovalLedger{consumed: map[string]struct{}{}}
	args := map[string]interface{}{"command": "git status"}
	approval := &Approval{
		Token:         "one-time",
		Subject:       "user-a",
		Owner:         "owner-a",
		Store:         "store-a",
		Tool:          "shell",
		Resource:      "git status",
		ArgumentsHash: ArgumentsHash(args),
		ExpiresAt:     now.Add(time.Minute),
	}
	decision, err := a.Authorize(context.Background(), Request{Subject: "user-a", Owner: "owner-a", Store: "store-a", Tool: "shell", Category: guard.CategoryExec, Resource: "git status", Arguments: args, Approval: approval})
	if err != nil || !decision.Approved {
		t.Fatalf("matching approval should allow execution, decision=%+v err=%v", decision, err)
	}

	approval.Resource = "git push"
	_, err = a.Authorize(context.Background(), Request{Subject: "user-a", Owner: "owner-a", Store: "store-a", Tool: "shell", Category: guard.CategoryExec, Resource: "git status", Arguments: args, Approval: approval})
	if err != ErrInvalidApproval {
		t.Fatalf("scope mismatch should reject approval, got %v", err)
	}
}

func TestExpiredApprovalRejected(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	a := NewDefaultAuthorizer()
	a.Now = func() time.Time { return now }
	args := map[string]interface{}{"path": "workspace/file.txt"}
	approval := &Approval{Token: "expired", Subject: "user-a", Tool: "local_file_write", ArgumentsHash: ArgumentsHash(args), ExpiresAt: now.Add(-time.Second)}
	_, err := a.Authorize(context.Background(), Request{Subject: "user-a", Tool: "local_file_write", Category: guard.CategoryWrite, Arguments: args, Approval: approval})
	if err != ErrInvalidApproval {
		t.Fatalf("expired approval should be rejected, got %v", err)
	}
}

type mockApprovalLedger struct {
	mu       sync.Mutex
	consumed map[string]struct{}
}

func (m *mockApprovalLedger) Validate(req Request, approval Approval, now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := approval.Token + "|" + req.Subject + "|" + req.Owner + "|" + req.Store + "|" + req.Server + "|" + req.Tool + "|" + req.Resource
	if _, exists := m.consumed[key]; exists {
		return ErrInvalidApproval
	}
	return nil
}

func (m *mockApprovalLedger) Consume(req Request, approval Approval, now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := approval.Token + "|" + req.Subject + "|" + req.Owner + "|" + req.Store + "|" + req.Server + "|" + req.Tool + "|" + req.Resource
	if _, exists := m.consumed[key]; exists {
		return ErrInvalidApproval
	}
	m.consumed[key] = struct{}{}
	return nil
}

func TestApprovalRequiresDurableLedger(t *testing.T) {
	a := NewDefaultAuthorizer()
	args := map[string]interface{}{"command": "git status"}
	approval := Approval{Token: "transient", Subject: "user-a", Tool: "shell", ArgumentsHash: ArgumentsHash(args), ExpiresAt: time.Now().Add(time.Minute)}
	decision, err := a.Authorize(context.Background(), Request{Subject: "user-a", Tool: "shell", Category: guard.CategoryExec, Arguments: args, Approval: &approval})
	if err != ErrInvalidApproval || decision.Approved {
		t.Fatalf("approval without durable ledger must fail closed, decision=%+v err=%v", decision, err)
	}
}

func TestApprovalLedgerRejectsReplay(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	ledger := &mockApprovalLedger{consumed: map[string]struct{}{}}
	a := NewDefaultAuthorizer()
	a.Now = func() time.Time { return now }
	a.ApprovalLedger = ledger

	args := map[string]interface{}{"command": "git status"}
	approval := &Approval{
		Token:         "ledger-token",
		Subject:       "user-a",
		Owner:         "owner-a",
		Store:         "store-a",
		Tool:          "shell",
		Resource:      "git status",
		ArgumentsHash: ArgumentsHash(args),
		ExpiresAt:     now.Add(time.Minute),
	}

	decision, err := a.Authorize(context.Background(), Request{Subject: "user-a", Owner: "owner-a", Store: "store-a", Tool: "shell", Category: guard.CategoryExec, Resource: "git status", Arguments: args, Approval: approval})
	if err != nil || !decision.Approved {
		t.Fatalf("first approval should be accepted, decision=%+v err=%v", decision, err)
	}

	_, err = a.Authorize(context.Background(), Request{Subject: "user-a", Owner: "owner-a", Store: "store-a", Tool: "shell", Category: guard.CategoryExec, Resource: "git status", Arguments: args, Approval: approval})
	if err != ErrInvalidApproval {
		t.Fatalf("reused approval should fail replay protection, got %v", err)
	}
}

func TestApprovalLedgerAllowsOneConcurrentConsumption(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	ledger := &mockApprovalLedger{consumed: map[string]struct{}{}}
	a := NewDefaultAuthorizer()
	a.Now = func() time.Time { return now }
	a.ApprovalLedger = ledger
	args := map[string]interface{}{"command": "git status"}
	approval := &Approval{Token: "concurrent-token", Subject: "user-a", Owner: "owner-a", Store: "store-a", Tool: "shell", Resource: "git status", ArgumentsHash: ArgumentsHash(args), ExpiresAt: now.Add(time.Minute)}
	req := Request{Subject: "user-a", Owner: "owner-a", Store: "store-a", Tool: "shell", Category: guard.CategoryExec, Resource: "git status", Arguments: args, Approval: approval}

	var wg sync.WaitGroup
	results := make(chan error, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := a.Authorize(context.Background(), req)
			results <- err
		}()
	}
	wg.Wait()
	close(results)

	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("concurrent approval consumption should allow exactly one success, got %d", successes)
	}
}

func TestSafeErrorDoesNotExposeArguments(t *testing.T) {
	err := SafeError(ErrDenied)
	if err.Error() != "tool execution is not authorized" {
		t.Fatalf("unexpected safe error: %v", err)
	}
	if ArgumentsHash(map[string]interface{}{"secret": "value"}) == "" {
		t.Fatal("argument hash should be available without exposing argument content")
	}
}
