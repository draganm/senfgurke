package world

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
)

type World struct {
	Attributes map[string]interface{}
	T          *testing.T
	Require    *require.Assertions
	Assert     *assert.Assertions
	cleanups   [](func() error)
}

func New(t *testing.T) (*World, func() error) {
	w := &World{
		Attributes: map[string]interface{}{},
		T:          t,
		Require:    require.New(t),
		Assert:     assert.New(t),
	}
	return w, func() error {
		var finalErr error
		for _, cleanup := range w.cleanups {
			err := cleanup()
			if err != nil {
				finalErr = multierr.Append(finalErr, err)
			}
		}
		return finalErr
	}
}

func (w World) GetInt(name string) int {
	v, found := w.Attributes[name]
	if !found {
		panic(fmt.Errorf("could not find value of %q in the world", name))
	}

	iv, ok := v.(int)
	if !ok {
		panic(fmt.Errorf("value of %q (%#v) is not an int", name, v))
	}

	return iv
}

func (w World) GetString(name string) string {
	v, found := w.Attributes[name]
	if !found {
		panic(fmt.Errorf("could not find value of %q in the world", name))
	}

	iv, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("value of %q (%#v) is not a string", name, v))
	}

	return iv
}

func (w World) Put(name string, value interface{}) {
	w.Attributes[name] = value
}

func (w *World) AddCleanup(f func() error) {
	w.cleanups = append(w.cleanups, f)
}
