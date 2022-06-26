package elmgen

import (
	"fmt"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (config Config) newModule() *Module {
	m := &Module{
		config:       config,
		protoNS:      make(map[protoreflect.FullName]ElmType),
		protoAliases: make(map[protoreflect.FullName]string)}
	s := struct{}{}
	m.elmNS = map[string]struct{}{
		// Reserved words
		// Source: https://github.com/elm/compiler/blob/770071accf791e8171440709effe71e78a9ab37c/compiler/src/Parse/Variable.hs
		"if": s, "then": s, "else": s, "case": s, "of": s, "let": s, "in": s,
		"type": s, "module": s, "where": s, "import": s, "exposing": s, "as": s,
		"port": s,
		// Prelude https://package.elm-lang.org/packages/elm/core/latest/
		"Basics": s, "List": s, "Maybe": s, "Result": s, "String": s,
		"Char": s, "Tuple": s, "Debug": s, "Platform": s, "Cmd": s, "Sub": s,
		// Basics(..)
		"Int": s, "Float": s, "toFloat": s, "round": s, "floor": s,
		"ceiling": s, "truncate": s, "max": s, "min": s, "compare": s, "LT": s,
		"EQ": s, "GT": s, "Bool": s, "True": s, "False": s, "not": s, "xor": s,
		"modBy": s, "remainderBy": s, "negate": s, "abs": s, "clamp": s,
		"sqrt": s, "logBase": s, "e": s, "degrees": s, "radians": s, "turns": s,
		"pi": s, "cos": s, "sin": s, "tan": s, "acos": s, "asin": s, "atan": s,
		"atan2": s, "toPolar": s, "fromPolar": s, "isNaN": s, "isInfinite": s,
		"identity": s, "always": s, "Never": s, "never": s,
		// Other imports (note these overlap with prelude)
		"Just": s, "Nothing": s, "Ok": s, "Err": s, "Program": s,
	}
	return m
}

func (config *Config) NewModule(proto *protogen.File) (*Module, error) {
	m := config.newModule()
	// Check config is valid
	if !validPartialElmID(m.config.QualifiedSeparator) {
		return nil, fmt.Errorf("qualified separator must be a valid Elm identifier, got `%s`",
			m.config.QualifiedSeparator)
	}
	if !validPartialElmID(m.config.CollisionSuffix) {
		return nil, fmt.Errorf("collision suffix must be a valid Elm identifier, got `%s`",
			m.config.CollisionSuffix)
	}
	//
	m.protoPkg = proto.Desc.Package() + "."
	m.Name, m.Path = config.nameAndPath(string(proto.Desc.Package()), proto.GeneratedFilenamePrefix)
	// First pass: get proto Idents
	m.regEnums(proto.Enums)
	m.regMessages(proto.Messages)
	// Next: translate proto -> elm. Ordering matters: name clashes are suffixed
	if err := m.addEnums(); err != nil {
		return nil, err
	}
	if err := m.addRecords(); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Config) nameAndPath(pkg, file string) (name, path string) {
	// Passed an override?
	if c.ModuleName != "" {
		name = c.ModuleName
	} else {
		// Derive pkg from generated file path if missing
		if pkg == "" {
			// Proto3 "package" is a fullIdent
			pkg = strings.TrimFunc(file, func(r rune) bool {
				// Non-alphanum
				return !(unicode.IsLetter(r) || unicode.IsNumber(r))
			})
			pkg = strings.ReplaceAll(pkg, "/", ".")
		}
		// Determine from pkg
		tokens := strings.Split(pkg, ".")
		for i, token := range tokens {
			tokens[i] = strings.Title(token)
		}
		name = c.ModulePrefix + strings.Join(tokens, ".")
	}
	path = strings.ReplaceAll(name, ".", "/")
	return
}
