package elmgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestNaming(t *testing.T) {
	cases := map[string]string{
		"hello":                         "Hello",
		"hello_world":                   "HelloWorld",
		"hello.world":                   "HelloWorld",
		"pkg.name.MyMessage.field_name": "PkgNameMyMessageFieldName",
		"ALL_CAPS":                      "AllCaps",
		"ALL__CAPS":                     "AllCaps",
		"ALL.CAPS":                      "AllCaps",
		"ALL..CAPS":                     "AllXCaps",
		"TT":                            "Tt",
		"TTaaa":                         "TTaaa",
		"_hello.1hello":                 "XHello1hello",
		"_Hello":                        "XHello",
		"__Hello":                       "XHello",
		"Hello_":                        "HelloX",
		"Hello__":                       "HelloX",
		"__Hello__":                     "XHelloX",
		"_":                             "XX",
		"___":                           "XX",
		"URL":                           "Url",
		"URLTag":                        "UrlTag",
		"URL1Tag":                       "Url1Tag",
		"A_B_C":                         "ABC", // Looks odd
		"MyURLIsHere":                   "MyUrlIsHere",
		"My_URL_Is_Here":                "MyUrlIsHere",
		"UpUpUp":                        "UpUpUp",
		".":                             "XX",
		"":                              "X",
		"...":                           "XXXX",
		"oops.oops":                     "OopsOops",
		"my._pkg":                       "MyXPkg",
		"1andonly":                      "X1andonly",
	}
	for check, exp := range cases {
		assert.Equal(t, exp, protoFullIdentToElmCasing(check, "", true), "check=%s", check)
	}
	// Again but with a NS separator
	cases = map[string]string{
		"hello.world":                   "Hello_World",
		"pkg.name.MyMessage.field_name": "Pkg_Name_MyMessage_FieldName",
		"shadow":                        "Shadow",
		"_":                             "XX",
		".":                             "X_X",
		"...":                           "X_X_X_X",
		"1andonly":                      "X1andonly",
	}
	for check, exp := range cases {
		assert.Equal(t, exp, protoFullIdentToElmCasing(check, "_", true), "check=%s", check)
	}
	// Again but with type / value treatment
	cases = map[string]string{
		"hello.world":                   "helloWorld",
		"pkg.name.MyMessage.field_name": "pkgNameMyMessageFieldName",
		"":                              "x",
		"_":                             "xX",
		".":                             "xX",
		"...":                           "xXXX",
		"shadow":                        "shadow",
		"Outer":                         "outer",
		"_Outer":                        "xOuter",
		"Outer.Inner":                   "outerInner",
		"1andonly":                      "x1andonly",
	}
	for check, exp := range cases {
		assert.Equal(t, exp, protoFullIdentToElmCasing(check, "", false), "check=%s", check)
	}
}

func TestNS(t *testing.T) {
	m := TestConfig.newModule()
	// Trying to register duplicate proto ident
	m.registerProtoName(protoreflect.FullName("test"), "")
	assert.Panics(t, func() {
		m.registerProtoName(protoreflect.FullName("test"), "")
	})
	// Getting an Elm ID from an unregistered proto ident
	assert.Panics(t, func() {
		m.getElmType("eek")
	})
	// Accidental codec registration gives a collision error
	m = Config{}.newModule()
	m.registerProtoName("Dupe", "")
	cases := []string{"Dupe", "emptyDupe", "dupeDecoder", "dupeEncoder"}
	for _, id := range cases {
		m.registerElmID(id)
		err := new(CodecIDs).register(m, "Dupe")
		assert.ErrorContains(t, err, id)
		// Unregister values
		for _, unreg := range cases {
			delete(m.elmNS, unreg)
		}
	}
}

func TestCollisionSuffix(t *testing.T) {
	m := Config{CollisionSuffix: "___"}.newModule()
	// Duplicate Elm ID gets a suffixed
	id, err := m.registerElmID("Hello")
	assert.NoError(t, err)
	assert.Equal(t, "Hello", id)
	id, err = m.registerElmID("Hello")
	assert.NoError(t, err)
	assert.Equal(t, "Hello___", id)
	// Invalid suffix
	_, err = (&Config{CollisionSuffix: " "}).NewModule(nil)
	assert.ErrorContains(t, err, "collision")
	// Must be valid Elm
	assert.Panics(t, func() {
		m.registerElmID(" ")
	})
	// Empty suffix gives an error
	m = Config{}.newModule()
	_, err = m.registerElmID("getReady")
	assert.NoError(t, err)
	_, err = m.registerElmID("getReady")
	assert.ErrorContains(t, err, "collision")
}
