package calculator_test

import (
	"fmt"
	"testing"

	"github.com/draganm/senfgurke/step"
	"github.com/draganm/senfgurke/testctx"
	"github.com/draganm/senfgurke/testrunner"
)

func TestFeatures(t *testing.T) {
	testrunner.RunScenarios(t, steps)
}

var steps = step.NewRegistry()

var _ = steps.When("I add 2 and 3", func(w testctx.Context) error {
	w.World.Put("result", 2+3)
	return nil
})

var _ = steps.Then("the result should be 5", func(w testctx.Context) error {
	result := w.World.GetInt("result")
	if result != 5 {
		return fmt.Errorf("value %d is not equal to 5", result)
	}
	return nil
})
