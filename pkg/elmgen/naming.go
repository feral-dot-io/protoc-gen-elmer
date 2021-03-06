// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
package elmgen

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Takes an proto full ident (dot (.) separated idents) and returns a list of Elm IDs.
// Each ID is a valid Elm ID, non-empty and won't include an underscore. It is formatted to loosely follow Elm naming conventions: a capital on start of words after an underscore or run of caps with the rest being lowercase.
// Examples: `my.pkg` -> `My, Pkg` and `My.URLIs_Here` -> `My, UrlIsHere`
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

// Converts a proto package to an Elm module
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

// Converts a proto ident to an Elm type and value
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

// Converts a proto ident to just an Elm value (convenience fn)
func protoIdentToElmValue(ident string) string {
	_, id := protoIdentToElmID(ident)
	return id
}

// Checks an Elm ID is valid. Does not check for reserved words
func validElmID(id string) bool {
	runes := []rune(id)
	return utf8.ValidString(id) && id != "" && // Non-empty utf8
		unicode.IsLetter(runes[0]) && // First char is a letter
		validPartialElmID(string(runes[1:])) // Remaining chars are valid
}

// Checks whether the ID is valid if it wasn't the first character
func validPartialElmID(partial string) bool {
	for _, r := range partial {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	// Allow empty as well
	return true
}

// List of reserved words. Should be reviewed with major Elm versions
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

// Checks whether an ID is a reserved word in Elm
func reservedWord(id string) bool {
	for _, word := range reservedWords {
		if word == id {
			return true
		}
	}
	return false
}
