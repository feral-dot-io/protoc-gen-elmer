package elmgen

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	Module struct {
		ns       map[string][]*ElmRef
		nonLocal bool

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

		Unions   Unions
		Oneofs   Oneofs
		Records  Records
		Services Services
	}

	ElmRef struct {
		Module, ID string
	}
	ElmType struct {
		*ElmRef
		asValue string
	}

	CommentSet struct {
		LeadingDetached []Comments
		Leading         Comments
		Trailing        Comments
	}
	Comments string

	// Unions are sortable by their Elm ID.
	Unions []*Union
	// Union is a sum type of simple tags (i.e., no other data). Has a default tag holding an unknown value.
	Union struct {
		Type           *ElmType
		DefaultVariant *Variant
		Variants       []*Variant
		Aliases        []*VariantAlias
		Comments       *CommentSet
	}
	// Variant represents an enum option.
	Variant struct {
		ID       *ElmRef
		Number   protoreflect.EnumNumber
		Comments *CommentSet
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
		Comments    *CommentSet
	}
	// Like a union's Variant but holds a field
	OneofVariant struct {
		ID    *ElmRef // Promoted Field label
		Field *Field
	}

	// Records are sortable by their Elm ID
	Records []*Record
	Record  struct {
		Type     *ElmType
		Oneofs   []*Oneof
		Fields   []*Field
		Comments *CommentSet
	}

	// Represents an Elm record field
	Field struct {
		Label    string
		Comments *CommentSet
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

	Services []*Service
	Service  struct {
		Label    string
		Methods  RPCs
		Comments *CommentSet
	}
	RPCs []*RPC
	RPC  struct {
		ID      *ElmRef
		In, Out *ElmType

		InStreaming, OutStreaming bool

		Service  protoreflect.FullName
		Method   protoreflect.Name
		Comments *CommentSet
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

func (a Services) Len() int           { return len(a) }
func (a Services) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Services) Less(i, j int) bool { return a[i].Label < a[j].Label }

func (a RPCs) Len() int           { return len(a) }
func (a RPCs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RPCs) Less(i, j int) bool { return a[i].ID.String() < a[j].ID.String() }

func NewModule(prefix string, file *protogen.File) *Module {
	m := new(Module)
	m.ns = make(map[string][]*ElmRef)
	// Paths
	pkg := prefix + string(file.Desc.Package())
	m.Name = protoFullIdentToElmCasing(pkg, ".", true)
	m.Path = protoFullIdentToElmCasing(pkg, "/", true)
	// Parse file
	m.addUnions(file.Enums)
	m.addRecords(file.Messages)
	m.addRPCs(file.Services)
	return m
}

func (m *Module) SetRefLocality(local bool) {
	for _, refs := range m.ns {
		for _, ref := range refs {
			if local && ref.Module == m.Name {
				ref.Module = ""
			} else if !local && ref.Module == "" {
				ref.Module = m.Name
			}
		}
	}
	m.nonLocal = !local
}

func NewCommentSet(set protogen.CommentSet) *CommentSet {
	out := new(CommentSet)
	for _, c := range set.LeadingDetached {
		out.LeadingDetached = append(out.LeadingDetached,
			Comments(c))
	}
	out.Leading = Comments(set.Leading)
	out.Trailing = Comments(set.Trailing)
	return out
}
