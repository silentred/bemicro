package middleware

import (
	"errors"

	"bemicro/gateway"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	// ErrBadAuth means not autherized request
	ErrBadAuth = errors.New("access not granted")
)

// Auth interceptor for grpc
func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if gateway.IsValidAuth(ctx) {
		return handler(ctx, req)
	}
	return nil, ErrBadAuth
}
