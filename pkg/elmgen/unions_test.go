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
		{"Abc", "emptyAbc", "decodeAbc", "encodeAbc", "fuzzAbc",
			"A", []string{"A", "B", "C"}},
		{"Choose", "emptyChoose", "decodeChoose", "encodeChoose", "fuzzChoose",
			"Hands", []string{"Hands", "Foil", "Epee", "Sabre"}},
		{"Minimal", "emptyMinimal", "decodeMinimal", "encodeMinimal", "fuzzMinimal",
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
	assert.Equal(t, "aliasRunning", alias.Aliases[0].Alias.ID)
	assert.Equal(t, "Started", alias.Aliases[0].Variant.ID.String())
	// Check comments
	assert.Contains(t, alias.Variants[1].Comments.Trailing, "The original")
	assert.Contains(t, alias.Aliases[0].Comments.Trailing, "This is the alias")
}

func TestPrefixAndSuffixCollision(t *testing.T) {
	// If we mix prefixes and suffixes from functions we can potentially get a collision
	testModule(t, `
		syntax = "proto3";
		// This generates an emptyDecoder zero fn
		enum Decoder {
			WHATEVER = 0;
		}
		// With suffixed decoder names this generates emptyDecoder (colliding)
		message Empty {
			bool field = 1;
		}`)
}
