package goaop

type Pointcut interface {
	OnEntry()
	OnReturn(err error)
}

type PointcutBuilder interface {
	Build(name string) Pointcut
}
