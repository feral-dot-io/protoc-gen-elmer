// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
package elmgen

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Adds RPCs to a an Elm module from proto services
func (m *Module) addRPCs(services []*protogen.Service) {
	for _, protoService := range services {
		// Build list of methods
		var methods RPCs
		for _, protoMethod := range protoService.Methods {
			methods = append(methods,
				m.newRPC(protoService.Desc, protoMethod))
		}
		// Add service
		sort.Sort(methods)
		m.Services = append(m.Services, &Service{
			string(protoService.Desc.Name()),
			methods,
			newCommentSet(protoService.Comments)})
	}
	sort.Sort(m.Services)
}

// Creates a new Elm RPC from a proto method
func (m *Module) newRPC(sd protoreflect.Descriptor, method *protogen.Method) *RPC {
	md := method.Desc
	in, out := md.Input(), md.Output()
	return &RPC{
		m.NewElmValue(md.ParentFile(), "twirp", md),
		m.NewElmType(in.ParentFile(), in),
		m.NewElmType(out.ParentFile(), out),
		md.IsStreamingClient(), md.IsStreamingServer(),

		sd.FullName(),
		md.Name(),
		newCommentSet(method.Comments)}
}
