package handler

import "context"

type Producer interface {
	Send(
		ctx context.Context,
		key string,
		value any,
	) error
}
