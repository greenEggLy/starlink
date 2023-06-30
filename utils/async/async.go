package async

import (
	"context"
	"starlink/pb"
)

// Future interface has the method signature for await

type Future[T bool | []*pb.PositionInfo] interface {
	Await() T
}

type future[T bool | []*pb.PositionInfo] struct {
	await func(ctx context.Context) T
}

func (f future[T]) Await() T {

	return f.await(context.Background())

}

// Exec executes the async function

func Exec[T bool | []*pb.PositionInfo](f func() T) Future[T] {

	var result T

	c := make(chan struct{})

	go func() {

		defer close(c)

		result = f()

	}()

	return future[T]{
		await: func(ctx context.Context) T {

			select {

			case <-ctx.Done():
				return result

			case <-c:
				return result

			}

		},
	}

}
