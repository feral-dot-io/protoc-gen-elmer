package elmgen

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) regMethods(protoServices []*protogen.Service) {
	for _, protoService := range protoServices {
		service := string(protoService.Desc.Name())
		for _, proto := range protoService.Methods {
			pd := proto.Desc
			// Prefix ID with service?
			alias := string(pd.Name())
			if m.config.RPCPrefixes {
				alias = service + "." + alias
			}

			m.registerProtoName(pd.FullName(), alias)
			m.protoMethods = append(m.protoMethods, proto)
		}
	}
}

func (m *Module) addRPCs() error {
	for _, proto := range m.protoMethods {
		rpc, err := m.newRPC(proto)
		if err != nil {
			return err
		}
		m.RPCs = append(m.RPCs, rpc)
	}
	sort.Sort(m.RPCs)
	return nil
}

func (m *Module) newRPC(method *protogen.Method) (*RPC, error) {
	md := method.Desc
	methodID := m.protoFullIdentToElmID(md.FullName(), false)
	in := m.getElmType(md.Input().FullName())
	out := m.getElmType(md.Output().FullName())
	// Already registered. Total mess
	inEncoder := m.getElmValue(protoreflect.FullName(in)) + "Encoder"
	outDecoder := m.getElmValue(protoreflect.FullName(out)) + "Decoder"
	return &RPC{
		method.Parent.Desc.FullName(),
		method.Desc.Name(),
		methodID,
		string(in), string(out), inEncoder, outDecoder,
		md.IsStreamingClient(), md.IsStreamingServer(),
	}, nil
}
