package elmgen

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	pkgSeparator    = "."
	identSeparator  = "_"
	collisionSuffix = "_"
)

type (
	Packager interface {
		Package() protoreflect.FullName
	}

	FullNamer interface {
		FullName() protoreflect.FullName
	}

	Namer interface {
		Name() protoreflect.Name
	}
)

var reservedWords = []string{
	// Source: https://github.com/elm/compiler/blob/770071accf791e8171440709effe71e78a9ab37c/compiler/src/Parse/Variable.hs
	"if", "then", "else", "case", "of", "let", "in", "type", "module",
	"where", "import", "exposing", "as", "port",
	// Prelude https://package.elm-lang.org/packages/elm/core/latest/
	"Basics", "List", "Maybe", "Result", "String", "Char", "Tuple", "Debug",
	"Platform", "Cmd", "Sub",
	// Basics(..)
	"Int", "Float", "toFloat", "round", "floor", "ceiling", "truncate",
	"max", "min", "compare", "LT", "EQ", "GT", "Bool", "True", "False",
	"not", "xor", "modBy", "remainderBy", "negate", "abs", "clamp", "sqrt",
	"logBase", "e", "degrees", "radians", "turns", "pi", "cos", "sin",
	"tan", "acos", "asin", "atan", "atan2", "toPolar", "fromPolar", "isNaN",
	"isInfinite", "identity", "always", "Never", "never",
	// Other imports (note these overlap with prelude)
	"Just", "Nothing", "Ok", "Err", "Program"}

func protoToElm(p Packager, d FullNamer) (mod, asType, asValue string) {
	pkg, fullIdent := string(p.Package()), string(d.FullName())
	mod = protoFullIdentToElmCasing(pkg, pkgSeparator, true)
	postPkg := strings.TrimPrefix(fullIdent, pkg+pkgSeparator)
	// Naming collision with reserved word?
	var post string
	for _, word := range reservedWords {
		if word == postPkg {
			post = collisionSuffix
			break
		}
	}
	// Elm IDs
	asType = protoFullIdentToElmCasing(postPkg, identSeparator, true) + post
	asValue = protoFullIdentToElmCasing(postPkg, identSeparator, false) + post
	return
}

func (m *Module) newElmRef(mod, id string) *ElmRef {
	ref := &ElmRef{mod, id}
	if mod == m.Name { // Local?
		ref.Module = ""
	}
	m.addImport(ref.Module)
	return ref
}

func (m *Module) NewElmValue(p Packager, d FullNamer) *ElmRef {
	mod, _, asValue := protoToElm(p, d)
	return m.newElmRef(mod, asValue)
}

func (m *Module) NewElmType(p Packager, d FullNamer) *ElmType {
	mod, asType, asValue := protoToElm(p, d)
	testMod := mod + "Tests"
	// Replace well known with our own library
	if mod == importGooglePB {
		mod, testMod = importElmer, importElmerTest
	}
	return &ElmType{
		m.newElmRef(mod, asType),
		m.newElmRef(mod, "empty"+asType),
		m.newElmRef(mod, asValue+"Decoder"),
		m.newElmRef(mod, asValue+"Encoder"),
		m.newElmRef(testMod, asValue+"Fuzzer")}
}

func (r *ElmRef) String() string {
	if r.Module == "" { // Local ref
		return r.ID
	}
	return r.Module + "." + r.ID
}

/*
	Transforms a Protobuf fullIdent into an Elm ID. Aside from an optional namespace separator, attempts to follow Elm's naming conventions.
	Protobuf idents can be made up of alphanum and underscores (this includes the first character).
	Return an Elm type or value with isType. Can be used to convert an Elm type to value as well.
	Rules:
		1) Split fullIdent into []ident on "."
		2) Break each ident into []words on underscores (__ counts as one) and runs of caps (URLTag is Url, Tag).
			a) Empty words are prefixed wtih an "X".
			b) Prefix first character of the first word with "X" if invalid Elm.
		3) First letter of each word is uppercased. Except on the first word if !isType (is a value).
		4) Idents are recombined with an optional separator.
*/
func protoFullIdentToElmCasing(fullIdent, sep string, isType bool) string {
	var idents [][][]rune
	// (1)
	for _, ident := range strings.Split(fullIdent, ".") {
		var words [][]rune
		appendWord := func(add []rune) {
			// (2a)
			if len(add) == 0 {
				add = []rune{'X'}
			}
			// (2b)
			if len(idents)+len(words) == 0 && !validElmID(string(add)) {
				add = append([]rune{'X'}, add...)
			}
			words = append(words, add)
		}
		var buf []rune
		var caps, underscore bool
		// (2)
		for _, r := range ident {
			if r == '_' { // Start of underscores
				if !underscore {
					underscore, caps = true, false
					appendWord(buf)
					buf = nil
				}
			} else {
				underscore = false
				if unicode.IsUpper(r) {
					if !caps && len(buf) > 0 {
						appendWord(buf)
						buf = nil
					}
					caps = true
				} else if caps && !unicode.IsDigit(r) { // Moving from caps to non-caps
					caps = false
					// New word boundary started -1 rune ago
					lastRune := len(buf) - 1
					if prior := buf[:lastRune]; len(prior) > 0 {
						appendWord(prior)
					}
					buf = []rune{buf[lastRune]}
				}
				// Accumulate
				buf = append(buf, r)
			}
		}
		// Add leftover buffer
		appendWord(buf)
		// Add to idents
		idents = append(idents, words)
	}
	// (3)
	var idents2 []string
	for i, words := range idents {
		var words2 []string
		for j, word := range words {
			var first rune
			if i == 0 && j == 0 && !isType { // Building an Elm value on first char
				first = unicode.ToLower(word[0])
			} else {
				first = unicode.ToUpper(word[0])
			}
			cased := string(first) + strings.ToLower(string(word[1:]))
			words2 = append(words2, cased)
		}
		ident := strings.Join(words2, "")
		idents2 = append(idents2, ident)
	}
	// (4)
	return strings.Join(idents2, sep)
}

func validElmID(id string) bool {
	runes := []rune(id)
	return utf8.ValidString(id) && id != "" && // Non-empty utf8
		unicode.IsLetter(runes[0]) && // First char is a letter
		validPartialElmID(string(runes[1:])) // Remaining chars are valid
}

func validPartialElmID(partial string) bool {
	for _, r := range partial {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	// Allow empty as well
	return true
}
