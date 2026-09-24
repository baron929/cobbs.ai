// Copyright 2026 The cobbs.ai Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0

package controllers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/baron929/cobbs.ai/audit"
	"github.com/baron929/cobbs.ai/object"
	"github.com/baron929/cobbs.ai/toolauth"
	"github.com/baron929/cobbs.ai/util"
)

type grantApprovalRequest struct {
	Store      string                 `json:"store"`
	Server     string                 `json:"server"`
	Tool       string                 `json:"tool"`
	Resource   string                 `json:"resource"`
	Arguments  map[string]interface{} `json:"arguments"`
	TTLSeconds int                    `json:"ttlSeconds"`
}

// GrantApproval creates a durable one-time approval for an exact tool call.
// The caller must be an authenticated global admin or a store admin granting
// only within their own owner scope.
func (c *ApiController) GrantApproval() {
	actor, ok := c.RequireSignedIn()
	if !ok {
		return
	}
	var request grantApprovalRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.ResponseError(err.Error())
		return
	}
	if request.Store == "" || request.Tool == "" || request.Arguments == nil {
		c.ResponseError("approval scope is incomplete")
		return
	}
	owner := "admin"
	if !c.IsGlobalAdmin() {
		if !c.IsStoreAdmin() {
			c.ResponseError("approval requires administrator privilege")
			return
		}
		owner = actor
	}
	store, err := object.GetStore(util.GetId(owner, request.Store))
	if err != nil || store == nil {
		c.ResponseError("approval store is not available in the authenticated owner scope")
		return
	}
	if !c.IsGlobalAdmin() && store.Owner != actor {
		c.ResponseError("approval requires authorization for the requested owner scope")
		return
	}
	if _, isolated := c.EnforceStoreIsolation(request.Store); !isolated {
		return
	}
	if request.TTLSeconds == 0 {
		request.TTLSeconds = 5 * 60
	}
	approval, err := object.GrantApproval(owner, actor, toolauth.Request{
		Subject: actor, Owner: owner, Store: request.Store,
		Server: request.Server, Tool: request.Tool, Resource: request.Resource, Arguments: request.Arguments,
	}, time.Duration(request.TTLSeconds)*time.Second)
	if err != nil {
		c.ResponseError("approval could not be created")
		return
	}
	audit.Record(audit.Event{Type: "approval_granted", Subject: actor, Owner: owner, Store: request.Store, Server: request.Server, Tool: request.Tool, ArgumentsLength: len(request.Arguments), ArgumentsHash: approval.ArgumentsHash, Outcome: "success"})
	c.ResponseOk(approval)
}

type updateApprovalStateRequest struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

// UpdateApprovalState permits an authorized owner administrator to reject or
// revoke a grant. Consumed approvals and other terminal states cannot move.
func (c *ApiController) UpdateApprovalState() {
	actor, ok := c.RequireSignedIn()
	if !ok {
		return
	}
	var request updateApprovalStateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil || request.ID == "" {
		c.ResponseError("invalid approval state request")
		return
	}
	if request.State != object.ApprovalStateRejected && request.State != object.ApprovalStateRevoked {
		c.ResponseError("unsupported approval state transition")
		return
	}
	record, err := object.GetApprovalRecord(request.ID)
	if err != nil || record == nil {
		c.ResponseError("approval not found")
		return
	}
	if !c.IsGlobalAdmin() && (!c.IsStoreAdmin() || record.Owner != actor) {
		c.ResponseError("approval requires authorization for the requested owner scope")
		return
	}
	if _, isolated := c.EnforceStoreIsolation(record.Store); !isolated {
		return
	}
	success, err := object.UpdateApprovalRecordState(request.ID, request.State)
	if err != nil {
		c.ResponseError("approval state could not be changed")
		return
	}
	audit.Record(audit.Event{Type: "approval_" + strings.ToLower(request.State), Subject: actor, Owner: record.Owner, Store: record.Store, Server: record.Server, Tool: record.Tool, ArgumentsHash: record.ArgumentsHash, Outcome: "success"})
	c.ResponseOk(success)
}
