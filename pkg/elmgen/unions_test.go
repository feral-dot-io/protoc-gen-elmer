package elmgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestUnions(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		enum Choose {
			HANDS = 0;
			FOIL = 1;
			EPEE = 2;
			SABRE = 3;
		}
		enum ABC {
			A = 0;
			B = 1;
			C = 2;
		}
		enum Minimal {
			lower = 0;
		}`)
	assert.Len(t, elm.Unions, 3)
	assert.Empty(t, elm.Records)
	type expVariant struct {
		Local  string
		Number protoreflect.EnumNumber
	}
	for i, exp := range []struct {
		Local, Zero, Decode, Encode, Fuzzer string

		Default  string
		Variants []expVariant
	}{
		{"Abc", "emptyAbc", "abcDecoder", "abcEncoder", "abcFuzzer",
			"A", []expVariant{{"B", 1}, {"C", 2}}},
		{"Choose", "emptyChoose", "chooseDecoder", "chooseEncoder", "chooseFuzzer",
			"Hands", []expVariant{
				{"Foil", 1},
				{"Epee", 2},
				{"Sabre", 3}}},
		{"Minimal", "emptyMinimal", "minimalDecoder", "minimalEncoder", "minimalFuzzer",
			"Lower", []expVariant{}},
	} {
		union := elm.Unions[i]
		// IDs
		assert.Equal(t, exp.Local, union.Type.Local())
		assert.Equal(t, exp.Zero, union.Type.Zero().Local())
		assert.Equal(t, exp.Decode, union.Type.Decoder().Local())
		assert.Equal(t, exp.Encode, union.Type.Encoder().Local())
		// Default
		assert.Equal(t, exp.Default, union.DefaultVariant.ID.Local())
		assert.Zero(t, union.DefaultVariant.Number)
		// Variants
		assert.Len(t, union.Variants, len(exp.Variants))
		for j, v := range union.Variants {
			expVar := exp.Variants[j]
			assert.Equal(t, expVar.Local, v.ID.Local())
			assert.Equal(t, expVar.Number, v.Number)
		}
	}
}

func TestUnionAllowAlias(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		enum Alias {
			option allow_alias = true;
			UNKNOWN = 0;
			STARTED = 1;
			RUNNING = 1;
		}`)
	assert.Len(t, elm.Unions, 1)
	alias := elm.Unions[0]
	assert.Len(t, alias.Variants, 1)
	assert.Len(t, alias.Aliases, 1)
	assert.Equal(t, "Started", alias.Variants[0].ID.Local())
	assert.Equal(t, "running", alias.Aliases[0].Alias.Local())
	assert.Equal(t, "X.Started", alias.Aliases[0].ID.String())
}
