package elmgen

import (
	"log"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) registerProtoName(name protoreflect.FullName, alias string) {
	// Don't overwrite. This should never happen. We rely on protobuf semantics being enforced by our caller
	if _, ok := m.protoNS[name]; ok {
		log.Panicf("duplicate protoreflect.FullName: %s", name)
	}
	// Defer Elm ID creation so that we can decide which IDs get suffixed on name collision
	m.protoNS[name] = ""
	m.protoAliases[name] = alias
}

func (m *Module) regMessages(protoMsgs []*protogen.Message) {
	for _, proto := range protoMsgs {
		// Recursive
		m.regEnums(proto.Enums)
		m.regMessages(proto.Messages)
		// Register msg
		prefix := m.aliasName(proto.Desc)
		m.registerProtoName(proto.Desc.FullName(), prefix)
		m.protoMessages = append(m.protoMessages, proto)
		// Register oneofs
		for _, oneof := range proto.Oneofs {
			od := oneof.Desc
			oneofAlias := prefix + "." + string(od.Name())
			if od.IsSynthetic() { // Optional
				// Make oneof an alias of single field
				fd := od.Fields().Get(0)
				oneofAlias = string(fd.FullName())
			}
			m.registerProtoName(od.FullName(), oneofAlias)
			for _, field := range oneof.Fields {
				fd := field.Desc
				m.registerProtoName(fd.FullName(), m.variantAlias(od, fd))
			}
		}
	}
}

func (m *Module) regEnums(protoEnums []*protogen.Enum) {
	for _, proto := range protoEnums {
		pd := proto.Desc
		m.registerProtoName(pd.FullName(), m.aliasName(pd))
		m.protoEnums = append(m.protoEnums, proto)
		for _, protoVal := range proto.Values {
			vd := protoVal.Desc
			m.registerProtoName(vd.FullName(), m.variantAlias(pd, vd))
		}
	}
}

func (m *Module) aliasName(pd protoreflect.Descriptor) string {
	// Use full (qualified) or a minimal name?
	if m.config.QualifyNested {
		full := string(pd.FullName())
		// No alias, drop pkg prefix from full
		return strings.TrimPrefix(full, string(m.protoPkg))
	} else {
		return string(pd.Name())
	}
}

func (m *Module) variantAlias(enum, value protoreflect.Descriptor) string {
	return string(value.Name()) + "." + m.aliasName(enum)
}
