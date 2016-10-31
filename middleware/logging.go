package middleware

import (
	"bemicro/gateway"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Logging interceptor for grpc
func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	traceID := gateway.HarvestTraceID(ctx)
	log.Printf("tid=%s, moethd=%s, req=%s, took=%v, resp=%v, err=%v \n", traceID, info.FullMethod, marshal(req), time.Since(start), resp, err)

	return resp, err
}
