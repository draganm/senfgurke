package calculator_test

import (
	"fmt"
	"testing"

	"github.com/draganm/senfgurke/step"
	"github.com/draganm/senfgurke/testrunner"
	"github.com/draganm/senfgurke/world"
)

func TestFeatures(t *testing.T) {
	testrunner.RunScenarios(t, steps)
}

var steps = step.NewRegistry()

var _ = steps.When("I add {int} and {int}", func(w world.World, a, b int) error {
	w.Put("result", a+b)
	return nil
})

var _ = steps.Then("the result should be {int}", func(w world.World, expected int) error {
	result := w.GetInt("result")
	if result != expected {
		return fmt.Errorf("value %d is not equal to %d", result, expected)
	}
	return nil
})

var _ = steps.When("I convert number {int} to string", func(w world.World, a int) error {
	w.Put("result", fmt.Sprintf("%d", a))
	return nil
})

var _ = steps.Then("the string should be {string}", func(w world.World, expected string) error {
	result := w.GetString("result")
	if result != expected {
		return fmt.Errorf("value %q is not equal to %q", result, expected)
	}
	return nil
})

var _ = steps.BeforeScenario(func(w world.World, featureName, scenarioName string, tags []string) error {
	fmt.Println("running", featureName, "->", scenarioName)
	return nil
})

var _ = steps.AfterScenario(func(w world.World, featureName, scenarioName string, tags []string, err error) error {
	fmt.Println("done", featureName, "->", scenarioName, err)
	return nil
})
