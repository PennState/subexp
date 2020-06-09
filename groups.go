package subexp

import "regexp"

/*
Groups contains the sub-expressions "captured" when matching text has
been found during the evaluation against a regular expression.  While
none of the fields are directly available, the individual groups can
be accessed by name or index.
*/
type Groups struct {
	s  []string
	lu map[string]int
}

/*
Capture attempts to find the pattern defined by the regular expression
in the text paramater.  If there are one or more capture groups, a
Groups object is returned.  If there are no matched capture groups, this
method returns nil.
*/
func Capture(re *regexp.Regexp, text string) *Groups {
	s := re.FindStringSubmatch(text)
	if len(s) == 0 {
		return nil
	}

	lu := map[string]int{}
	for i, n := range re.SubexpNames() {
		lu[n] = i
	}

	return &Groups{
		s:  s,
		lu: lu,
	}
}

func compileAndCapture(expr, text string, fn func(string) (*regexp.Regexp, error)) (*Groups, error) {
	re, err := fn(expr)
	if err != nil {
		return nil, err
	}

	return Capture(re, text), nil
}

/*
CompileAndCapture is a convenience method for performing the the
regular expression compilation and capture in a single step.  If the
regular expression fails to compile, an error is returned, otherwise the
Groups object will be valid or nil as described by the Capture method.
*/
func CompileAndCapture(expr, text string) (*Groups, error) {
	return compileAndCapture(expr, text, regexp.Compile)
}

/*
CompileAndCapturePOSIX is a convenience method with the same behavior
as CompileAndCapture but uses POSIX syntax and evaluation.
*/
func CompileAndCapturePOSIX(expr, text string) (*Groups, error) {
	return compileAndCapture(expr, text, regexp.CompilePOSIX)
}

/*
MustCompileAndCapture performs the same function as CompileAndCapture
but panics if the regular expression cannot be compiled instead of
returning an error.
*/
func MustCompileAndCapture(expr, text string) *Groups {
	return Capture(regexp.MustCompile(expr), text)
}

/*
MustCompileAndCapturePOSIX is a convenience method with the same behavior
as CompileAndCapturePOSIX except it panics instead of returning an error
if the regular expression can not be compiled.
*/
func MustCompileAndCapturePOSIX(expr, text string) *Groups {
	return Capture(regexp.MustCompilePOSIX(expr), text)
}

/*
ByIndex returns the captured group (string) by index counting from left
to right.  By conventsion, index zero represents the entire matched text, so
index one would be the first captured group.
*/
func (g Groups) ByIndex(i int) string {
	// TODO: add error checking
	return g.s[i]
}

/*
ByName returns the text matching a named capture group by name.  Due to
the evaluation of the captured groups, this method currently returns the
last capture group with the provide name if the name is repeated.  This
behavior will be changing shortly.
*/
func (g Groups) ByName(name string) string {
	// TODO: add error checking
	// TODO: split into FirstByName and AllByName
	return g.ByIndex(g.lu[name])
}
