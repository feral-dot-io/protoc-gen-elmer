package elmgen

import (
	"sort"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addRPCs(services protoreflect.ServiceDescriptors) {
	for i := 0; i < services.Len(); i++ {
		sd := services.Get(i)
		methods := sd.Methods()
		for i := 0; i < methods.Len(); i++ {
			md := methods.Get(i)
			rpc := m.newRPC(sd, md)
			m.RPCs = append(m.RPCs, rpc)
		}
	}
	sort.Sort(m.RPCs)
}

func (m *Module) newRPC(sd protoreflect.Descriptor, md protoreflect.MethodDescriptor) *RPC {
	in, out := md.Input(), md.Output()
	return &RPC{
		NewElmValue(md.ParentFile(), md),
		NewElmType(in.ParentFile(), in),
		NewElmType(out.ParentFile(), out),
		md.IsStreamingClient(), md.IsStreamingServer(),

		sd.FullName(),
		md.Name()}
}
