package elmgen

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func protoIdentToElmCasing(fullIdent string) []string {
	var idents [][][]rune
	for _, ident := range strings.Split(fullIdent, ".") {
		var words [][]rune
		appendWord := func(add []rune) {
			// Prefix empty segments with "X"
			if len(add) == 0 {
				add = []rune{'X'}
			}
			words = append(words, add)
		}
		var buf []rune
		var caps, underscore bool
		// Break ident into []words on underscores (__ counts as one) and runs of caps (URLTag is Url, Tag).
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
		// Add result to idents
		idents = append(idents, words)
	}
	// Build Elm ID from our words
	var idents2 []string
	for _, words := range idents {
		var words2 []string
		// Uppercase first letter of each word
		for _, word := range words {
			first := unicode.ToUpper(word[0])
			cased := string(first) + strings.ToLower(string(word[1:]))
			words2 = append(words2, cased)
		}
		// Stitch together
		ident := strings.Join(words2, "")
		idents2 = append(idents2, ident)
	}
	return idents2
}

func protoPkgToElmModule(pkg string) string {
	var parts []string
	for _, part := range protoIdentToElmCasing(pkg) {
		if !validElmID(part) {
			part = "X" + part
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, ".")
}

func protoIdentToElmID(ident string) (asType, asValue string) {
	parts := protoIdentToElmCasing(ident)
	asType = strings.Join(parts, "_")
	// Lowercase first rune
	runes := []rune(asType)
	asValue = strings.ToLower(string(runes[:1])) + string(runes[1:])
	// Valid Elm?
	if !validElmID(asType) || !validElmID(asValue) ||
		reservedWord(asType) || reservedWord(asValue) {
		asType = "X" + asType
		asValue = "x" + asValue
	}
	return
}

func protoIdentToElmValue(ident string) string {
	_, id := protoIdentToElmID(ident)
	return id
}

/* Takes an ident from protoreflect and converts to an Elm ID */

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

func protoReflectToElm(p Packager, d FullNamer) (mod, asType, asValue string) {
	pkg, fullIdent := string(p.Package()), string(d.FullName())
	mod = protoPkgToElmModule(pkg)
	postPkg := strings.TrimPrefix(fullIdent, pkg+".")
	asType, asValue = protoIdentToElmID(postPkg)
	return
}

/* Helper naming functions */

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

func reservedWord(check string) bool {
	for _, word := range reservedWords {
		if word == check {
			return true
		}
	}
	return false
}
