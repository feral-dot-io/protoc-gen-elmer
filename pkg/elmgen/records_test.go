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

func TestScalarRecord(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.scalar;
		message Scalar {
			double my_double = 1;
			float my_float = 2;
			int32 my_int32 = 3;
			//int64 my_int64 = 4;
			uint32 my_uint32 = 5;
			//uint64 my_uint64 = 6;
			sint32 my_sint32 = 7;
			//sint64 my_sint64 = 8;
			fixed32 my_fixed32 = 9;
			//fixed64 my_fixed64 = 10;
			sfixed32 my_sfixed32 = 11;
			//sfixed64 my_sfixed64 = 12;
			bool my_bool = 13;
			string my_string = 14;
			// Field label intentionally used as a reserved word
			bytes type = 15;
		}
	`)
	assert.Equal(t, []string{"Bytes", importElmer, importElmerTests,
		"Test.ScalarTests"}, elm.Imports)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 1)
	scalar := elm.Records[0]
	// IDs
	assert.Equal(t, "Scalar", scalar.Type.ID)
	assert.Equal(t, "emptyScalar", scalar.Type.Zero.ID)
	assert.Equal(t, "decodeScalar", scalar.Type.Decoder.ID)
	assert.Equal(t, "encodeScalar", scalar.Type.Encoder.ID)
}

func TestRecordFieldEnum(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.enum_field;
		message Loaded {
			Question my_question = 123;
		}
		enum Question {
			MAYBE = 0;
			YES = 1;
			NO = 2;
		}`)
	assert.Len(t, elm.Unions, 1)
	assert.Len(t, elm.Records, 1)
}

func TestRecordFieldMessage(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.msg_field;
		message Second {
			double my_two = 123;
		}
		message First {
			double my_one = 123;
		}
		message Zombined {
			First first = 1;
			Second second = 2;
		}`)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 3)
}

func TestRecordNestedMessage(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.nested;
		message Nested {
			message First {
				double my_one = 123;
			}
			message Second {
				double my_two = 123;
			}

			First first = 1;
			Second second = 2;
		}
		message Other {
			Nested.First far_out = 123;
		}`)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 4)
}

func TestFieldErrors(t *testing.T) {
	// Can't create a kind so use an unsupported
	plugin := testPlugin(t, `
		syntax = "proto3";
		package test.oops;
		message Oops {
			int64 my_int64 = 1;
		}`)
	field := plugin.Files[0].Messages[0].Fields[0]
	assert.Equal(t, "my_int64", string(field.Desc.Name()))
	// Field type
	assert.Panics(t, func() {
		fieldTypeDesc(nil, field.Desc)
	})
	// Field zero
	assert.Panics(t, func() {
		fieldZero(nil, field.Desc)
	})
	// Field codec
	assert.Panics(t, func() {
		fieldCodecKind(nil, "PE.", field.Desc)
	})
	// Field fuzzer
	assert.Panics(t, func() {
		fieldFuzzer(nil, field.Desc)
	})
	// General error path (fails on fieldType)
	assert.Panics(t, func() {
		elm := NewModule("", FilesToPackages(plugin.Files)[0])
		g := plugin.NewGeneratedFile("file", "")
		GenerateCodec(elm, g)
	})
}

func TestListField(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.list;
		message Lister {
			repeated bool on_repeat = 11;
		}`)
	assert.Len(t, elm.Records, 1)
}

func TestMapField(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package test.map;
		message A {
			bool first_field = 1;
			map<string, int32> a = 2;
		}
		message B {
			map<string, A> b = 1;
		}
		message C {
			message MimicEntry {
				string key = 1;
				int32 value = 2;
			}
			repeated MimicEntry mimic = 1;			  
		}`)
	assert.Len(t, elm.Records, 4)
	// Can we mimic a map entry? No. https://developers.google.com/protocol-buffers/docs/proto3#backwards_compatibility
}

func TestOneOf(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		message Multi {
			bool leading = 100;
			oneof pick_one {
				string a_str = 1;
				bool a_bool = 2;
			}
			oneof pick_another {
				string b_str = 3;
				int32 b_num = 4;
			}
		}`)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 1)
	assert.Len(t, elm.Oneofs, 2)
	r := elm.Records[0]
	assert.Len(t, r.Oneofs(), 2)
	assert.Equal(t, "Multi", r.Type.ID)
}

func TestOptionalField(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		message Night {
			optional bool shadow = 666;
		}`)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 1)
	assert.Len(t, elm.Oneofs, 1)
	r := elm.Records[0]
	assert.Equal(t, "Night", r.Type.ID)
	// Again but nested
	elm = testModule(t, `
		syntax = "proto3";
		message Day {
			message Night {
				optional bool shadow = 666;
			}
			Night sun = 1;
		}`)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.Oneofs, 1)
	r = elm.Records[0]
	assert.Equal(t, "Day", r.Type.ID)
	r = elm.Records[1]
	assert.Equal(t, "Day_Night", r.Type.ID)
}
