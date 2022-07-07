package elmgen

import (
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	// Top-level structures describing an Elm Module
	Module struct {
		importsSeen map[string]bool

		Name, Path, Proto string
		Imports           []string

		Unions   Unions
		Oneofs   Oneofs
		Records  Records
		Services Services
	}

	// Elm reference pointing an identifier e.g., a type or function in another module. Module is blank for local references.
	ElmRef struct {
		Module, ID string
	}
	// Speciality reference for a codegen type that has our derived functions. Refs are never nil. The fuzzer assumes a reference to another module with the "Tests" suffix (see `NewModule`)
	ElmType struct {
		*ElmRef
		Zero, Decoder, Encoder, Fuzzer *ElmRef
	}

	// Describes a set of comments from the Protobuf source
	CommentSet struct {
		LeadingDetached []Comments
		Leading         Comments
		Trailing        Comments
	}
	// Specialised Elm comments. Has a stringer method that wraps itself in dash comments (--)
	Comments string

	// Unions sortable by ID
	Unions []*Union
	// Union is a sum type with simple tags that won't hold data. Always holds at least one variant with the first being the default. Derived from Protobuf enums.
	Union struct {
		Type     *ElmType
		Variants []*Variant
		Aliases  []*VariantAlias
		Comments *CommentSet
	}
	// Describes a Union tag
	Variant struct {
		ID       *ElmRef
		Number   protoreflect.EnumNumber
		Comments *CommentSet
	}
	// VariantAlias is a Variant with an alternative name and the same wire number. First Variant seen is the real one, subsequent are alternate names.
	VariantAlias struct {
		Alias    *ElmRef
		Variant  *Variant
		Comments *CommentSet
	}

	// Oneofs sortable by ID
	Oneofs []*Oneof
	// Oneof is a sum type, similar to a Union, whose tags hold a reference to another data type. The data type can be complex (another type like a record or union) or simple (scalar data like a bool). Oneofs are used in records where a field should hold "one and only one" from a set of types.
	// If IsSynthetic is set then it's an optional field with a single choice for its data type which can be present or missing.
	Oneof struct {
		Type        *ElmType
		IsSynthetic bool
		Variants    []*OneofVariant
		Comments    *CommentSet
	}
	// Describes a Oneof option. Very similar to Variant except it holds a data type as well.
	OneofVariant struct {
		ID    *ElmRef // Promoted Field label
		Field *Field
	}

	// Records sortable by ID
	Records []*Record
	// A record is derived from a Protobuf message
	Record struct {
		Type     *ElmType
		Oneofs   []*Oneof
		Fields   []*Field
		Comments *CommentSet
	}

	// A record field. Desc may be nil if it's a non-synthetic Oneof.
	Field struct {
		Label    string
		Desc     protoreflect.FieldDescriptor
		Oneof    *Oneof
		Comments *CommentSet
	}

	// Services sortable by ID
	Services []*Service
	// Represents a grouping of RPC methods. Not necessarily important but used to retain comments
	Service struct {
		Label    string
		Methods  RPCs
		Comments *CommentSet
	}

	// RPC sortable by ID
	RPCs []*RPC
	// Describes an RPC method
	RPC struct {
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

// Entry point for elmgen. Builds an Elm module from a given proto File. The module name may be suffixed to allow for different derivative use cases e.g., a codec with no suffix and the suffix "Twirp" for a client could live alongside each other.
func NewModule(suffix string, file *protogen.File) *Module {
	m := new(Module)
	m.importsSeen = make(map[string]bool)
	// Paths
	pkg := string(file.Desc.Package())
	// Adding a prefix / suffix can prevent an empty can lead to X.elm and Tests.elm
	// We want a consistent form i.e. X.elm and XTests.elm
	if pkg == "" {
		pkg = "X"
	}
	pkg = pkg + suffix
	m.Name = protoPkgToElmModule(pkg)
	m.Path = strings.ReplaceAll(m.Name, ".", "/") + ".elm"
	m.Proto = file.Desc.Path()
	// Parse file
	m.addUnions(file.Enums)
	m.addRecords(file.Messages)
	m.addRPCs(file.Services)
	// Imports
	m.findImports()
	for key := range m.importsSeen {
		m.Imports = append(m.Imports, key)
	}
	sort.Strings(m.Imports)
	return m
}

// Converts a protogen CommentSet to our own
func newCommentSet(set protogen.CommentSet) *CommentSet {
	out := new(CommentSet)
	for _, c := range set.LeadingDetached {
		out.LeadingDetached = append(out.LeadingDetached,
			Comments(c))
	}
	out.Leading = Comments(set.Leading)
	out.Trailing = Comments(set.Trailing)
	return out
}
