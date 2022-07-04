package elmgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindImports(t *testing.T) {
	plugin := testPlugin(t, `
		syntax = "proto3";
		package find;
		message Find {
			optional bytes my_bytes = 1;
			map<bool, int32> my_map = 2;
		}
	`)
	elm := NewModule("", "", plugin.Files[0])
	assert.Equal(t, []string{"Bytes", "Dict", "FindTests"}, elm.Imports)
}

func TestFindImportsNested(t *testing.T) {
	plugin := testPlugin(t, `
		syntax = "proto3";
		import "test1.proto";
		package my;
		message My {
			other.Other my_other = 1;
			oneof asdf {
				int32 a = 2;
				int32 b = 3;
				bytes triggered = 4;
			};
		}
	`, `
		syntax = "proto3";
		package other;
		message Other {
			map<bool, int32> not_triggered = 1;
		}
	`)
	elm := NewModule("", "", plugin.Files[1])
	assert.Equal(t, []string{"Bytes", "MyTests", "Other", "OtherTests"}, elm.Imports)
}

func TestImports(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		import "test1.proto";
		message MyMessage {
			AnotherPkg.Other out_of_this_world = 1;
		}`, `
		syntax = "proto3";
		package AnotherPkg;
		message Other {
			int32 a = 1;
			int32 b = 2;
			int32 c = 3;
		}`)
	assert.Equal(t, []string{"AnotherPkg", "AnotherPkgTests", "XTests"}, elm.Imports)
	assert.Len(t, elm.Records, 1)
	assert.Equal(t, "MyMessage", elm.Records[0].Type.ID)
	f := elm.Records[0].Fields[0]
	assert.Equal(t, "outOfThisWorld", f.Label)
	assert.Equal(t, "AnotherPkg.Other", f.Type)
	assert.Equal(t, "AnotherPkg.emptyOther", f.Zero)
	assert.Equal(t, "AnotherPkg.otherDecoder", f.Decoder)
	assert.Equal(t, "AnotherPkg.otherEncoder", f.Encoder)
	assert.Equal(t, "AnotherPkgTests.otherFuzzer", f.Fuzzer)
}
