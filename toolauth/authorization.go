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

// Package toolauth is the mandatory policy boundary for agent tool execution.
// It has no dependency on model, MCP, HTTP, or persistence packages so every
// execution path can use the same authorization contract.
package toolauth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/baron929/cobbs.ai/guard"
)

var (
	ErrDenied           = errors.New("tool execution denied")
	ErrApprovalRequired = errors.New("tool execution requires approval")
	ErrApprovalDenied   = errors.New("tool approval denied")
	ErrInvalidApproval  = errors.New("invalid tool approval")
)

// Request is the security context for one attempted tool execution.
type Request struct {
	Subject   string
	Owner     string
	Store     string
	Device    string
	Server    string
	Tool      string
	Category  string
	Resource  string
	Arguments map[string]interface{}
	Approval  *Approval
}

// Approval binds approval to the exact caller, resource, and argument set.
type Approval struct {
	Token         string    `json:"token"`
	Subject       string    `json:"subject"`
	Owner         string    `json:"owner"`
	Store         string    `json:"store"`
	Device        string    `json:"device"`
	Server        string    `json:"server"`
	Tool          string    `json:"tool"`
	Resource      string    `json:"resource"`
	ArgumentsHash string    `json:"argumentsHash"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

// Decision is the policy result plus whether execution is permitted now.
type Decision struct {
	Effect         guard.Effect
	Reason         string
	Rule           string
	Approved       bool
	ApprovalNeeded bool
}

// Approver is an optional host-provided approval workflow. Returning false or
// an error never permits execution.
type Approver interface {
	RequestApproval(context.Context, Request, guard.Decision) (bool, error)
}

// ApprovalLedger enforces durable, one-time approval consumption.
// Implementations must reject reused approvals even if the approval token is
// otherwise valid and not expired.
type ApprovalLedger interface {
	Validate(req Request, approval Approval, now time.Time) error
	Consume(req Request, approval Approval, now time.Time) error
}

// Authorizer is the only interface execution paths need to depend on.
type Authorizer interface {
	Authorize(context.Context, Request) (Decision, error)
}

// PolicyAuthorizer evaluates the existing guard engine and an optional
// approval workflow. It never executes tools.
type PolicyAuthorizer struct {
	Guard          guard.Guard
	Approver       Approver
	ApprovalLedger ApprovalLedger
	Now            func() time.Time
}

// NewDefaultAuthorizer allows only explicitly classified read-only tools. All
// network, write, exec, sensitive, and unknown categories require approval;
// without an approver they fail closed.
func NewDefaultAuthorizer() *PolicyAuthorizer {
	g, err := guard.NewCasbinGuard(guard.Policy{
		Default: guard.EffectAsk,
		Rules: []guard.Rule{{
			Name:     "builtin-read-only",
			Category: guard.CategoryRead,
			Effect:   guard.EffectAllow,
			Priority: 100,
		}},
	})
	if err != nil {
		return NewDenyAuthorizer()
	}
	return &PolicyAuthorizer{Guard: g, Now: time.Now}
}

// NewDenyAuthorizer is used when an execution path forgot to provide policy.
func NewDenyAuthorizer() *PolicyAuthorizer {
	g, _ := guard.NewCasbinGuard(guard.Policy{Default: guard.EffectDeny})
	return &PolicyAuthorizer{Guard: g, Now: time.Now}
}

func (a *PolicyAuthorizer) Authorize(ctx context.Context, req Request) (Decision, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if req.Subject == "" || req.Tool == "" || req.Arguments == nil {
		return Decision{Effect: guard.EffectDeny, Reason: "missing tool security context"}, ErrDenied
	}
	if a == nil || a.Guard == nil {
		return Decision{Effect: guard.EffectDeny, Reason: "no authorization policy"}, ErrDenied
	}

	policyDecision, err := a.Guard.Check(ctx, guard.Request{
		Subject:  req.Subject,
		Tool:     req.Tool,
		Category: req.Category,
		Resource: req.Resource,
	})
	if err != nil {
		return Decision{Effect: guard.EffectDeny, Reason: "authorization policy error"}, ErrDenied
	}

	decision := Decision{
		Effect: policyDecision.Effect,
		Reason: policyDecision.Reason,
		Rule:   policyDecision.Rule,
	}
	switch policyDecision.Effect {
	case guard.EffectAllow:
		decision.Approved = true
		return decision, nil
	case guard.EffectDeny:
		return decision, ErrDenied
	case guard.EffectAsk:
		decision.ApprovalNeeded = true
		if req.Approval != nil {
			if a.ApprovalLedger == nil {
				return decision, ErrInvalidApproval
			}
			if err := VerifyApproval(req, *req.Approval, a.now()); err != nil {
				return decision, err
			}
			if err := a.ApprovalLedger.Validate(req, *req.Approval, a.now()); err != nil {
				return decision, err
			}
			if err := a.ApprovalLedger.Consume(req, *req.Approval, a.now()); err != nil {
				return decision, err
			}
			decision.Approved = true
			return decision, nil
		}
		if a.Approver == nil {
			return decision, ErrApprovalRequired
		}
		approved, err := a.Approver.RequestApproval(ctx, req, policyDecision)
		if err != nil || !approved {
			return decision, ErrApprovalDenied
		}
		// The legacy approver contract returns only a boolean and cannot provide
		// a durable one-time grant. Never execute from that transient result.
		return decision, ErrApprovalRequired
	default:
		return Decision{Effect: guard.EffectDeny, Reason: "unknown authorization effect"}, ErrDenied
	}
}

func (a *PolicyAuthorizer) now() time.Time {
	if a != nil && a.Now != nil {
		return a.Now()
	}
	return time.Now()
}

// VerifyApproval rejects expired, replayable-by-context, or scope-mismatched
// approvals. Token consumption belongs to the persistence-backed approver.
func VerifyApproval(req Request, approval Approval, now time.Time) error {
	if approval.Token == "" || approval.ExpiresAt.IsZero() || !approval.ExpiresAt.After(now) {
		return ErrInvalidApproval
	}
	if approval.Subject != req.Subject || approval.Owner != req.Owner || approval.Store != req.Store || approval.Device != req.Device || approval.Server != req.Server || approval.Tool != req.Tool || approval.Resource != req.Resource {
		return ErrInvalidApproval
	}
	if approval.ArgumentsHash != ArgumentsHash(req.Arguments) {
		return ErrInvalidApproval
	}
	return nil
}

// ArgumentsHash provides a stable non-secret binding for audit and approval.
func ArgumentsHash(arguments map[string]interface{}) string {
	encoded, err := json.Marshal(arguments)
	if err != nil {
		return ""
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}

// Classify returns the least-privileged broad category for a known builtin or
// an unknown category for MCP/future tools. Unknown tools require explicit
// policy and cannot inherit the read-only default.
func Classify(server, tool string) string {
	if server != "" {
		return guard.CategoryUnknown
	}
	switch tool {
	case "time", "local_file_scan", "local_file_read", "web_search":
		return guard.CategoryRead
	case "web_fetch", "web_browser", "browser_screenshot", "video_download":
		return guard.CategoryNetwork
	case "shell", "browser_evaluate", "browser_click", "browser_use_open", "gui":
		return guard.CategoryExec
	case "local_file_write", "local_file_move", "office", "browser_use_click":
		return guard.CategoryWrite
	default:
		if strings.Contains(tool, "secret") || strings.Contains(tool, "credential") {
			return guard.CategorySensitive
		}
		return guard.CategoryUnknown
	}
}

// Resource derives a policy-matchable resource without putting raw arguments
// into audit records. The raw value is only passed to the guard matcher.
func Resource(arguments map[string]interface{}) string {
	for _, key := range []string{"command", "path", "source", "target", "url"} {
		if value, ok := arguments[key].(string); ok && strings.TrimSpace(value) != "" {
			if key == "url" {
				if parsed, err := url.Parse(value); err == nil && parsed.Host != "" {
					return parsed.Host
				}
			}
			return value
		}
	}
	return ""
}

// SafeError maps authorization failures to stable messages without arguments,
// secrets, policy internals, or resource values.
func SafeError(err error) error {
	switch {
	case errors.Is(err, ErrApprovalRequired):
		return fmt.Errorf("tool execution requires approval")
	case errors.Is(err, ErrApprovalDenied):
		return fmt.Errorf("tool approval denied")
	case errors.Is(err, ErrInvalidApproval):
		return fmt.Errorf("tool approval is invalid or expired")
	default:
		return fmt.Errorf("tool execution is not authorized")
	}
}
