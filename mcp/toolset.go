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

package mcp

import (
	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/baron929/cobbs.ai/toolauth"
	"github.com/baron929/cobbs.ai/tool"
)

// ToolSet holds everything needed to execute MCP and builtin tool calls
// during a model conversation: the open connections, the tool list for the AI,
// and the builtin tool registry.
type ToolSet struct {
	// Connections maps server name → open MCP connection (closed after the conversation).
	Connections map[string]*client.Client
	// Tools is the protocol-level tool list passed to the AI model.
	Tools []*protocol.Tool
	// BuiltinTools is the registry of server-side builtin tools.
	BuiltinTools     *tool.ToolRegistry
	WebSearchEnabled bool
	// Authorizer is mandatory for execution. A nil value is treated as deny.
	Authorizer toolauth.Authorizer
}
