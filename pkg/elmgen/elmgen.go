package elmgen

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	Module struct {
		Name, Path string
		Imports    struct {
			Bytes bool
			Dict  bool
		}
		Fuzzers struct {
			Int32   bool
			Uint32  bool
			Float32 bool
		}

		Unions  Unions
		Oneofs  Oneofs
		Records Records

		RPCs RPCs
	}

	ElmRef struct {
		Module, ID string
	}
	ElmType struct {
		ElmRef
		asValue string
	}

	// Unions are sortable by their Elm ID.
	Unions []*Union
	// Union is a sum type of simple tags (i.e., no other data). Has a default tag holding an unknown value.
	Union struct {
		Type           *ElmType
		DefaultVariant *Variant
		Variants       []*Variant
		Aliases        []*VariantAlias
	}
	// Variant represents an enum option.
	Variant struct {
		ID     *ElmRef
		Number protoreflect.EnumNumber
	}
	// VariantAlias is a Variant with an alternative name. Identified by having the same wire number. First Variant seen is the real one, subsequent are alternate names.
	VariantAlias struct {
		*Variant
		Alias *ElmRef
	}

	// Oneofs are sortable by their Elm ID.
	Oneofs []*Oneof
	// Oneof is a sum type whose tags (variants) hold complex data. Oneofs are held by records as "one and only one from a selection of fields".
	Oneof struct {
		Type        *ElmType
		IsSynthetic bool
		Variants    []*OneofVariant
	}
	// Like a union's Variant but holds a field
	OneofVariant struct {
		ID    *ElmRef // Promoted Field label
		Field *Field
	}

	// Records are sortable by their Elm ID
	Records []*Record
	Record  struct {
		Type   *ElmType
		Oneofs []*Oneof
		Fields []*Field
	}

	// Represents an Elm record field
	Field struct {
		Label string
		// Wire handling
		IsOneof     bool
		IsMap       bool
		WireNumber  protowire.Number
		Cardinality protoreflect.Cardinality
		// Elm handling
		Type             string
		Zero             string
		Decoder, Encoder string
		Fuzzer           string

		Key *MapKey
	}

	MapKey struct {
		Zero             string
		Decoder, Encoder string
		Fuzzer           string
	}

	RPCs []*RPC
	RPC  struct {
		ID      *ElmRef
		In, Out *ElmType

		InStreaming, OutStreaming bool

		Service protoreflect.FullName
		Method  protoreflect.Name
	}
)

func (a Unions) Len() int           { return len(a) }
func (a Unions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Unions) Less(i, j int) bool { return a[i].Type.String() < a[j].Type.String() }

func (a Oneofs) Len() int           { return len(a) }
func (a Oneofs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Oneofs) Less(i, j int) bool { return a[i].Type.String() < a[j].Type.String() }

func (a Records) Len() int           { return len(a) }
func (a Records) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Records) Less(i, j int) bool { return a[i].Type.String() < a[j].Type.String() }

func (a RPCs) Len() int           { return len(a) }
func (a RPCs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RPCs) Less(i, j int) bool { return a[i].ID.String() < a[j].ID.String() }

func NewModule(prefix string, fd protoreflect.FileDescriptor) *Module {
	m := new(Module)
	// Paths
	pkg := prefix + string(fd.Package())
	m.Name = protoFullIdentToElmCasing(pkg, ".", true)
	m.Path = protoFullIdentToElmCasing(pkg, "/", true)
	// Parse file
	m.addUnions(fd.Enums())
	m.addRecords(fd.Messages())
	m.addRPCs(fd.Services())
	return m
}
