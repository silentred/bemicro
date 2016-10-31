package middleware

import (
	"bemicro/gateway"
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	// ErrBadUser means invalid user claim
	ErrBadUser = errors.New("invalid user token")
	// ErrKeyNotExist means no key in metadata
	ErrKeyNotExist = errors.New("user key not exists in metadata")
)

// UserVerify interceptor for grpc
func UserVerify(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	_, err = verifyUserClaim(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func verifyUserClaim(ctx context.Context) (*gateway.UserClaims, error) {
	if md, ok := metadata.FromContext(ctx); ok {
		values := md[gateway.UserClaimKey]
		if len(values) > 0 {
			user, err := gateway.VerifyUserToken(values[0])
			if err != nil {
				fmt.Println(err)
				return nil, ErrBadUser
			}
			return user, nil
		}
	}

	return nil, ErrKeyNotExist
}
