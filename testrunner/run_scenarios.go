package testrunner

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/gherkin-go"
	"github.com/draganm/senfgurke/step"
	"github.com/draganm/senfgurke/world"
	"github.com/stretchr/testify/require"
)

func RunScenarios(t *testing.T, steps *step.Registry) {
	entries, err := os.ReadDir(".")
	require.NoError(t, err)

	features := []string{}
	for _, e := range entries {

		if !strings.HasSuffix(e.Name(), ".feature") {
			continue
		}

		i, err := e.Info()
		require.NoError(t, err)
		if !i.Mode().IsRegular() {
			continue
		}

		features = append(features, e.Name())
	}

	for _, f := range features {
		f := f
		doc, err := parseGherkin(f)
		require.NoError(t, err)

		t.Run(fmt.Sprintf("%s(%s)", f, doc.Feature.Name), func(t *testing.T) {
			for _, p := range doc.Pickles() {
				p := p
				t.Run(fmt.Sprintf("Scenario: %s", p.Name), func(t *testing.T) {
					w := world.World{}
					for _, s := range p.Steps {
						err = steps.Execute(s.Text, w)
						require.NoError(t, err)
					}
				})
			}

		})
	}

}

func parseGherkin(name string) (gherkinDocument *gherkin.GherkinDocument, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return gherkin.ParseGherkinDocument(f)

}
