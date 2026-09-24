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

package object

import (
	"strings"

	"github.com/baron929/cobbs.ai/mcp"
	"github.com/baron929/cobbs.ai/model"
	"github.com/baron929/cobbs.ai/toolauth"
	"github.com/baron929/cobbs.ai/tool"
	"github.com/baron929/cobbs.ai/util"
)

func buildToolSetForBuiltinTool(toolName, user, origin, lang string) (*mcp.ToolSet, error) {
	newAuthorizer := func() *toolauth.PolicyAuthorizer {
		a := toolauth.NewDefaultAuthorizer()
		a.ApprovalLedger = NewApprovalLedger()
		return a
	}

	if toolName == "" {
		return nil, nil
	}

	id := util.GetIdFromOwnerAndName("admin", toolName)
	t, err := GetTool(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	tp, err := tool.New(getToolConfig(t), lang)
	if err != nil {
		return nil, err
	}

	reg := tool.NewToolRegistry()
	for _, t := range tp.BuiltinTools() {
		wrapped := wrapSnapshotBuiltin("admin", t)
		wrapped = wrapGeneratedResourceBuiltin(wrapped, "admin", user, origin)
		reg.RegisterTool(wrapped)
	}

	allTools := reg.GetToolsAsProtocolTools()
	if len(allTools) == 0 {
		return nil, nil
	}

	return &mcp.ToolSet{
		Tools:        allTools,
		BuiltinTools: reg,
		Authorizer:   newAuthorizer(),
	}, nil
}

func GetAnswerWithTool(modelProviderName, toolName, question, user, origin, lang string) (string, *model.ModelResult, error) {
	modelProvider, modelProviderObj, err := GetModelProviderFromContext("admin", modelProviderName, lang)
	if err != nil {
		return "", nil, err
	}

	mcpToolSet, err := buildToolSetForBuiltinTool(toolName, user, origin, lang)
	if err != nil {
		return "", nil, err
	}

	prompt := "You are an expert in your field and you specialize in using your knowledge to answer or solve people's problems."
	history := []*model.RawMessage{}
	knowledge := []*model.RawMessage{}

	var writer MyWriter
	var modelResult *model.ModelResult

	if mcpToolSet != nil {
		messages := &model.ToolMessages{
			Messages:  []*model.RawMessage{},
			ToolCalls: nil,
		}
		toolSession := &model.ToolSession{
			McpToolSet:   mcpToolSet,
			ToolMessages: messages,
			IsVision:     modelProvider != nil && model.IsVisionModel(modelProvider.SubType),
			Subject:      user,
			Owner:        "admin",
			Store:        toolName,
		}
		modelResult, err = model.QueryTextWithTools(modelProviderObj, question, &writer, history, prompt, knowledge, toolSession, lang)
	} else {
		modelResult, err = modelProviderObj.QueryText(question, &writer, history, prompt, knowledge, nil, lang)
	}
	if err != nil {
		return "", nil, err
	}

	res := writer.String()
	res = strings.Trim(res, "\"")
	return res, modelResult, nil
}
