package grpcservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"shortener/internal/pkg/user"
)

const mdKeyUID = "uid"

func UID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		uid := metautils.ExtractIncoming(ctx).Get(mdKeyUID)
		if uid == "" {
			uid = uuid.New().String()
		}
		ctx = user.WriteUID(ctx, uid)
		resp, err := handler(ctx, req)
		return resp, err
	}
}
