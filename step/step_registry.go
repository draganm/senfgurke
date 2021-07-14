package step

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/draganm/senfgurke/world"
)

type step struct {
	matcher *stepMatcher
	impl    interface{}
}

type Registry struct {
	beforeScenarios []func(w *world.World, featureName, scenarioName string, tags []string) error
	afterScenarios  []func(w *world.World, featureName, scenarioName string, tags []string, err error) error
	steps           []step
}

func NewRegistry() *Registry {
	return &Registry{}
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func (r *Registry) BeforeScenario(fn func(w *world.World, featureName, scenarioName string, tags []string) error) error {
	r.beforeScenarios = append(r.beforeScenarios, fn)
	return nil
}

func (r *Registry) AfterScenario(fn func(w *world.World, featureName, scenarioName string, tags []string, err error) error) error {
	r.afterScenarios = append(r.afterScenarios, fn)
	return nil
}

func (r *Registry) ExecuteBeforeScenarios(w *world.World, featureName, scenarioName string, tags []string) error {
	for _, bs := range r.beforeScenarios {
		err := bs(w, featureName, scenarioName, tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) ExecuteAfterScenarios(w *world.World, featureName, scenarioName string, tags []string, err error) error {
	for _, bs := range r.afterScenarios {
		err := bs(w, featureName, scenarioName, tags, err)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *Registry) addStep(pattern string, impl interface{}) error {
	matcher, err := newStepMatcher(pattern)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(impl)
	if t.Kind() != reflect.Func {
		return errors.New("step implementation is not a function")
	}

	if t.NumIn() != len(matcher.types)+1 {
		return fmt.Errorf(
			"expected step implementation to have %d parameters, but it has %d",
			len(matcher.types)+1,
			t.NumIn(),
		)
	}

	if !t.In(0).AssignableTo(reflect.TypeOf(&world.World{})) {
		return errors.New("first parameter of step implementation must be world.World")
	}

	for i, ty := range matcher.types {
		argKind := t.In(i + 1).Kind()
		switch ty {
		case "int":
			if argKind != reflect.Int {
				return fmt.Errorf("argument %d is expected to be %s, but was %s", i+1, ty, argKind.String())
			}
		case "string":
			if argKind != reflect.String {
				return fmt.Errorf("argument %d is expected to be %s, but was %s", i+1, ty, argKind.String())
			}
		default:
			return fmt.Errorf("unsupported step implementation argument type %q", ty)
		}
	}

	if t.NumOut() > 1 {
		return fmt.Errorf("expected step implementation to return at most one value, but it returns %d", t.NumOut())
	}

	if t.NumOut() == 1 && !t.Out(0).AssignableTo(errorInterface) {
		return errors.New("step implementation must return error type")
	}

	r.steps = append(r.steps, step{matcher: matcher, impl: impl})
	return nil

}

func (r *Registry) Given(pattern string, impl interface{}) error {
	return r.addStep(pattern, impl)
}

func (r *Registry) When(pattern string, impl interface{}) error {
	return r.addStep(pattern, impl)
}

func (r *Registry) Then(pattern string, impl interface{}) error {
	return r.addStep(pattern, impl)
}

func (r *Registry) ExecuteStep(text string, w *world.World) error {

	for _, s := range r.steps {
		err := s.execute(text, w)
		if err == errNotMatching {
			continue
		}

		if err != nil {
			return fmt.Errorf("while executing step %q: %s", text, err.Error())
		}

		return nil
	}

	return fmt.Errorf("no step found matching %q", text)
}

// var stringOrIntMatcher = regexp.MustCompile(`((?:-?\d+)|"((?:[^"]|(?:\\"))*))"`)

var stringOrIntMatcher = regexp.MustCompile(`(-?\d+)|("([^"]|\\")*")`)

func (r *Registry) CheckExisting(text string) error {

	for _, s := range r.steps {
		_, err := s.matcher.match(text)
		if err == errNotMatching {
			continue
		}

		if err == nil {
			return nil
		}

		return err
	}

	argIndexes := stringOrIntMatcher.FindAllStringIndex(text, -1)

	pattern := new(strings.Builder)
	args := new(strings.Builder)
	args.WriteString("w *world.World")

	prev := 0

	for i, ai := range argIndexes {

		pattern.WriteString(text[prev:ai[0]])

		if ai[0] > 0 {
			prevChar := text[ai[0]-1]
			if unicode.IsLetter(rune(prevChar)) || unicode.IsDigit(rune(prevChar)) {
				pattern.WriteString(text[ai[0]:ai[1]])
				prev = ai[1]
				continue
			}
		}

		if len(text) > ai[1] {
			nextChar := text[ai[1]]
			if unicode.IsLetter(rune(nextChar)) || unicode.IsDigit(rune(nextChar)) {
				pattern.WriteString(text[ai[0]:ai[1]])
				prev = ai[1]
				continue
			}
		}

		args.WriteString(", ")
		if text[ai[0]] == '"' {
			pattern.WriteString("{string}")
			args.WriteString(fmt.Sprintf("arg%d string", i))
		} else {
			pattern.WriteString("{int}")
			args.WriteString(fmt.Sprintf("arg%d int", i))
		}
		prev = ai[1]
	}

	pattern.WriteString(text[prev:])

	return errors.New(fmt.Sprintf(
		`// could not find step matching %q.
// to implement it, add following code:
var _ = steps.Then(%q, func(%s) {	
	w.Require.Fail("not yet implemented")
})
`,
		text,
		pattern.String(),
		args.String(),
	))
}

var errNotMatching = errors.New("not matching")

func (s step) execute(text string, w *world.World) error {

	params, err := s.matcher.match(text)
	if err != nil {
		return err
	}

	iv := reflect.ValueOf(s.impl)
	values := make([]reflect.Value, len(params)+1)
	values[0] = reflect.ValueOf(w)
	for i, p := range params {
		values[i+1] = reflect.ValueOf(p)
	}
	res := iv.Call(values)

	if len(res) == 0 {
		return nil
	}

	ri := res[0].Interface()
	if ri == nil {
		return nil
	}
	err = (ri).(error)

	return err
}
