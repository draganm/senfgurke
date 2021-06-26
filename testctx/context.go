package testctx

type Context struct {
	Params
	World
}

func New(params []interface{}, w World) Context {
	return Context{
		Params: Params(params),
		World:  w,
	}
}
