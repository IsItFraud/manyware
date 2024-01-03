package manyware

type (
	Executor[T any]   func(ctx T) error
	Middleware[T any] func(next Executor[T]) Executor[T]
)

func Prepare[T any](executor Executor[T], middleware ...Middleware[T]) Executor[T] {
	for i := len(middleware) - 1; i >= 0; i-- {
		executor = middleware[i](executor)
	}
	return executor
}
