package testrunner

import (
	"errors"
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

	// check for missing steps
	for _, f := range features {
		doc, err := parseGherkin(f)
		require.NoError(t, err)
		for _, p := range doc.Pickles() {
			missingSteps := []string{}
			for _, s := range p.Steps {
				err = steps.CheckExisting(s.Text)
				if err != nil {
					missingSteps = append(missingSteps, err.Error())
				}
			}
			if len(missingSteps) != 0 {
				require.NoError(t, errors.New(strings.Join(missingSteps, "\n")))
			}
		}
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
					tags := []string{}
					for _, t := range p.Tags {
						tags = append(tags, t.Name)
					}
					var err error
					defer func() {
						r := recover()
						if r != nil {
							ne, isErr := r.(error)
							if !isErr {
								ne = fmt.Errorf("PANIC: %v", r)
							}
							err = ne
						}

						steps.ExecuteAfterScenarios(w, doc.Feature.Name, p.Name, tags, err)
					}()
					err = steps.ExecuteBeforeScenarios(w, doc.Feature.Name, p.Name, tags)
					require.NoError(t, err)
					for _, s := range p.Steps {
						err = steps.ExecuteStep(s.Text, w)
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
