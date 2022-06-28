package elmgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRPC(t *testing.T) {
	elm := TestConfig.testModule(t, `
		syntax = "proto3";
		package test.service;
		service HelloWorld {
			rpc Hello1(HelloReq) returns (HelloResp);
			rpc Hello2(HelloReq) returns (HelloResp);
			rpc Hello3(HelloReq) returns (HelloResp);
		}

		message HelloReq {
			string subject = 1;
		}
		message HelloResp {
			string text = 1;
		}
	`)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.RPCs, 3)
	for i, rpc := range elm.RPCs {
		assert.Equal(t, protoreflect.FullName("test.service.HelloWorld"), rpc.Service)
		assert.Equal(t, protoreflect.Name(fmt.Sprintf("Hello%d", i+1)), rpc.Method)
		assert.Equal(t, fmt.Sprintf("hello%d", i+1), rpc.MethodID)
		assert.Equal(t, "HelloReq", rpc.In)
		assert.Equal(t, "HelloResp", rpc.Out)
		assert.Equal(t, "helloReqEncoder", rpc.InEncoder)
		assert.Equal(t, "helloRespDecoder", rpc.OutDecoder)
		assert.False(t, rpc.InStreaming)
		assert.False(t, rpc.OutStreaming)
	}
	// It would be great to test the generated code but can't process a Cmd
}

func TestRPCPrefixes(t *testing.T) {
	config := TestConfig
	config.QualifiedSeparator = "_"
	config.RPCPrefixes = true
	elm := config.testModule(t, `
		syntax = "proto3";
		message Empty {}
		service Hello {
			rpc World(Empty) returns (Empty);
		}
	`)
	assert.Len(t, elm.Records, 1)
	assert.Len(t, elm.RPCs, 1)
	rpc := elm.RPCs[0]
	assert.Equal(t, protoreflect.FullName("Hello"), rpc.Service)
	assert.Equal(t, protoreflect.Name("World"), rpc.Method)
	assert.Equal(t, "hello_World", rpc.MethodID)
}
