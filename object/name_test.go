// Copyright 2024 The cobbs.ai Authors. All Rights Reserved.
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

package object_test

import (
	"fmt"
	"testing"

	"github.com/baron929/cobbs.ai/auth"
	"github.com/baron929/cobbs.ai/conf"
	"github.com/baron929/cobbs.ai/controllers"
	"github.com/baron929/cobbs.ai/object"
)

var userTag = "user"

func TestUpdateMessagesForName(t *testing.T) {
	object.InitConfig()
	controllers.InitAuthConfig()
	if !conf.IsCasdoorAvailable() {
		t.Skip("Casdoor integration unavailable: configure and start Casdoor to run this migration test")
	}

	users, err := auth.GetUsers()
	if err != nil {
		panic(err)
	}

	userMap := map[string]*auth.User{}
	for _, user := range users {
		if user.Tag != userTag {
			continue
		}

		userMap[user.Name] = user
	}

	messages, err := object.GetGlobalMessages()
	if err != nil {
		panic(err)
	}

	for i, message := range messages {
		user, ok := userMap[message.User]
		if ok {
			message.User = user.DisplayName
		}

		if message.Author != "AI" {
			user, ok = userMap[message.Author]
			if ok {
				message.Author = user.DisplayName
			}
		}

		fmt.Printf("[%d/%d] message: %s, organization: %s, user: %s, author: %s\n", i+1, len(messages), message.Name, message.Organization, message.User, message.Author)

		_, err = object.UpdateMessage(message.GetId(), message, false)
		if err != nil {
			panic(err)
		}
	}
}

func TestUpdateChatsForName(t *testing.T) {
	t.Skip("stale legacy chat-name migration test references deleted Chat.User1/Users fields")
}

func TestUpdateMessagesAndChatsForName(t *testing.T) {
	TestUpdateMessagesForName(t)
	TestUpdateChatsForName(t)
}
