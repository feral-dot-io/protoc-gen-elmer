// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
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
	elm := NewModule("", plugin.Files[0])
	assert.Equal(t, []string{"Bytes", "Dict", "FindTests", importElmer,
		importElmerTests}, elm.Imports)
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
				uint32 b = 3;
				bytes triggered = 4;
				float f = 5;
			};
		}
	`, `
		syntax = "proto3";
		package other;
		message Other {
			map<bool, int32> not_triggered = 1;
		}
	`)
	elm := NewModule("", plugin.Files[1])
	assert.Equal(t, []string{"Bytes", "MyTests", "Other", "OtherTests",
		importElmer, importElmerTests}, elm.Imports)
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
	assert.Equal(t, []string{"AnotherPkg", "AnotherPkgTests", importElmerTests,
		"XTests"}, elm.Imports)
	assert.Len(t, elm.Records, 1)
	assert.Equal(t, "MyMessage", elm.Records[0].Type.ID)
}

func TestWellKnown(t *testing.T) {
	testModule(t, `
		syntax = "proto3";
		import "google/protobuf/any.proto";
		import "google/protobuf/api.proto";
		import "google/protobuf/duration.proto";
		import "google/protobuf/empty.proto";
		import "google/protobuf/field_mask.proto";
		import "google/protobuf/timestamp.proto";
		import "google/protobuf/type.proto";
		import "google/protobuf/source_context.proto";
		import "google/protobuf/struct.proto";
		import "google/protobuf/wrappers.proto";
		message Famous {
			google.protobuf.Any any = 1;
			google.protobuf.Api api = 2;
			google.protobuf.BoolValue bool_value = 3;
			google.protobuf.BytesValue bytes_value = 4;
			google.protobuf.DoubleValue double_value = 5;
			google.protobuf.Duration stitch_in_time = 6;
			google.protobuf.Empty empty = 7;
			google.protobuf.Enum enum = 8;
			google.protobuf.EnumValue enum_value = 9;
			google.protobuf.Field field = 10;
			google.protobuf.Field.Cardinality field_cardinality = 11;
			google.protobuf.Field.Kind field_kind = 12;
			google.protobuf.FieldMask field_mask = 13;
			google.protobuf.FloatValue float_value = 14;
			google.protobuf.Int32Value int32_value = 15;
			google.protobuf.Int64Value int64_value = 16;
			google.protobuf.ListValue list_value = 17;
			google.protobuf.Method method = 18;
			google.protobuf.Mixin mixin = 19;
			// Missing decoder / encoder in Google.Protobuf
			//google.protobuf.NullValue null_value = 20;
			google.protobuf.Option option = 21;
			//google.protobuf.SourceContext source_context = 22;
			google.protobuf.StringValue string_value = 23;
			google.protobuf.Struct struct = 24;
			google.protobuf.Syntax syntax = 25;
			google.protobuf.Timestamp now_or_never = 26;
			google.protobuf.Type type = 27;
			google.protobuf.UInt32Value uint32_value = 28;
			google.protobuf.UInt64Value uint64_value = 29;
			google.protobuf.Value value = 30;
		}`)
}
