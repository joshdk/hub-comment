// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package hub

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	tests := []struct {
		title    string
		body     string
		expected string
	}{
		{
			title: "empty string",
		},

		{
			title: "only spaces",
			body:  "    ",
		},
		{
			title: "only newlines",
			body:  "\n\n\n\n",
		},
		{
			title: "spaces with leading newlines",
			body:  "\n\n\n\n    ",
		},
		{
			title: "newlines with leading spaces",
			body:  "    \n\n\n\n",
		},
		{
			title:    "single line",
			body:     "hello world",
			expected: "hello world",
		},
		{
			title:    "single line with leading spaces",
			body:     "    hello world",
			expected: "    hello world",
		},
		{
			title:    "single line with trailing spaces",
			body:     "hello world    ",
			expected: "hello world",
		},
		{
			title:    "single line with surrounding spaces",
			body:     "    hello world    ",
			expected: "    hello world",
		},
		{
			title:    "single line with leading newlines",
			body:     "\n\nhello world",
			expected: "hello world",
		},
		{
			title:    "single line with leading newlines and spaces",
			body:     "    \n    \n    hello world",
			expected: "    hello world",
		},
		{
			title:    "single line with trailing newlines",
			body:     "hello world\n\n",
			expected: "hello world",
		},
		{
			title:    "single line with trailing newlines and spaces",
			body:     "hello world    \n    \n    ",
			expected: "hello world",
		},
		{
			title:    "single unicode",
			body:     "🤖",
			expected: "🤖",
		},
		{
			title:    "multiple unicode",
			body:     "    \n    \n    🤖    🤖    \n    \n    ",
			expected: "    🤖    🤖",
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("%d %s", index+1, test.title)
		t.Run(name, func(t *testing.T) {
			actual := trim(test.body)

			assert.False(t, strings.HasPrefix(actual, "\n"))
			assert.False(t, strings.HasSuffix(actual, "\n"))
			assert.Equal(t, test.expected, actual)
		})
	}
}
