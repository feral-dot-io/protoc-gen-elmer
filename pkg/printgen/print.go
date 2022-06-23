package printgen

import (
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func print(str string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, str, args...)
}

func printTab(indent int, str string, args ...interface{}) {
	print(strings.Repeat("\t", indent)+str, args...)
}

func File(f *protogen.File) {
	print("File %s:\n", f.Desc.Path())
	FileDescriptor(1, f.Desc)
	printTab(1, "GoDescriptorIdent: %s\n", f.GoDescriptorIdent)
	printTab(1, "GoPackageName: %s\n", f.GoPackageName)
	printTab(1, "GoImportPath: %s\n", f.GoImportPath)
	if f.Generate {
		printTab(1, "GeneratedFilenamePrefix: %s\n", f.GeneratedFilenamePrefix)
	} else {
		printTab(1, "no file generation\n")
	}
	Enums(1, f.Enums)
	Messages(1, f.Messages)
	Extensions(1, f.Extensions)
	Services(1, f.Services)
	print("\n")
}

func FileDescriptor(indent int, fd protoreflect.FileDescriptor) {
	printTab(indent, "fd.Name: %s\n", fd.Name())
	printTab(indent, "fd.Syntax: %s\n", fd.Syntax())
	printTab(indent, "fd.Path: %s\n", fd.Path())
	printTab(indent, "fd.Package: %s\n", fd.Package())
	// Imports
	imports := fd.Imports()
	printTab(indent, "fd.Imports (%d):\n", imports.Len())
	for i := 0; i < imports.Len(); i++ {
		imp := imports.Get(i)
		printTab(indent+1, "%s\n", imp.FileDescriptor.Path())
	}
}

/* Iterators */

func iterator(indent, size int, name string, fn func(int)) {
	if size > 0 {
		printTab(indent, "%s (%d):\n", name, size)
		for i := 0; i < size; i++ {
			fn(i)
		}
	}
}

func Enums(indent int, items []*protogen.Enum) {
	iterator(indent, len(items), "Enums", func(i int) {
		Enum(indent+1, items[i])
	})
}

func Extensions(indent int, items []*protogen.Extension) {
	iterator(indent, len(items), "Extensions", func(i int) {
		Field(indent+1, items[i])
	})
}

func Fields(indent int, items []*protogen.Field) {
	iterator(indent, len(items), "Fields", func(i int) {
		Field(indent+1, items[i])
	})
}

func OneOfs(indent int, items []*protogen.Oneof) {
	iterator(indent, len(items), "Oneofs", func(i int) {
		Oneof(indent+1, items[i])
	})
}

func Services(indent int, items []*protogen.Service) {
	iterator(indent, len(items), "Services", func(i int) {
		Service(indent+1, items[i])
	})
}

func Messages(indent int, items []*protogen.Message) {
	iterator(indent, len(items), "Messages", func(i int) {
		Message(indent+1, items[i])
	})
}

/* Structs */

func Enum(indent int, enum *protogen.Enum) {
	var vals []string
	for _, enumVal := range enum.Values {
		vals = append(vals,
			fmt.Sprintf("%s(%d)",
				enumVal.Desc.Name(), enumVal.Desc.Number()))
	}
	printTab(indent, "%s: %s\n", enum.Desc.Name(), strings.Join(vals, ", "))
}

func Service(indent int, s *protogen.Service) {
	printTab(indent, "%s: %s\n", s.Desc.Name(), s.GoName)
	indent += 1
	for _, m := range s.Methods {
		var inStream, outStream string
		if m.Desc.IsStreamingClient() {
			inStream = "stream "
		}
		if m.Desc.IsStreamingServer() {
			outStream = "stream "
		}
		printTab(indent, "%s:\t%s%s\t-> %s%s\n",
			m.Desc.Name(),
			inStream, m.Input.Desc.Name(),
			outStream, m.Output.Desc.Name())
	}
}

func Message(indent int, msg *protogen.Message) {
	if msg.Desc.IsMapEntry() {
		printTab(indent, "%s: IsMapEntry\n", msg.Desc.Name())
		return
	}

	printTab(indent, "%s\n", msg.Desc.Name())
	Fields(indent+1, msg.Fields)
	OneOfs(indent+1, msg.Oneofs)
	Enums(indent+1, msg.Enums)
	Messages(indent+1, msg.Messages)
	Extensions(indent+1, msg.Extensions)
}

func Field(indent int, f *protogen.Field) {
	fd := f.Desc
	var pre, post string
	if fd.IsList() {
		pre = "repeated "
	} else if fd.IsMap() {
		post = fmt.Sprintf(", map<%s, %s>",
			fieldType(fd.MapKey()), fieldType(fd.MapValue()))
	}

	printTab(indent, "%s%s(%d): %s%s\n", pre, fd.Name(), fd.Number(),
		f.GoName, post)
}

func fieldType(fd protoreflect.FieldDescriptor) string {
	switch k := fd.Kind(); k {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return string(fd.Message().Name())

	default:
		return k.String()
	}
}

func Oneof(indent int, o *protogen.Oneof) {
	extra := ""
	if o.Desc.IsSynthetic() {
		extra = ", proto3-optional"
	}

	printTab(indent, "%s: %s%s\n", o.Desc.Name(), o.GoName, extra)
	if !o.Desc.IsSynthetic() {
		Fields(indent+1, o.Fields)
	}
}
