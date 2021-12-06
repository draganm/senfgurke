package testrunner

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	// "github.com/cucumber/gherkin-go"
	gherkin "github.com/cucumber/common/gherkin/go/v22"
	messages "github.com/cucumber/common/messages/go/v17"
	"github.com/gofrs/uuid"

	// "github.com/cucumber/gherkin-go"
	"github.com/draganm/senfgurke/step"
	"github.com/draganm/senfgurke/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunScenarios(t *testing.T, steps *step.Registry) {
	entries, err := os.ReadDir(".")
	require.NoError(t, err)

	featureFiles := []string{}
	for _, e := range entries {

		if !strings.HasSuffix(e.Name(), ".feature") {
			continue
		}

		i, err := e.Info()
		require.NoError(t, err)
		if !i.Mode().IsRegular() {
			continue
		}

		featureFiles = append(featureFiles, e.Name())
	}

	// check for missing steps
	for _, f := range featureFiles {
		doc, err := parseGherkin(f)
		require.NoError(t, err)
		missingSteps := []string{}
		for _, fc := range doc.Feature.Children {

			switch {
			case fc.Scenario != nil:
				for _, s := range fc.Scenario.Steps {
					err = steps.CheckExisting(s.Text)
					if err != nil {
						missingSteps = append(missingSteps, err.Error())
					}
				}
			case fc.Background != nil:
				for _, s := range fc.Background.Steps {
					err = steps.CheckExisting(s.Text)
					if err != nil {
						missingSteps = append(missingSteps, err.Error())
					}
				}
			}

			if len(missingSteps) != 0 {
				require.NoError(t, errors.New(strings.Join(missingSteps, "\n")))
			}
		}

	}

	runWIP := false

	// check for WIP tag
outer:
	for _, f := range featureFiles {
		doc, err := parseGherkin(f)
		require.NoError(t, err)

		for _, ft := range doc.Feature.Tags {
			if ft.Name == "@WIP" {
				runWIP = true
				break outer
			}
		}

		for _, fc := range doc.Feature.Children {
			if fc.Scenario == nil {
				continue
			}

			for _, st := range fc.Scenario.Tags {
				if st.Name == "@WIP" {
					runWIP = true
					break outer
				}
			}
		}
	}

	for _, f := range featureFiles {
		f := f
		doc, err := parseGherkin(f)
		require.NoError(t, err)

		gotWIP := false

		for _, ft := range doc.Feature.Tags {
			if ft.Name == "@WIP" {
				gotWIP = true
			}
		}

		t.Run(fmt.Sprintf("%s(%s)", f, doc.Feature.Name), func(t *testing.T) {
			if !runWIP {
				t.Parallel()
			}

			backgroundSteps := []*messages.Step{}
			for _, fc := range doc.Feature.Children {
				if fc.Background != nil {
					backgroundSteps = append(backgroundSteps, fc.Background.Steps...)
				}
			}

			for _, fc := range doc.Feature.Children {

				if fc.Scenario == nil {
					continue
				}

				fc := fc

				for _, pt := range fc.Scenario.Tags {
					if pt.Name == "@WIP" {
						gotWIP = true
					}
				}

				t.Run(fmt.Sprintf("Scenario: %s", fc.Scenario.Name), func(t *testing.T) {
					if !runWIP {
						t.Parallel()
					}
					if runWIP && !gotWIP {
						t.Skip("not marked as @WIP")
					}
					w, cleanup := world.New(t)
					tags := []string{}
					for _, t := range fc.Scenario.Tags {
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

						assert.NoError(t, err)

						err = cleanup()
						assert.NoError(t, err)

						err = steps.ExecuteAfterScenarios(w, doc.Feature.Name, fc.Scenario.Name, tags, err)
						assert.NoError(t, err)
					}()
					err = steps.ExecuteBeforeScenarios(w, doc.Feature.Name, fc.Scenario.Name, tags)
					require.NoError(t, err)
					allSteps := append([]*messages.Step{}, backgroundSteps...)

					allSteps = append(allSteps, fc.Scenario.Steps...)
					for _, s := range allSteps {
						err = steps.ExecuteStep(s.Text, w)
						require.NoError(t, err)
					}
				})
			}

		})
	}

	if runWIP {
		t.Error("was running features/scenarios wit @WIP tags")
	}

}

func parseGherkin(name string) (gherkinDocument *messages.GherkinDocument, err error) {
	// gherkin.GherkinDocument
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return gherkin.ParseGherkinDocument(f, func() string { return uuid.Must(uuid.NewV4()).String() })
}
