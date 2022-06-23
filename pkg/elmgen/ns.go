package elmgen

import (
	"log"

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
		m.registerProtoName(proto.Desc.FullName(), "")
		m.protoMessages = append(m.protoMessages, proto)
		// Register oneofs
		prefix := m.aliasName(proto.Desc)
		for _, oneof := range proto.Oneofs {
			od := oneof.Desc
			m.registerProtoName(od.FullName(), string(prefix+"."+od.Name()))
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
		m.registerProtoName(pd.FullName(), "")
		m.protoEnums = append(m.protoEnums, proto)
		for _, protoVal := range proto.Values {
			vd := protoVal.Desc
			m.registerProtoName(vd.FullName(), m.variantAlias(pd, vd))
		}
	}
}

func (m *Module) aliasName(pd protoreflect.Descriptor) protoreflect.Name {
	// Use full (qualified) or a minimal name?
	if m.config.QualifyNested {
		return protoreflect.Name(pd.FullName())
	} else {
		return pd.Name()
	}
}

func (m *Module) variantAlias(enum, value protoreflect.Descriptor) string {
	var suffix protoreflect.Name
	// Suffix variants?
	if m.config.VariantSuffixes {
		suffix = "." + m.aliasName(enum)
	}
	return string(value.Name() + suffix)
}
