// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
package elmgen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRPCWithComments(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.service;
		// service comment 0
		service HelloWorld {
			// method comment 0
			// method comment 1
			rpc Hello2(HelloReq) returns (HelloResp);
			// method comment 2
			rpc Hello1(HelloReq) returns (HelloResp);
			rpc Hello3(HelloReq) returns (HelloResp);
		}

		service Before {
			rpc Ignored(HelloReq) returns (HelloResp);
		}

		message HelloReq {
			string subject = 1;
		}
		message HelloResp {
			string text = 1;
		}
	`)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.Services, 2)
	assert.Len(t, elm.Services[1].Methods, 3)
	for i, rpc := range elm.Services[1].Methods {
		assert.Equal(t, protoreflect.FullName("test.service.HelloWorld"), rpc.Service)
		assert.Equal(t, protoreflect.Name(fmt.Sprintf("Hello%d", i+1)), rpc.Method)
		assert.Equal(t, fmt.Sprintf("twirpHelloWorld_Hello%d", i+1), rpc.ID.ID)
		assert.Equal(t, "HelloReq", rpc.In.ID)
		assert.Equal(t, "HelloResp", rpc.Out.ID)
		assert.Equal(t, "encodeHelloReq", rpc.In.Encoder.String())
		assert.Equal(t, "decodeHelloResp", rpc.Out.Decoder.String())
		assert.False(t, rpc.InStreaming)
		assert.False(t, rpc.OutStreaming)
	}
	// It would be great to test the generated code but can't process a Cmd

	// Check comments
	content := string(testFileContents["Test/ServiceTwirp.elm"])
	assert.True(t, strings.Contains(content, "service comment 0"))
	assert.True(t, strings.Contains(content, "method comment 0"))
	assert.True(t, strings.Contains(content, "method comment 1"))
	assert.True(t, strings.Contains(content, "method comment 2"))
}
