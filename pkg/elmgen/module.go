package elmgen

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (config Config) newModule() *Module {
	return &Module{
		config:       config,
		protoNS:      make(map[protoreflect.FullName]ElmType),
		protoAliases: make(map[protoreflect.FullName]string),
		elmNS:        make(map[string]struct{})}
}

func (config *Config) NewModule(proto *protogen.File) (*Module, error) {
	m := config.newModule()
	// Check config is valid
	if !ValidPartialElmID(m.config.QualifiedSeparator) {
		return nil, fmt.Errorf("qualified separator must be a valid Elm identifier, got `%s`",
			m.config.QualifiedSeparator)
	}
	if !ValidPartialElmID(m.config.CollisionSuffix) {
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
	m.addEnums()
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

// Convert a protobuf full ident to a valid Elm ID following naming conventions. Diverging from naming conversions, proto namespaces can be joined with a separator.
// The proto3 spec says idents are limited to "alphanum plus _". A full ident is joins idents by a dot.
// If isType is true then first character is uppercased, otherwise it's lowercased.
func (m *Module) protoFullIdentToElmCasing(from string, isType bool) string {
	var segments []string
	var buf []rune
	var prev, prev2 rune
	appendBuf := func() {
		if len(buf) == 0 { // No characters, e.g., just underscore(s)
			x := 'x'
			if isType {
				x = 'X'
			}
			buf = append(buf, x)
		}
		segments = append(segments, string(buf))
		buf = nil
	}
	for _, r := range from {
		original := r
		if r == '.' { // End of namespace
			appendBuf()
		} else if r != '_' { // Skip underscores
			if len(buf) == 0 {
				// First character must be upper or lower
				if isType {
					r = unicode.ToUpper(r)
				} else {
					isType = true // Uppercase subsequent segments
					r = unicode.ToLower(r)
				}
			} else if prev == '_' {
				// Uppercase after _
				r = unicode.ToUpper(r)
			} else if unicode.IsUpper(prev) {
				// Lowercase sequences of uppercase
				r = unicode.ToLower(r)
			} else if unicode.IsUpper(prev2) {
				// Revert last lowercased in sequence of uppercase
				buf[len(buf)-2] = prev2
			} // Otherwise, don't change case
			buf = append(buf, r)
		}
		prev2 = prev
		prev = original
	}
	// Add leftover buffer
	appendBuf()
	return strings.Join(segments, m.config.QualifiedSeparator)
}

func (m *Module) protoFullIdentToElmType(name protoreflect.FullName, isType bool) string {
	// Aliased?
	alias := m.protoAliases[name]
	if alias == "" {
		// No alias, drop pkg prefix from full
		alias = strings.TrimPrefix(string(name), string(m.protoPkg))
	}
	return m.protoFullIdentToElmCasing(alias, isType)
}

func (m *Module) getElmType(name protoreflect.FullName) (ElmType, error) {
	elmID, ok := m.protoNS[name]
	// This should never happen. All proto names should be registered before use
	if !ok {
		log.Panicf("missing protoreflect.FullName: %s", name)
	}
	// First retrieval: create
	var err error
	if elmID == "" {
		candidate := m.protoFullIdentToElmType(name, true)
		candidate, err = m.registerElmID(candidate)
		elmID = ElmType(candidate)
		m.protoNS[name] = elmID
	}
	return elmID, err
}

func (m *Module) getElmValue(name protoreflect.FullName) (string, error) {
	candidate := m.protoFullIdentToElmType(name, false)
	return m.registerElmID(candidate)
}

func (m *Module) registerElmID(id string) (string, error) {
	if !ValidElmID(id) {
		log.Panicf("invalid Elm ID: %s", id)
	}
	// Already registered?
	if _, ok := m.elmNS[id]; ok {
		// Generate an error?
		if m.config.CollisionSuffix == "" {
			return "", fmt.Errorf("protobuf schema generates a name collision with ID `%s` and current config (%#v) prevents us from resolving it",
				id, m.config)
		}
		// Add suffix to ID
		id += m.config.CollisionSuffix
		return m.registerElmID(id)
	}
	m.elmNS[id] = struct{}{}
	return id, nil
}

func (d *CodecIDs) register(m *Module, name protoreflect.FullName) (err error) {
	if d.ID, err = m.getElmType(name); err != nil {
		return err
	}
	id := string(d.ID)
	if d.ZeroID, err = m.registerElmID("empty" + id); err != nil {
		return err
	}
	if d.DecodeID, err = m.registerElmID("decode" + id); err != nil {
		return err
	}
	d.EncodeID, err = m.registerElmID("encode" + id)
	return err
}
