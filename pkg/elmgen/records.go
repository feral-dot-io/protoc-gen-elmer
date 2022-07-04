package elmgen

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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

func (m *Module) newRecord(msg *protogen.Message) *Record {
	md := msg.Desc
	var record Record
	record.Type = m.NewElmType(md.ParentFile(), md)
	record.Comments = NewCommentSet(msg.Comments)
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
			record.Oneofs = append(record.Oneofs, oneof)
			record.Fields = append(record.Fields, field)
		} else { // Regular field
			field := m.newField(field)
			record.Fields = append(record.Fields, field)
		}
	}
	return &record
}

func (m *Module) newField(field *protogen.Field) *Field {
	fd := field.Desc
	return &Field{
		protoFullIdentToElmCasing(string(fd.Name()), "", false),
		fd, nil,
		NewCommentSet(field.Comments)}
}

func (m *Module) newOneofField(protoOneof *protogen.Oneof) (*Oneof, *Field) {
	od := protoOneof.Desc
	oneof := m.newOneof(protoOneof)
	field := &Field{
		protoFullIdentToElmCasing(string(od.Name()), "", false),
		nil, oneof,
		NewCommentSet(protoOneof.Comments)}
	// Optional field?
	if oneof.IsSynthetic {
		// Unwrap type
		field.Desc = od.Fields().Get(0)
		field.Label = string(field.Desc.Name())
	}
	return oneof, field
}
