package elmgen

import (
	"log"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (config Config) newModule() *Module {
	if config.CollisionSuffix == "" {
		config.CollisionSuffix = DefaultConfig.CollisionSuffix
	}
	return &Module{
		config:       config,
		protoNS:      make(map[protoreflect.FullName]ElmType),
		protoAliases: make(map[protoreflect.FullName]string),
		elmNS:        make(map[string]struct{})}
}

func (config *Config) NewModule(proto *protogen.File) (*Module, error) {
	m := config.newModule()
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

func (m *Module) registerProtoName(name protoreflect.FullName, alias string) {
	// Don't overwrite. This should never happen. We rely on protobuf semantics being enforced by our caller
	if _, ok := m.protoNS[name]; ok {
		log.Panicf("duplicate protoreflect.FullName: %s", name)
	}
	// Defer Elm ID creation so that we can decide which IDs get suffixed on name collision
	m.protoNS[name] = ""
	m.protoAliases[name] = alias
}

// Convert a protobuf full ident to a valid Elm ID following naming conventions. Diverging from naming conversions, proto namespaces can be joined with a separator.
// The proto3 spec says idents are limited to "alphanum plus _". A full ident is joins idents by a dot.
// If isType is true then first character is uppercased, otherwise it's lowercased.
func (m *Module) protoFullIdentToElmCasing(from string, isType bool) string {
	var segments []string
	var buf []rune
	var prev, prev2 rune
	for _, r := range from {
		original := r
		if r == '.' { // End of namespace
			if len(buf) > 0 {
				segments = append(segments, string(buf))
				buf = nil
			}
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
	if len(buf) > 0 {
		segments = append(segments, string(buf))
	}
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

func (m *Module) getElmType(name protoreflect.FullName) ElmType {
	elmID, ok := m.protoNS[name]
	// This should never happen. All proto names should be registered before use
	if !ok {
		log.Panicf("missing protoreflect.FullName: %s", name)
	}
	// First retrieval: create
	if elmID == "" {
		candidate := m.protoFullIdentToElmType(name, true)
		elmID = ElmType(m.registerElmID(candidate))
		m.protoNS[name] = elmID
	}
	return elmID
}

func (m *Module) getElmValue(name protoreflect.FullName) string {
	candidate := m.protoFullIdentToElmType(name, false)
	return m.registerElmID(candidate)
}

func (m *Module) registerElmID(id string) string {
	// TODO: check if name looks normal ie., starts with a capital
	// TODO: check for invalid Elm IDs
	if _, ok := m.elmNS[id]; ok {
		// Already registered, add a suffix to last ID
		id += m.config.CollisionSuffix
		return m.registerElmID(id)
	}
	m.elmNS[id] = struct{}{}
	return id
}

func (d *CodecIDs) register(m *Module, name protoreflect.FullName) {
	d.ID = m.getElmType(name)
	id := string(d.ID)
	d.ZeroID = m.registerElmID("empty" + id)
	d.DecodeID = m.registerElmID("decode" + id)
	d.EncodeID = m.registerElmID("encode" + id)
}
