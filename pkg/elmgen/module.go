package elmgen

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (config Config) newModule() *Module {
	m := &Module{
		config:       config,
		protoNS:      make(map[protoreflect.FullName]ElmType),
		protoAliases: make(map[protoreflect.FullName]string),
		elmNS:        make(map[string]struct{})}
	reserved := []string{
		// Reserved words
		// Source: https://github.com/elm/compiler/blob/770071accf791e8171440709effe71e78a9ab37c/compiler/src/Parse/Variable.hs
		"if", "then", "else", "case", "of", "let", "in", "type", "module",
		"where", "import", "exposing", "as", "port",
		// Prelude https://package.elm-lang.org/packages/elm/core/latest/
		"Basics", "List", "Maybe", "Result", "String", "Char", "Tuple", "Debug",
		"Platform", "Cmd", "Sub",
		// Basics(..)
		"Int", "Float", "toFloat", "round", "floor", "ceiling", "truncate",
		"max", "min", "compare", "LT", "EQ", "GT", "Bool", "True", "False",
		"not", "xor", "modBy", "remainderBy", "negate", "abs", "clamp", "sqrt",
		"logBase", "e", "degrees", "radians", "turns", "pi", "cos", "sin",
		"tan", "acos", "asin", "atan", "atan2", "toPolar", "fromPolar", "isNaN",
		"isInfinite", "identity", "always", "Never", "never",
		// Other imports (note these overlap with prelude)
		"Just", "Nothing", "Ok", "Err", "Program"}
	for _, id := range reserved {
		m.elmNS[id] = struct{}{}
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
	// Paths
	m.protoPkg = proto.Desc.Package() + "."
	m.Name, m.Path = config.nameAndPath(string(proto.Desc.Package()), proto.GeneratedFilenamePrefix)
	// First pass: get proto Idents
	m.regEnums(proto.Enums)
	m.regMessages(proto.Messages)
	m.regMethods(proto.Services)
	// Next: translate proto -> elm. Ordering matters: name clashes are suffixed
	if err := m.addEnums(); err != nil {
		return nil, err
	}
	if err := m.addRecords(); err != nil {
		return nil, err
	}
	if err := m.addRPCs(); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Config) nameAndPath(pkg, file string) (name, path string) {
	// Passed an override?
	if c.ModuleName != "" {
		name = c.ModuleName
	} else {
		// Derive pkg from generated file path if missing (no package = ...)
		if pkg == "" {
			file = file[:len(file)-len(filepath.Ext(file))]
			file = strings.ReplaceAll(file, "/", ".")
			// Turn file into a fullIdent
			pkg = strings.TrimFunc(file, func(r rune) bool {
				// Remove non-alphanum and _
				return !(unicode.IsLetter(r) ||
					unicode.IsNumber(r)) ||
					r == '_'
			})
		}
		name = protoFullIdentToElmCasing(pkg, ".", true)
	}
	path = strings.ReplaceAll(name, ".", "/")
	return
}
