// Copyright 2026 The cobbs.ai Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0

package object

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/baron929/cobbs.ai/toolauth"
	"github.com/baron929/cobbs.ai/util"
	"xorm.io/core"
)

const (
	ApprovalStatePending  = "Pending"
	ApprovalStateApproved = "Approved"
	ApprovalStateRejected = "Rejected"
	ApprovalStateConsumed = "Consumed"
	ApprovalStateExpired  = "Expired"
	ApprovalStateRevoked  = "Revoked"
)

// ApprovalRecord is the durable ledger row for a single one-time approval.
// Raw tokens are never stored; TokenHash is the lookup and persistence value.
type ApprovalRecord struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`

	TokenHash     string `xorm:"varchar(128) unique notnull" json:"-"`
	Subject       string `xorm:"varchar(100) index notnull" json:"subject"`
	Store         string `xorm:"varchar(100) index notnull" json:"store"`
	Device        string `xorm:"varchar(100) index" json:"device"`
	OwnerScope    string `xorm:"varchar(100) index notnull" json:"ownerScope"`
	Approver      string `xorm:"varchar(100) index" json:"approver"`
	Server        string `xorm:"varchar(100) index" json:"server"`
	Tool          string `xorm:"varchar(200) index" json:"tool"`
	Resource      string `xorm:"varchar(500) index" json:"resource"`
	ArgumentsHash string `xorm:"varchar(128) index notnull" json:"argumentsHash"`
	ExpiresAt     string `xorm:"varchar(100) index notnull" json:"expiresAt"`
	State         string `xorm:"varchar(50) index notnull" json:"state"`
	UsedAt        string `xorm:"varchar(100)" json:"usedAt"`
	Reason        string `xorm:"varchar(200)" json:"reason"`
}

func (a *ApprovalRecord) GetId() string { return util.GetId(a.Owner, a.Name) }

func (a *ApprovalRecord) normalize() error {
	if a == nil || a.Owner == "" || a.Name == "" || a.TokenHash == "" {
		return fmt.Errorf("approval owner, name and token hash are required")
	}
	if a.State == "" {
		a.State = ApprovalStatePending
	}
	if !validApprovalState(a.State) {
		return fmt.Errorf("invalid approval state %q", a.State)
	}
	return nil
}

func validApprovalState(state string) bool {
	switch state {
	case ApprovalStatePending, ApprovalStateApproved, ApprovalStateRejected, ApprovalStateConsumed, ApprovalStateExpired, ApprovalStateRevoked:
		return true
	default:
		return false
	}
}

func validApprovalTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case ApprovalStatePending:
		return to == ApprovalStateApproved || to == ApprovalStateRejected || to == ApprovalStateExpired || to == ApprovalStateRevoked
	case ApprovalStateApproved:
		return to == ApprovalStateConsumed || to == ApprovalStateExpired || to == ApprovalStateRevoked
	default:
		return false
	}
}

func hashApprovalToken(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}

func getApprovalRecordByToken(token string) (*ApprovalRecord, error) {
	if token == "" {
		return nil, nil
	}
	var records []*ApprovalRecord
	if err := adapter.engine.Where("token_hash = ?", hashApprovalToken(token)).Find(&records); err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return records[0], nil
}

func AddApprovalRecord(record *ApprovalRecord) (bool, error) {
	if err := record.normalize(); err != nil {
		return false, err
	}
	if record.CreatedTime == "" {
		record.CreatedTime = util.GetCurrentTimeWithMilli()
	}
	if record.UpdatedTime == "" {
		record.UpdatedTime = record.CreatedTime
	}
	affected, err := adapter.engine.Insert(record)
	return affected != 0, err
}

func GetApprovalRecord(id string) (*ApprovalRecord, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return nil, err
	}
	record := &ApprovalRecord{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(record)
	if err != nil || !existed {
		return nil, err
	}
	return record, nil
}

func UpdateApprovalRecordState(id, state string) (bool, error) {
	if !validApprovalState(state) {
		return false, fmt.Errorf("invalid approval state %q", state)
	}
	record, err := GetApprovalRecord(id)
	if err != nil || record == nil {
		return false, err
	}
	if !validApprovalTransition(record.State, state) {
		return false, fmt.Errorf("invalid approval state transition from %q to %q", record.State, state)
	}
	updatedAt := util.GetCurrentTimeWithMilli()
	updates := &ApprovalRecord{State: state, UpdatedTime: updatedAt}
	if state == ApprovalStateConsumed {
		updates.UsedAt = updatedAt
	}
	affected, err := adapter.engine.ID(core.PK{record.Owner, record.Name}).Cols("state", "updated_time", "used_at").Update(updates)
	return affected == 1, err
}

func DeleteApprovalRecord(record *ApprovalRecord) (bool, error) {
	if record == nil {
		return false, nil
	}
	affected, err := adapter.engine.ID(core.PK{record.Owner, record.Name}).Delete(&ApprovalRecord{})
	return affected != 0, err
}

type ApprovalLedger struct{}

func NewApprovalLedger() *ApprovalLedger { return &ApprovalLedger{} }

func (l *ApprovalLedger) Validate(req toolauth.Request, approval toolauth.Approval, now time.Time) error {
	if err := toolauth.VerifyApproval(req, approval, now); err != nil {
		return err
	}
	record, err := getApprovalRecordByToken(approval.Token)
	if err != nil {
		return err
	}
	if record == nil || record.State != ApprovalStateApproved || record.ExpiresAt == "" {
		return toolauth.ErrInvalidApproval
	}
	expiresAt, err := time.Parse(util.TimeFormat, record.ExpiresAt)
	if err != nil || !expiresAt.After(now) {
		return toolauth.ErrInvalidApproval
	}
	if record.TokenHash != hashApprovalToken(approval.Token) || record.Owner != req.Owner || record.OwnerScope != req.Owner || record.Subject != req.Subject || record.Store != req.Store || record.Device != req.Device || record.Server != req.Server || record.Tool != req.Tool || record.Resource != req.Resource || record.ArgumentsHash != approval.ArgumentsHash {
		return toolauth.ErrInvalidApproval
	}
	return nil
}

// Consume reserves the approval before tool execution. The conditional update
// makes concurrent callers mutually exclusive and consumes nothing on DB error.
func (l *ApprovalLedger) Consume(req toolauth.Request, approval toolauth.Approval, now time.Time) error {
	if err := l.Validate(req, approval, now); err != nil {
		return err
	}
	record, err := getApprovalRecordByToken(approval.Token)
	if err != nil {
		return err
	}
	if record == nil {
		return toolauth.ErrInvalidApproval
	}
	updatedAt := util.FormatTimeForCompare(now)
	updates := &ApprovalRecord{State: ApprovalStateConsumed, UpdatedTime: updatedAt, UsedAt: updatedAt}
	affected, err := adapter.engine.
		Where("owner = ? AND name = ? AND token_hash = ? AND owner_scope = ? AND subject = ? AND store = ? AND device = ? AND server = ? AND tool = ? AND resource = ? AND arguments_hash = ? AND state = ? AND expires_at > ?",
			record.Owner, record.Name, hashApprovalToken(approval.Token), req.Owner, req.Subject, req.Store, req.Device, req.Server, req.Tool, req.Resource, approval.ArgumentsHash, ApprovalStateApproved, util.FormatTimeForCompare(now)).
		Cols("state", "updated_time", "used_at").Update(updates)
	if err != nil {
		return err
	}
	if affected != 1 {
		return toolauth.ErrInvalidApproval
	}
	return nil
}

// GrantApproval creates an approved one-time grant and returns the raw token
// only to the authenticated approval workflow caller. It is never persisted.
func GrantApproval(owner, approver string, req toolauth.Request, ttl time.Duration) (*toolauth.Approval, error) {
	if owner == "" || approver == "" || req.Owner != owner || req.Subject == "" || req.Store == "" || req.Tool == "" || req.Arguments == nil {
		return nil, fmt.Errorf("invalid approval scope")
	}
	if ttl <= 0 || ttl > 15*time.Minute {
		return nil, fmt.Errorf("approval expiry must be between 1 second and 15 minutes")
	}
	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(rawToken)
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	approval := &toolauth.Approval{Token: token, Subject: req.Subject, Owner: req.Owner, Store: req.Store, Device: req.Device, Server: req.Server, Tool: req.Tool, Resource: req.Resource, ArgumentsHash: toolauth.ArgumentsHash(req.Arguments), ExpiresAt: expiresAt}
	record := &ApprovalRecord{Owner: owner, Name: util.GenerateId(), CreatedTime: util.FormatTimeForCompare(now), UpdatedTime: util.FormatTimeForCompare(now), TokenHash: hashApprovalToken(token), Subject: req.Subject, Store: req.Store, Device: req.Device, OwnerScope: req.Owner, Approver: approver, Server: req.Server, Tool: req.Tool, Resource: req.Resource, ArgumentsHash: approval.ArgumentsHash, ExpiresAt: util.FormatTimeForCompare(expiresAt), State: ApprovalStateApproved}
	if _, err := AddApprovalRecord(record); err != nil {
		return nil, err
	}
	return approval, nil
}
