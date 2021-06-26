package testctx

import (
	"errors"
	"fmt"
)

type Params []interface{}

func (p Params) GetInt(idx int) int {
	if idx < 0 {
		panic(errors.New("index can't be negative"))
	}
	if len(p) >= idx {
		panic(fmt.Errorf("step does not have param %d", idx))
	}

	pv, ok := p[idx].(int)
	if !ok {
		panic(fmt.Errorf("param %d is not an int", idx))
	}

	return pv
}
