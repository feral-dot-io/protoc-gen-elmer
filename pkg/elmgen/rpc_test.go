package elmgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRPC(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.service;
		service HelloWorld {
			rpc Hello2(HelloReq) returns (HelloResp);
			rpc Hello1(HelloReq) returns (HelloResp);
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
		assert.Equal(t, fmt.Sprintf("helloWorld_Hello%d", i+1), rpc.ID.Local())
		assert.Equal(t, "HelloReq", rpc.In.Local())
		assert.Equal(t, "HelloResp", rpc.Out.Local())
		assert.Equal(t, "Test.Service.helloReqEncoder", rpc.In.Encoder().String())
		assert.Equal(t, "Test.Service.helloRespDecoder", rpc.Out.Decoder().String())
		assert.False(t, rpc.InStreaming)
		assert.False(t, rpc.OutStreaming)
	}
	// It would be great to test the generated code but can't process a Cmd
}
