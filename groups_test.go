package subexp_test

import (
	"testing"

	"github.com/PennState/subexp"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name     string
	expr     string
	text     string
	err      error
	panicked bool
}

type Permutator func([]struct{})

func Permutate(trgt []struct{}, fns []Permutator) {
	for _, fn := range fns {
		fn(trgt)
	}
}

func TestPermutator(t *testing.T) {
	assert.True(t, true)
	// fns := []Permutator{
	// 	func(s []struct{}) {
	// 		s = append(s, test{name: "A"})
	// 		s = append(s, test{name: "B"})
	// 	},
	// 	func(s []struct{}) {

	// 	},
	// }
}

// expr * POSIX * constructor * compiles * matches

// func tests() []test {
// 	return []test{
// 		{"Valid expression - no match", "", ""},
// 		{"Valid expression - 3 matches", "", ""},
// 	}
// }

func adaptPanicSignature(fn func(expr, text string) *subexp.Groups) func(expr, text string) (*subexp.Groups, error) {
	return func(expr, text string) (*subexp.Groups, error) {
		return fn(expr, text), nil
	}
}

func TestCompileAnCapture(t *testing.T) {
	tests := []struct {
		name   string
		expr   string
		text   string
		exp    map[string][]string
		panics bool
		errs   bool
		fn     func(expr, text string) (*subexp.Groups, error)
	}{
		{
			name: "CompileAndCapture - works",
			expr: "^(?P<a>[0-9]) (?P<a>[0-9]) (?P<b>[0-9])(?: (?P<c>[0-9]))?$",
			text: "0 1 2",
			exp: map[string][]string{
				"a": {"0", "1"},
				"b": {"2"},
			},
			fn: subexp.CompileAndCapture,
		},
		{
			name: "CompileAndCapture - no match",
			expr: "^(?P<a>[0-9]) (?P<a>[0-9]) (?P<b>[0-9])$",
			text: "a b c",
			exp:  nil,
			fn:   subexp.CompileAndCapture,
		},
		{
			name: "CompileAndCapture - errs",
			expr: "^(?P<a>[0-9]$", // unbalanced group parenthesis
			errs: true,
			fn:   subexp.CompileAndCapture,
		},
		{
			name: "CompileAndCapturePOSIX - errs",
			expr: "^([0-9]$", // unbalanced group parenthesis
			errs: true,
			fn:   subexp.CompileAndCapturePOSIX,
		},
		{
			name: "MustCompileAndCapture - works",
			expr: "^(?P<a>[0-9]) (?P<a>[0-9]) (?P<b>[0-9])(?: (?P<c>[0-9]))?$",
			text: "0 1 2",
			exp: map[string][]string{
				"a": {"0", "1"},
				"b": {"2"},
			},
			fn: subexp.CompileAndCapture,
		},
		{
			name:   "MustCompileAndCapture - panics",
			expr:   "^(?P<a>[0-9]$", // unbalanced group parenthesis
			panics: true,
			fn:     adaptPanicSignature(subexp.MustCompileAndCapture),
		},
		{
			name:   "MustCompileAndCapturePOSIX - panics",
			expr:   "^([0-9]$", // unbalanced group parenthesis
			panics: true,
			fn:     adaptPanicSignature(subexp.MustCompileAndCapturePOSIX),
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			if test.panics {
				assert.Panics(t, func() {
					g, err := test.fn(test.expr, test.text)
					assert.NoError(t, err)
					assert.Nil(t, g)
				})
				return
			}

			if test.errs {
				g, err := test.fn(test.expr, test.text)
				assert.Error(t, err)
				assert.Nil(t, g)
				return
			}

			g, err := test.fn(test.expr, test.text)
			assert.NoError(t, err)

			if test.exp == nil {
				assert.Nil(t, g)
				return
			}

			for k, v := range test.exp {
				a, err := g.AllByName(k) // AllByName works
				assert.NoError(t, err)
				assert.Equal(t, v, a)

				s, err := g.FirstByName(k) // FirstByName works
				assert.NoError(t, err)
				assert.Equal(t, v[0], s)
			}

			s, err := g.ByIndex(1)
			assert.NoError(t, err)
			assert.Equal(t, "0", s)

			s, err = g.ByIndex(2)
			assert.NoError(t, err)
			assert.Equal(t, "1", s)

			s, err = g.ByIndex(3)
			assert.NoError(t, err)
			assert.Equal(t, "2", s)

			s, err = g.ByIndex(4) // optional capture group
			assert.NoError(t, err)
			assert.Empty(t, s)

			s, err = g.ByIndex(5) // Invalid index
			assert.Error(t, err)
			assert.Empty(t, s)

			s, err = g.FirstByName("c") // optional capture group
			assert.Error(t, err)
			assert.Empty(t, s)

			s, err = g.FirstByName("d") // missing name
			assert.Error(t, err)
			assert.Empty(t, s)
		})
	}
}

// func TestNameUnknown(t *testing.T) {
// 	m := map[string][]string{
// 		"a": {"1", "2", "3"},
// 		"b": {"2", "3", "4"},
// 		"c": {"3", "4", "5"},
// 	}

// 	t.Log("UnknownName", m["d"])
// 	t.Fail()
// }
