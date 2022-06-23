package elmgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var opt = protoreflect.Optional

func assertFields(t *testing.T, fields []*Field, exp ...*Field) {
	t.Helper()
	assert.Len(t, fields, len(exp))
	for i, exp := range exp {
		act := fields[i]
		assert.Equal(t, exp.Label, act.Label)
		assert.Equal(t, exp.IsOneof, act.IsOneof)
		if !exp.IsOneof {
			assert.Equal(t, exp.WireNumber, act.WireNumber)
			assert.Equal(t, exp.Cardinality, act.Cardinality)
		}
		assert.Equal(t, exp.Type, act.Type)
		if exp.Type != "Bytes" {
			zero := fmt.Sprintf("%s", act.Zero)
			assert.Equal(t, exp.Zero, zero)
		} else {
			assert.NotEmpty(t, act.Zero)
		}
		// We rely on Elm tests to check the contents
		if exp.Decoder == "" {
			assert.NotEmpty(t, act.Decoder)
		} else {
			assert.Equal(t, exp.Decoder, act.Decoder)
		}
		if exp.Encoder == "" {
			assert.NotEmpty(t, act.Encoder)
		} else {
			assert.Equal(t, exp.Encoder, act.Encoder)
		}
	}
}

func TestScalarRecord(t *testing.T) {
	elm := TestConfig.testModule(t, `
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
			bytes my_bytes = 15;
		}
	`)
	assert.True(t, elm.Imports.Bytes)
	assert.False(t, elm.Imports.Dict)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 1)
	scalar := elm.Records[0]
	// IDs
	assert.Equal(t, ElmType("Scalar"), scalar.ID)
	assert.Equal(t, "emptyScalar", scalar.ZeroID)
	assert.Equal(t, "decodeScalar", scalar.DecodeID)
	assert.Equal(t, "encodeScalar", scalar.EncodeID)
	// Fields
	assertFields(t, scalar.Fields,
		&Field{"myDouble", false, 1, opt, "Float", "0", "", ""},
		&Field{"myFloat", false, 2, opt, "Float", "0", "", ""},
		&Field{"myInt32", false, 3, opt, "Int", "0", "", ""},
		&Field{"myUint32", false, 5, opt, "Int", "0", "", ""},
		&Field{"mySint32", false, 7, opt, "Int", "0", "", ""},
		&Field{"myFixed32", false, 9, opt, "Int", "0", "", ""},
		&Field{"mySfixed32", false, 11, opt, "Int", "0", "", ""},
		&Field{"myBool", false, 13, opt, "Bool", "False", "", ""},
		&Field{"myString", false, 14, opt, "String", `""`, "", ""},
		&Field{"myBytes", false, 15, opt, "Bytes", "", "", ""})
}

const testRecordFieldEnumProto = `
	syntax = "proto3";
	package test.enum_field;
	message Loaded {
		Question my_question = 123;
	}
	enum Question {
		MAYBE = 0;
		YES = 1;
		NO = 2;
	}`

func TestRecordFieldEnum(t *testing.T) {
	elm := TestConfig.testModule(t, testRecordFieldEnumProto)
	assert.Len(t, elm.Unions, 1)
	assert.Len(t, elm.Records, 1)
	assertFields(t, elm.Records[0].Fields,
		&Field{"myQuestion", false, 123, opt, "Question", "emptyQuestion", "", ""})
}

func TestRecordFieldMessage(t *testing.T) {
	elm := TestConfig.testModule(t, `
		syntax = "proto3";
		package test.msg_field;
		message First {
			double my_one = 123;
		}
		message Second {
			double my_two = 123;
		}
		message Zombined {
			First first = 1;
			Second second = 2;
		}`)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 3)
	assertFields(t, elm.Records[2].Fields,
		&Field{"first", false, 1, opt, "First", "emptyFirst", "decodeFirst", "encodeFirst"},
		&Field{"second", false, 2, opt, "Second", "emptySecond", "decodeSecond", "encodeSecond"})
}

func TestRecordNestedMessage(t *testing.T) {
	elm := TestConfig.testModule(t, `
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
	assertFields(t, elm.Records[0].Fields,
		&Field{"first", false, 1, opt, "NestedFirst",
			"emptyNestedFirst", "decodeNestedFirst", "encodeNestedFirst"},
		&Field{"second", false, 2, opt, "NestedSecond",
			"emptyNestedSecond", "decodeNestedSecond", "encodeNestedSecond"})
	assertFields(t, elm.Records[3].Fields,
		&Field{"farOut", false, 123, opt, "NestedFirst",
			"emptyNestedFirst", "decodeNestedFirst", "encodeNestedFirst"})
}

func TestRecordErrors(t *testing.T) {
	plugin := testPlugin(t, `
		syntax = "proto3";
		package test.oops;
		message Oops {
			int64 my_int64 = 1;
		}`)
	field := plugin.Files[0].Messages[0].Fields[0]
	assert.Equal(t, "my_int64", string(field.Desc.Name()))
	// Field type
	_, err := fieldType(nil, field.Desc)
	assert.ErrorContains(t, err, "protoreflect.Kind")
	// Field zero
	_, err = fieldZero(nil, field.Desc)
	assert.ErrorContains(t, err, "protoreflect.Kind")
	// Field codec
	_, err = fieldCodec(nil, "PE.", "encode", field.Desc)
	assert.ErrorContains(t, err, "protoreflect.Kind")
	// General error path (fails on fieldType)
	_, err = TestConfig.NewModule(plugin.Files[0])
	assert.ErrorContains(t, err, "protoreflect.Kind")
}

func TestListField(t *testing.T) {
	elm := TestConfig.testModule(t, `
		syntax = "proto3";
		package test.list;
		message List {
			repeated bool on_repeat = 11;
		}`)
	assert.Len(t, elm.Records, 1)
	assertFields(t, elm.Records[0].Fields,
		&Field{"onRepeat", false, 11, protoreflect.Repeated, "(List Bool)", "[]", "", ""})
}

func TestMapField(t *testing.T) {
	// First, the general case
	elm := TestConfig.testModule(t, `
		syntax = "proto3";
		package test.map;
		message A {
			map<string, int32> a = 1;
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
	assert.True(t, elm.Imports.Dict)
	// Basic
	assertFields(t, elm.Records[0].Fields,
		&Field{"a", false, 1, protoreflect.Repeated,
			"(Dict String Int)", "Dict.empty", "", ""})
	// Value is a nested message
	assertFields(t, elm.Records[1].Fields,
		&Field{"b", false, 1, protoreflect.Repeated,
			"(Dict String A)", "Dict.empty", "", ""})
	// Can we mimic a map entry? No. https://developers.google.com/protocol-buffers/docs/proto3#backwards_compatibility
	assertFields(t, elm.Records[2].Fields,
		&Field{"mimic", false, 1, protoreflect.Repeated,
			"(List CMimicEntry)", "[]", "", ""})
}

func TestOneOf(t *testing.T) {
	// First, the general case
	elm := TestConfig.testModule(t, `
		syntax = "proto3";
		message Multi {
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
	assert.Len(t, r.Oneofs, 2)
	assert.Equal(t, ElmType("Multi"), r.ID)
	assertFields(t, r.Fields,
		&Field{"pickOne", true, 0, 0, "Maybe MultiPickOne", "Nothing",
			"decodeMultiPickOne", "encodeMultiPickOne"},
		&Field{"pickAnother", true, 0, 0, "Maybe MultiPickAnother", "Nothing",
			"decodeMultiPickAnother", "encodeMultiPickAnother"})
}

func TestOptionalField(t *testing.T) {
	// First, the general case
	elm := TestConfig.testModule(t, `
		syntax = "proto3";
		message Night {
			optional bool shadow = 666;
		}`)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 1)
	assert.Len(t, elm.Oneofs, 1)
	r := elm.Records[0]
	assert.Equal(t, ElmType("Night"), r.ID)
	assertFields(t, r.Fields,
		&Field{"shadow", true, 0, 0, "Maybe Bool", "Nothing",
			"decodeNightShadow", "encodeNightShadow"})
}
