package grpcservice

import (
	"google.golang.org/grpc"
	pb "shortener/api/proto"
	"shortener/internal/pkg/grpcserver"
)

type ShortenerService struct {
	pb.UnimplementedShortenerServer
}

func NewShortenerService() *ShortenerService {
	s := &ShortenerService{}

	return s
}

func (s *ShortenerService) Init() grpcserver.ServiceInit {
	return func(registrar grpc.ServiceRegistrar) {
		pb.RegisterShortenerServer(registrar, s)
	}
}
