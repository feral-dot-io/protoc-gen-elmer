package elmgen

import (
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Config struct {
	// Literal value to prefix our generated module's name
	ModulePrefix string
	// Override our generated module's name with this value. Does not apply `ModulePrefix`
	ModuleName string
	// Proto allows nesting. When we decide on a name, do we use the last (deepest) name or qualify all?
	QualifyNested bool
	// Replaces a Protobuf namespace separator (a dot, `.`) with this string in Elm code
	QualifiedSeparator string
	// Variants created from enums and oneofs are suffixed with their parent
	VariantSuffixes bool
	// When mapping Protobuf idents to Elm we may up with a naming collision. Resolve this collision by appending this to the second ID generated or returning an error if blank
	CollisionSuffix string // TODO return an error on blank
}

var DefaultConfig = Config{
	VariantSuffixes: true,
	CollisionSuffix: "_"} // Must be non-empty

type (
	Module struct {
		config Config

		protoPkg      protoreflect.FullName
		protoNS       map[protoreflect.FullName]ElmType
		protoAliases  map[protoreflect.FullName]string
		protoEnums    []*protogen.Enum
		protoMessages []*protogen.Message
		elmNS         map[string]struct{} // Top-level types and functions

		Name, Path string
		Imports    struct {
			Bytes bool
			Dict  bool
		}

		Unions  Unions
		Oneofs  Oneofs
		Records Records
	}

	ElmType  string
	CodecIDs struct {
		ID                         ElmType
		ZeroID, DecodeID, EncodeID string
	}

	// Unions are sortable by their Elm ID.
	Unions []*Union
	// Union is a sum type of simple tags (i.e., no other data). Has a default tag holding an unknown value.
	Union struct {
		CodecIDs
		DefaultVariant *Variant
		Variants       []*Variant
		Aliases        []*VariantAlias
	}
	// Variant represents an enum option.
	Variant struct {
		ID     ElmType
		Number protoreflect.EnumNumber
	}
	// VariantAlias is a Variant with an alternative name. Identified by having the same wire number. First Variant seen is the real one, subsequent are alternate names.
	VariantAlias struct {
		*Variant
		Alias string
	}

	// Oneofs are sortable by their Elm ID.
	Oneofs []*Oneof
	// Oneof is a sum type whose tags (variants) hold complex data. Oneofs are held by records as "one and only one from a selection of fields".
	Oneof struct {
		CodecIDs
		IsSynthetic bool
		Variants    []*OneofVariant
	}
	// Like a union's Variant but holds a field
	OneofVariant struct {
		ID    ElmType // Promoted Field label
		Field *Field
	}

	// Records are sortable by their Elm ID
	Records []*Record
	Record  struct {
		CodecIDs
		Oneofs []*Oneof
		Fields []*Field
	}

	// Represents an Elm record field
	Field struct {
		Label string
		// Wire handling
		IsOneof     bool
		WireNumber  protowire.Number
		Cardinality protoreflect.Cardinality
		// Elm handling
		Type             string
		Zero             interface{}
		Decoder, Encoder string
	}
)

func (a Unions) Len() int           { return len(a) }
func (a Unions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Unions) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (a Oneofs) Len() int           { return len(a) }
func (a Oneofs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Oneofs) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (a Records) Len() int           { return len(a) }
func (a Records) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Records) Less(i, j int) bool { return a[i].ID < a[j].ID }

func ValidElmID(id string) bool {
	runes := []rune(id)
	return utf8.ValidString(id) && id != "" && // Non-empty utf8
		unicode.IsLetter(runes[0]) && // First char is a letter
		ValidPartialElmID(string(runes[1:])) // Remaining chars are valid
}

func ValidPartialElmID(partial string) bool {
	for _, r := range partial {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	// Allow empty as well
	return true
}
