package world

import "fmt"

type World map[string]interface{}

func (w World) GetInt(name string) int {
	v, found := w[name]
	if !found {
		panic(fmt.Errorf("could not find value of %q in the world", name))
	}

	iv, ok := v.(int)
	if !ok {
		panic(fmt.Errorf("could not find value of %q (%#v) is not int", name, v))
	}

	return iv
}
