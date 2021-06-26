package step_test

import (
	"testing"

	"github.com/draganm/senfgurke/step"
	"github.com/stretchr/testify/require"
)

func TestMatcher(t *testing.T) {

	type testCase struct {
		text        string
		shouldMatch bool
		params      []interface{}
	}

	cases := []struct {
		pattern string
		cases   []testCase
	}{
		{
			pattern: "foo",
			cases: []testCase{
				{
					text:        "foo",
					shouldMatch: true,
				},
				{
					text:        "fooz",
					shouldMatch: false,
				},
			},
		},
		{
			pattern: "foo {int}",
			cases: []testCase{
				{
					text:        "foo 1",
					shouldMatch: true,
					params:      []interface{}{1},
				},
				{
					text:        "fooz",
					shouldMatch: false,
				},
			},
		},

		{
			pattern: "foo {int} bar {int}",
			cases: []testCase{
				{
					text:        "foo 1 bar 2",
					shouldMatch: true,
					params:      []interface{}{1, 2},
				},
			},
		},
	}

	for _, c := range cases {

		t.Run(c.pattern, func(t *testing.T) {
			for _, ca := range c.cases {
				t.Run(ca.text, func(t *testing.T) {
					p, err := step.Match(c.pattern, ca.text)
					if !ca.shouldMatch {
						require.Error(t, err)
						return
					} else {
						require.NoError(t, err)
					}
					if ca.params == nil {
						ca.params = []interface{}{}
					}
					require.Equal(t, ca.params, p)
				})
			}
		})
	}
}
