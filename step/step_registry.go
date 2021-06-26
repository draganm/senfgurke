package step

import (
	"errors"
	"fmt"

	"github.com/draganm/senfgurke/testctx"
)

// type Step func()

type step struct {
	pattern string
	impl    func(w testctx.Context) error
}
type Registry struct {
	steps []step
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Given(pattern string, impl func(w testctx.Context) error) error {
	r.steps = append(r.steps, step{pattern: pattern, impl: impl})
	return nil
}

func (r *Registry) When(pattern string, impl func(w testctx.Context) error) error {
	r.steps = append(r.steps, step{pattern: pattern, impl: impl})
	return nil
}

func (r *Registry) Then(pattern string, impl func(w testctx.Context) error) error {
	r.steps = append(r.steps, step{pattern: pattern, impl: impl})
	return nil
}

func (r *Registry) Execute(text string, w testctx.World) error {

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

var errNotMatching = errors.New("not matching")

func (s step) execute(text string, w testctx.World) error {

	params, err := Match(s.pattern, text)
	if err != nil {
		return err
	}

	tc := testctx.New(params, w)
	return s.impl(tc)
}
