package testctx

type Context struct {
	Params
	World
}

func New(params []interface{}) Context {
	return Context{
		Params: Params(params),
		World:  World{},
	}
}
