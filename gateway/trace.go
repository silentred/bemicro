package gateway

import (
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const TraceKey = "trace_id"

// GetTraceIDPair returns key value pair of trace_id
func GetTraceIDPair() []string {
	uuid := GetUUID()
	return []string{TraceKey, uuid}
}

// HarvestTraceID get trace_id from context
func HarvestTraceID(ctx context.Context) string {
	if md, ok := metadata.FromContext(ctx); ok {
		values := md[TraceKey]
		if len(values) > 0 {
			return values[0]
		}
	}

	return ""
}

// GetUUID returns uuid
func GetUUID() string {
	return uuid.NewV4().String()
}

func MergeStrings(ctx context.Context, pairs ...[]string) context.Context {
	var result []string
	for _, val := range pairs {
		result = append(result, val...)
	}
	md := metadata.Pairs(result...)

	return metadata.NewContext(ctx, md)
}
