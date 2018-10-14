// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package hub

import (
	"bytes"
	"strings"
	"text/template"
)

// Context represents a logical grouping of data for use with comment templates.
type Context struct {

	// Env is a map of available environment variables.
	Env map[string]string
}

// makeEnv takes in a list of strings of the form "key=value", and returns a
// map of keys to their respective values. Intended to be passed the return
// value of os.Environ().
func makeEnv(environ []string) map[string]string {
	env := make(map[string]string, len(environ))

	for _, entry := range environ {
		pieces := strings.SplitN(entry, "=", 2)
		name, value := pieces[0], pieces[1]
		env[name] = value
	}
	return env
}

// NewContext is a helper for constructing a context object.
func NewContext(environ []string) Context {
	return Context{
		Env: makeEnv(environ),
	}
}

// NewTemplate is a helper for constructing a template object.
func NewTemplate(body []byte) (*template.Template, error) {
	return template.New("comment").Parse(string(body))
}

// Execute applies the given context to the given template and returns the
// result as a string.
func Execute(tpl *template.Template, ctx Context) (string, error) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}
