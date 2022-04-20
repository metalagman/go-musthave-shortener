package grpcservice

import (
	"google.golang.org/grpc"
	pb "shortener/api/proto"
	"shortener/internal/pkg/grpcserver"
)

type ShortenerService struct {
	pb.UnimplementedShortenerServer
}

func (s *ShortenerService) Shorten(ctx context.Context, request *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShortenerService) BatchShorten(ctx context.Context, request *pb.BatchShortenRequest) (*pb.BatchShortenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShortenerService) Expand(ctx context.Context, request *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShortenerService) BatchRemove(ctx context.Context, request *pb.BatchRemoveRequest) (*pb.BatchRemoveResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShortenerService) mustEmbedUnimplementedShortenerServer() {
	//TODO implement me
	panic("implement me")
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
