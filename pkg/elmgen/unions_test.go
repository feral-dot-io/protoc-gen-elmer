package elmgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnions(t *testing.T) {
	elm := TestConfig.testModule(t, `
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
		CodecIDs
		Default  string
		Variants []Variant
	}{
		{CodecIDs{"Abc", "emptyAbc", "abcDecoder", "abcEncoder", "abcFuzzer"},
			"A_Abc", []Variant{{"B_Abc", 1}, {"C_Abc", 2}}},
		{CodecIDs{"Choose", "emptyChoose", "chooseDecoder", "chooseEncoder", "chooseFuzzer"},
			"Hands_Choose", []Variant{
				{"Foil_Choose", 1},
				{"Epee_Choose", 2},
				{"Sabre_Choose", 3}}},
		{CodecIDs{"Minimal", "emptyMinimal", "minimalDecoder", "minimalEncoder", "minimalFuzzer"},
			"Lower_Minimal", []Variant{}},
	} {
		union := elm.Unions[i]
		// IDs
		assert.Equal(t, exp.ID, union.ID)
		assert.Equal(t, exp.ZeroID, union.ZeroID)
		assert.Equal(t, exp.DecodeID, union.DecodeID)
		assert.Equal(t, exp.EncodeID, union.EncodeID)
		// Default
		assert.Equal(t, ElmType(exp.Default), union.DefaultVariant.ID)
		assert.Zero(t, union.DefaultVariant.Number)
		// Variants
		assert.Len(t, union.Variants, len(exp.Variants))
		for j, v := range union.Variants {
			expVar := exp.Variants[j]
			assert.Equal(t, ElmType(expVar.ID), v.ID)
			assert.Equal(t, expVar.Number, v.Number)
		}
	}
}

func TestUnionAllowAlias(t *testing.T) {
	elm := TestConfig.testModule(t, `
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
	assert.Equal(t, ElmType("Started_Alias"), alias.Variants[0].ID)
	assert.Equal(t, "running_Alias", alias.Aliases[0].Alias)
	assert.Equal(t, ElmType("Started_Alias"), alias.Aliases[0].ID)
}
