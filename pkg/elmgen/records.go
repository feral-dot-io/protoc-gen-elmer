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

// Adds records to a module converting from proto messages
func (m *Module) addRecords(msgs []*protogen.Message) {
	for _, msg := range msgs {
		// Defer map handling to fields
		if msg.Desc.IsMapEntry() {
			continue
		}

		m.Records = append(m.Records, m.newRecord(msg))
		// Add nested
		m.addUnions(msg.Enums)
		m.addRecords(msg.Messages)
	}
	sort.Sort(m.Records)
	sort.Sort(m.Oneofs)
}

// Creates a new record deriving it from a proto message
func (m *Module) newRecord(msg *protogen.Message) *Record {
	md := msg.Desc
	var record Record
	record.Type = m.NewElmType(md.ParentFile(), md)
	record.Comments = newCommentSet(msg.Comments)
	oneofsSeen := make(map[protoreflect.FullName]bool)

	for _, field := range msg.Fields {
		// Oneof field?
		if field.Oneof != nil {
			od := field.Oneof.Desc
			// Only add once
			if oneofsSeen[od.FullName()] {
				continue
			}
			oneofsSeen[od.FullName()] = true

			oneof, field := m.newOneofField(field.Oneof)
			m.Oneofs = append(m.Oneofs, oneof)
			record.Fields = append(record.Fields, field)
		} else { // Regular field
			field := m.newField(field)
			record.Fields = append(record.Fields, field)
		}
	}
	return &record
}

// Creates a new Elm field from a proto field
func (m *Module) newField(field *protogen.Field) *Field {
	fd := field.Desc
	return &Field{
		protoIdentToElmValue(string(fd.Name())),
		fd, nil,
		newCommentSet(field.Comments)}
}

// Creates a field that holds a proto Oneof type. This diverges from protobuf that unwraps the oneof fields and embeds them into the parent message; we want to be closer to our expected Elm code where a oneof has a single field in the record (message).
func (m *Module) newOneofField(protoOneof *protogen.Oneof) (*Oneof, *Field) {
	od := protoOneof.Desc
	oneof := m.newOneof(protoOneof)
	field := &Field{
		protoIdentToElmValue(string(od.Name())),
		nil, oneof,
		newCommentSet(protoOneof.Comments)}
	// Optional field?
	if oneof.IsSynthetic {
		// Unwrap type
		field.Desc = od.Fields().Get(0)
		field.Label = protoIdentToElmValue(string(field.Desc.Name()))
	}
	return oneof, field
}

func (r *Record) Oneofs() (fields []*Field) {
	for _, f := range r.Fields {
		if o := f.Oneof; o != nil {
			fields = append(fields, f)
		}
	}
	return
}
