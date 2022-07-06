package elmgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	for i, exp := range []struct {
		Local, Zero, Decode, Encode, Fuzzer string

		Default  string
		Variants []string
	}{
		{"Abc", "emptyAbc", "abcDecoder", "abcEncoder", "abcFuzzer",
			"A", []string{"A", "B", "C"}},
		{"Choose", "emptyChoose", "chooseDecoder", "chooseEncoder", "chooseFuzzer",
			"Hands", []string{"Hands", "Foil", "Epee", "Sabre"}},
		{"Minimal", "emptyMinimal", "minimalDecoder", "minimalEncoder", "minimalFuzzer",
			"Lower", []string{"Lower"}},
	} {
		union := elm.Unions[i]
		// IDs
		assert.Equal(t, exp.Local, union.Type.ID)
		assert.Equal(t, exp.Zero, union.Type.Zero.ID)
		assert.Equal(t, exp.Decode, union.Type.Decoder.ID)
		assert.Equal(t, exp.Encode, union.Type.Encoder.ID)
		// Default
		assert.Equal(t, exp.Default, union.Default().ID.ID)
		assert.Zero(t, union.Default().Number)
		// Variants
		assert.Len(t, union.Variants, len(exp.Variants))
		for j, v := range union.Variants {
			exp := exp.Variants[j]
			assert.Equal(t, exp, v.ID.ID)
		}
	}
}

func TestUnionAllowAlias(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		enum Alias {
			option allow_alias = true;
			UNKNOWN = 0;
			STARTED = 1; // The original
			RUNNING = 1; // This is the alias 
			STOPPED = 2;
		}`)
	assert.Len(t, elm.Unions, 1)
	alias := elm.Unions[0]
	assert.Len(t, alias.Variants, 3)
	assert.Len(t, alias.Aliases, 1)
	assert.Equal(t, "Unknown", alias.Variants[0].ID.ID)
	assert.Equal(t, "Started", alias.Variants[1].ID.ID)
	assert.Equal(t, "running", alias.Aliases[0].Alias.ID)
	assert.Equal(t, "Started", alias.Aliases[0].Variant.ID.String())
	// Check comments
	assert.Contains(t, alias.Variants[1].Comments.Trailing, "The original")
	assert.Contains(t, alias.Aliases[0].Comments.Trailing, "This is the alias")
}
