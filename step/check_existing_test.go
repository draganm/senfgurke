package step_test

import (
	"errors"
	"testing"

	"github.com/draganm/senfgurke/step"
	"github.com/stretchr/testify/require"
)

func TestCheckExisting(t *testing.T) {

	testCases := []struct {
		text     string
		expected error
	}{
		{
			text:     "plain text",
			expected: errors.New("// could not find step matching \"plain text\".\n// to implement it, add following code:\nvar _ = steps.Then(\"plain text\", func(w world.World) error {\t\n\treturn errors.New(\"not yet implemented\")\n})\n"),
		},
		{
			text:     "plain text \"stringarg\"",
			expected: errors.New("// could not find step matching \"plain text \\\"stringarg\\\"\".\n// to implement it, add following code:\nvar _ = steps.Then(\"plain text {string}\", func(w world.World, arg0 string) error {\t\n\treturn errors.New(\"not yet implemented\")\n})\n"),
		},
		{
			text:     "plain text \"stringarg\" sometext",
			expected: errors.New("// could not find step matching \"plain text \\\"stringarg\\\" sometext\".\n// to implement it, add following code:\nvar _ = steps.Then(\"plain text {string} sometext\", func(w world.World, arg0 string) error {\t\n\treturn errors.New(\"not yet implemented\")\n})\n"),
		},
		{
			text:     "\"stringarg\"",
			expected: errors.New("// could not find step matching \"\\\"stringarg\\\"\".\n// to implement it, add following code:\nvar _ = steps.Then(\"{string}\", func(w world.World, arg0 string) error {\t\n\treturn errors.New(\"not yet implemented\")\n})\n"),
		},
		{
			text:     "\"stringarg\"abc",
			expected: errors.New("// could not find step matching \"\\\"stringarg\\\"abc\".\n// to implement it, add following code:\nvar _ = steps.Then(\"\\\"stringarg\\\"abc\", func(w world.World) error {\t\n\treturn errors.New(\"not yet implemented\")\n})\n"),
		},
		{
			text:     "123abc",
			expected: errors.New("// could not find step matching \"123abc\".\n// to implement it, add following code:\nvar _ = steps.Then(\"123abc\", func(w world.World) error {\t\n\treturn errors.New(\"not yet implemented\")\n})\n"),
		},
		{
			text:     "123 abc",
			expected: errors.New("// could not find step matching \"123 abc\".\n// to implement it, add following code:\nvar _ = steps.Then(\"{int} abc\", func(w world.World, arg0 int) error {\t\n\treturn errors.New(\"not yet implemented\")\n})\n"),
		},
	}

	r := step.NewRegistry()

	for _, tc := range testCases {
		t.Run(tc.text, func(t *testing.T) {
			err := r.CheckExisting(tc.text)
			require.Equal(t, tc.expected, err)
		})
	}

}
