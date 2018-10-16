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

	// Git is a map of GitHub specific parameters.
	Git map[string]string

	// Build is a map of CircleCI specific parameters.
	Build map[string]string
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

// Get is a helper function for looking up the given key in a map. If the key
// does not exist, fallback is returned, but only if specified.
func get(env map[string]string, key string, fallback ...string) string {
	value, found := env[key]
	switch {
	case found:
		return value
	case !found && len(fallback) > 0:
		return fallback[0]
	default:
		return ""
	}
}

// NewContext is a helper for constructing a context object.
func NewContext(environ []string) Context {
	env := makeEnv(environ)
	return Context{
		Env: env,
		Git: map[string]string{
			"Branch": get(env, "CIRCLE_BRANCH"),
			"PR":     get(env, "CIRCLE_PULL_REQUEST"),
			"SHA":    get(env, "CIRCLE_SHA1"),
			"Tag":    get(env, "CIRCLE_TAG"),
		},
		Build: map[string]string{
			"CI":       get(env, "CIRCLECI"),
			"Index":    get(env, "CIRCLE_NODE_INDEX", "0"),
			"Job":      get(env, "CIRCLE_JOB"),
			"Nodes":    get(env, "CIRCLE_NODE_TOTAL", "1"),
			"Number":   get(env, "CIRCLE_BUILD_NUM"),
			"Owner":    get(env, "CIRCLE_PROJECT_USERNAME"),
			"Repo":     get(env, "CIRCLE_PROJECT_REPONAME"),
			"Stage":    get(env, "CIRCLE_STAGE"),
			"URL":      get(env, "CIRCLE_BUILD_URL"),
			"User":     get(env, "CIRCLE_USERNAME"),
			"Workflow": get(env, "CIRCLE_WORKFLOW_ID"),
		},
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
