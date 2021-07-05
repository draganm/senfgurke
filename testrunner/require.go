package testrunner

import (
	"github.com/draganm/senfgurke/world"
	"github.com/stretchr/testify/require"
)

const requireWorldKey = "require"

func Require(w world.World) *require.Assertions {
	return w[requireWorldKey].(*require.Assertions)
}
