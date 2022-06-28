package elmgen

import (
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
		assert.Equal(t, exp.IsMap, act.IsMap)
		if !exp.IsOneof {
			assert.Equal(t, exp.WireNumber, act.WireNumber)
			assert.Equal(t, exp.Cardinality, act.Cardinality)
		}
		assert.Equal(t, exp.Type, act.Type)
		if exp.Type != "Bytes" {
			assert.Equal(t, exp.Zero, act.Zero)
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
		assert.Equal(t, exp.Key, act.Key)
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
	assert.True(t, elm.Fuzzers.Int32)
	assert.True(t, elm.Fuzzers.Uint32)
	assert.True(t, elm.Fuzzers.Float32)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 1)
	scalar := elm.Records[0]
	// IDs
	assert.Equal(t, ElmType("Scalar"), scalar.ID)
	assert.Equal(t, "emptyScalar", scalar.ZeroID)
	assert.Equal(t, "scalarDecoder", scalar.DecodeID)
	assert.Equal(t, "scalarEncoder", scalar.EncodeID)
	// Fields
	assertFields(t, scalar.Fields,
		&Field{"myDouble", false, false, 1, opt, "Float", "0", "", "", "Fuzz.float", nil},
		&Field{"myFloat", false, false, 2, opt, "Float", "0", "", "", "Fuzz.float", nil},
		&Field{"myInt32", false, false, 3, opt, "Int", "0", "", "", "fuzzInt32", nil},
		&Field{"myUint32", false, false, 5, opt, "Int", "0", "", "", "fuzzUint32", nil},
		&Field{"mySint32", false, false, 7, opt, "Int", "0", "", "", "fuzzInt32", nil},
		&Field{"myFixed32", false, false, 9, opt, "Int", "0", "", "", "fuzzUint32", nil},
		&Field{"mySfixed32", false, false, 11, opt, "Int", "0", "", "", "fuzzUint32", nil},
		&Field{"myBool", false, false, 13, opt, "Bool", "False", "", "", "Fuzz.bool", nil},
		&Field{"myString", false, false, 14, opt, "String", `""`, "", "", "Fuzz.string", nil},
		&Field{"myBytes", false, false, 15, opt, "Bytes", "", "", "", "fuzzBytes", nil})
}

func TestRecordFieldEnum(t *testing.T) {
	elm := TestConfig.testModule(t, `
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
	assertFields(t, elm.Records[0].Fields,
		&Field{"myQuestion", false, false, 123, opt,
			"Question", "emptyQuestion", "", "", "questionFuzzer", nil})
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
		&Field{"first", false, false, 1, opt, "First",
			"emptyFirst", "firstDecoder", "firstEncoder", "firstFuzzer", nil},
		&Field{"second", false, false, 2, opt, "Second",
			"emptySecond", "secondDecoder", "secondEncoder", "secondFuzzer", nil})
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
		&Field{"first", false, false, 1, opt, "Nested_First",
			"emptyNested_First", "nested_FirstDecoder", "nested_FirstEncoder", "nested_FirstFuzzer", nil},
		&Field{"second", false, false, 2, opt, "Nested_Second",
			"emptyNested_Second", "nested_SecondDecoder", "nested_SecondEncoder", "nested_SecondFuzzer", nil})
	assertFields(t, elm.Records[3].Fields,
		&Field{"farOut", false, false, 123, opt, "Nested_First",
			"emptyNested_First", "nested_FirstDecoder", "nested_FirstEncoder", "nested_FirstFuzzer", nil})
}

func TestRecordErrors(t *testing.T) {
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
	_, err := fieldType(nil, field.Desc)
	assert.ErrorContains(t, err, "protoreflect.Kind")
	// Field zero
	_, err = fieldZero(nil, field.Desc)
	assert.ErrorContains(t, err, "protoreflect.Kind")
	// Field codec
	_, err = fieldKindCodec(nil, "PE.", "encode", field.Desc)
	assert.ErrorContains(t, err, "protoreflect.Kind")
	// Field fuzzer
	_, err = fieldFuzzer(nil, field.Desc)
	assert.ErrorContains(t, err, "protoreflect.Kind")
	// General error path (fails on fieldType)
	_, err = TestConfig.NewModule(plugin.Files[0])
	assert.ErrorContains(t, err, "protoreflect.Kind")
}

func TestListField(t *testing.T) {
	config := TestConfig
	elm := config.testModule(t, `
		syntax = "proto3";
		package test.list;
		message Lister {
			repeated bool on_repeat = 11;
		}`)
	assert.Len(t, elm.Records, 1)
	assertFields(t, elm.Records[0].Fields,
		&Field{"onRepeat", false, false, 11, protoreflect.Repeated,
			"(List Bool)", "[]", "", "", "", nil})
}

func TestMapField(t *testing.T) {
	elm := TestConfig.testModule(t, `
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
	assert.True(t, elm.Imports.Dict)
	// Basic
	stringKey := &MapKey{`""`, "PD.string", "PE.string", "Fuzz.string"}
	assertFields(t, elm.Records[0].Fields,
		&Field{"firstField", false, false, 1, opt, "Bool",
			"False", "", "", "Fuzz.bool", nil},
		&Field{"a", false, true, 2, protoreflect.Repeated,
			"(Dict String Int)", "0", "", "", "", stringKey})
	// Value is a nested message
	assertFields(t, elm.Records[1].Fields,
		&Field{"b", false, true, 1, protoreflect.Repeated,
			"(Dict String A)", "emptyA", "", "", "", stringKey})
	// Can we mimic a map entry? No. https://developers.google.com/protocol-buffers/docs/proto3#backwards_compatibility
	assertFields(t, elm.Records[2].Fields,
		&Field{"mimic", false, false, 1, protoreflect.Repeated,
			"(List C_MimicEntry)", "[]", "", "", "", nil})
}

func TestOneOf(t *testing.T) {
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
		&Field{"pickOne", true, false, 0, 0, "(Maybe Multi_PickOne)", "Nothing",
			"multi_PickOneDecoder", "multi_PickOneEncoder", "multi_PickOneFuzzer", nil},
		&Field{"pickAnother", true, false, 0, 0, "(Maybe Multi_PickAnother)", "Nothing",
			"multi_PickAnotherDecoder", "multi_PickAnotherEncoder", "multi_PickAnotherFuzzer", nil})
}

func TestOptionalField(t *testing.T) {
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
		&Field{"shadow", true, false, 0, 0, "(Maybe Bool)", "Nothing",
			"night_ShadowDecoder", "night_ShadowEncoder", "night_ShadowFuzzer", nil})
	// Again but nested
	elm = TestConfig.testModule(t, `
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
	assert.Equal(t, ElmType("Day"), r.ID)
	assertFields(t, r.Fields,
		&Field{"sun", false, false, 1, opt, "Day_Night", "emptyDay_Night",
			"day_NightDecoder", "day_NightEncoder", "day_NightFuzzer", nil})
	r = elm.Records[1]
	assert.Equal(t, ElmType("Day_Night"), r.ID)
	assertFields(t, r.Fields,
		&Field{"shadow", true, false, 0, 0, "(Maybe Bool)", "Nothing",
			"day_Night_ShadowDecoder", "day_Night_ShadowEncoder", "day_Night_ShadowFuzzer", nil})
}
