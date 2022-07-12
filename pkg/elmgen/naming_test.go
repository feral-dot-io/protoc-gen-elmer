// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
package elmgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestProtoIdentToElmCasing(t *testing.T) {
	cases := map[string][]string{
		"hello":                         {"Hello"},
		"hello_world":                   {"HelloWorld"},
		"hello.world":                   {"Hello", "World"},
		"pkg.name.MyMessage.field_name": {"Pkg", "Name", "MyMessage", "FieldName"},
		"ALL_CAPS":                      {"AllCaps"},
		"ALL__CAPS":                     {"AllCaps"},
		"ALL.CAPS":                      {"All", "Caps"},
		"ALL..CAPS":                     {"All", "X", "Caps"},
		"TT":                            {"Tt"},
		"TTaaa":                         {"TTaaa"},
		"_hello.1hello":                 {"XHello", "1hello"},
		"_Hello":                        {"XHello"},
		"__Hello":                       {"XHello"},
		"Hello_":                        {"HelloX"},
		"Hello__":                       {"HelloX"},
		"__Hello__":                     {"XHelloX"},
		"_":                             {"XX"},
		"___":                           {"XX"},
		"URL":                           {"Url"},
		"URLTag":                        {"UrlTag"},
		"URL1Tag":                       {"Url1Tag"},
		"A_B_C":                         {"ABC"}, // Looks odd
		"MyURLIsHere":                   {"MyUrlIsHere"},
		"My.URLIs_Here":                 {"My", "UrlIsHere"},
		"My_URL_Is_Here":                {"MyUrlIsHere"},
		"UpUpUp":                        {"UpUpUp"},
		".":                             {"X", "X"},
		"":                              {"X"},
		"...":                           {"X", "X", "X", "X"},
		"oops.oops":                     {"Oops", "Oops"},
		"my._pkg":                       {"My", "XPkg"},
		"1andonly":                      {"1andonly"},
	}
	for check, exp := range cases {
		act := protoIdentToElmCasing(check)
		assert.Equal(t, exp, act, "check=%s", check)
	}
}

func TestProtoPkgToElmModule(t *testing.T) {
	cases := map[string]string{
		"hello":       "Hello",
		"helloWorld":  "HelloWorld",
		"hello.world": "Hello.World",
		"hello.1":     "Hello.X1",
		"..":          "X.X.X",
		"type.Int":    "Type.Int",
	}
	for check, exp := range cases {
		act := protoPkgToElmModule(check)
		assert.Equal(t, exp, act, "check=%s", check)
	}
}

func TestProtoIdentToElmID(t *testing.T) {
	cases := map[string][]string{
		"hello":       {"Hello", "hello"},
		"helloWorld":  {"HelloWorld", "helloWorld"},
		"hello.world": {"Hello_World", "hello_World"},
		"hello.1":     {"Hello_1", "hello_1"},
		"..":          {"X_X_X", "x_X_X"},
		"type.Int":    {"Type_Int", "type_Int"},
		"type":        {"XType", "xtype"},
		"to.Float":    {"To_Float", "to_Float"},
	}
	for check, exp := range cases {
		actType, actVal := protoIdentToElmID(check)
		assert.Equal(t, exp[0], actType, "check=%s", check)
		assert.Equal(t, exp[1], actVal, "check=%s", check)
	}
}

type protoTest struct {
	pkg, name string
}

func (t *protoTest) Package() protoreflect.FullName {
	return protoreflect.FullName(t.pkg)
}

func (t *protoTest) FullName() protoreflect.FullName {
	return protoreflect.FullName(t.name)
}

func TestProtoToElm(t *testing.T) {
	ex := &protoTest{"package.name", "full.name"}
	mod, asType, asVal := protoReflectToElm(ex, ex)
	assert.Equal(t, "Package.Name", mod)
	assert.Equal(t, "Full_Name", asType)
	assert.Equal(t, "full_Name", asVal)
	// Reserved word collision
	ex.name = "case"
	mod, asType, asVal = protoReflectToElm(ex, ex)
	assert.Equal(t, "Package.Name", mod)
	assert.Equal(t, "XCase", asType)
	assert.Equal(t, "xcase", asVal)
}

func TestValidElmID(t *testing.T) {
	assert.True(t, validElmID("Hello"))
	assert.True(t, validElmID("HelloWorld"))
	assert.False(t, validElmID("_Hello"))
	assert.True(t, validElmID("Hello_World"))
	assert.True(t, validElmID("type"))
	// Partial
	assert.True(t, validPartialElmID("Hello123_"))
	assert.True(t, validPartialElmID("_Hello"))
	assert.False(t, validPartialElmID("Hello?"))
}
