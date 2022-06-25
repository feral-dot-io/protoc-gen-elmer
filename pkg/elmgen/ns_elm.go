package elmgen

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"google.golang.org/protobuf/reflect/protoreflect"
)

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

func (m *Module) protoFullIdentToElmID(name protoreflect.FullName, isType bool) string {
	// Aliased?
	alias := m.protoAliases[name]
	if alias == "" {
		// No alias, drop pkg prefix from full
		alias = strings.TrimPrefix(string(name), string(m.protoPkg))
	}
	return protoFullIdentToElmCasing(alias, m.config.QualifiedSeparator, isType)
}

func (m *Module) getElmType(name protoreflect.FullName) (ElmType, error) {
	elmID, ok := m.protoNS[name]
	// This should never happen. All proto names should be registered before use
	if !ok {
		log.Panicf("missing protoreflect.FullName: %s", name)
	}
	// First retrieval: create
	var err error
	if elmID == "" {
		candidate := m.protoFullIdentToElmID(name, true)
		candidate, err = m.registerElmID(candidate)
		elmID = ElmType(candidate)
		m.protoNS[name] = elmID
	}
	return elmID, err
}

func (m *Module) getElmValue(name protoreflect.FullName) string {
	return m.protoFullIdentToElmID(name, false)
}

func (m *Module) registerElmID(id string) (string, error) {
	if !validElmID(id) {
		log.Panicf("invalid Elm ID: %s", id)
	}
	// Already registered?
	if _, ok := m.elmNS[id]; ok {
		// Generate an error?
		if m.config.CollisionSuffix == "" {
			return "", fmt.Errorf("protobuf schema generates a name collision with ID `%s` and current config (%#v) prevents us from resolving it",
				id, m.config)
		}
		// Add suffix to ID
		id += m.config.CollisionSuffix
		return m.registerElmID(id)
	}
	m.elmNS[id] = struct{}{}
	return id, nil
}

func (d *CodecIDs) register(m *Module, name protoreflect.FullName) (err error) {
	if d.ID, err = m.getElmType(name); err != nil {
		return err
	}
	if d.ZeroID, err = m.registerElmID("empty" + string(d.ID)); err != nil {
		return err
	}
	id := m.getElmValue(name)
	if d.DecodeID, err = m.registerElmID(id + "Decoder"); err != nil {
		return err
	}
	d.EncodeID, err = m.registerElmID(id + "Encoder")
	if err != nil {
		return err
	}
	d.FuzzerID, err = m.registerElmID(id + "Fuzzer")
	return err
}
