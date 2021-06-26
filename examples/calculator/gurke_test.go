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

var _ = steps.When("I add {int} and {int}", func(w testctx.Context) error {
	w.World.Put("result", w.Params.GetInt(0)+w.Params.GetInt(1))
	return nil
})

var _ = steps.Then("the result should be {int}", func(w testctx.Context) error {
	result := w.World.GetInt("result")
	expected := w.Params.GetInt(0)
	if result != expected {
		return fmt.Errorf("value %d is not equal to %d", result, expected)
	}
	return nil
})

var _ = steps.When("I convert number {int} to string", func(w testctx.Context) error {
	w.World.Put("result", fmt.Sprintf("%d", w.Params.GetInt(0)))
	return nil
})

var _ = steps.Then("the string should be {string}", func(w testctx.Context) error {
	result := w.World.GetString("result")
	expected := w.Params.GetString(0)
	if result != expected {
		return fmt.Errorf("value %q is not equal to %q", result, expected)
	}
	return nil
})
