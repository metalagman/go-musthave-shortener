package grpcservice

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "shortener/api/proto"
	"shortener/internal/app/service/store"
	"shortener/internal/pkg/user"
)

type ShortenerService struct {
	pb.UnimplementedShortenerServer

	store store.Store
}

func NewShortenerService(st store.Store) *ShortenerService {
	s := &ShortenerService{
		store: st,
	}

	return s
}

func (s *ShortenerService) Init() ServiceInit {
	return func(registrar grpc.ServiceRegistrar) {
		pb.RegisterShortenerServer(registrar, s)
	}
}

func (s *ShortenerService) Shorten(ctx context.Context, request *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	resp := &pb.ShortenResponse{}

	uid := user.ReadUID(ctx)
	shortURL, err := s.store.WriteURL(request.GetOriginalUrl(), uid)
	if err != nil {
		var errConflict *store.ConflictError
		if errors.As(err, &errConflict) {
			resp.ShortUrl = errConflict.ExistingURL
			return resp, status.Error(codes.Internal, err.Error())
		}
		return nil, err
	}

	resp.ShortUrl = shortURL
	return resp, nil
}

func (s *ShortenerService) BatchShorten(ctx context.Context, request *pb.BatchShortenRequest) (*pb.BatchShortenResponse, error) {
	uid := user.ReadUID(ctx)

	storeReq := make([]store.Record, len(request.Items))
	for i, rec := range request.Items {
		storeReq[i] = store.Record{
			CorrelationID: rec.GetCorrelationId(),
			OriginalURL:   rec.GetOriginalUrl(),
		}
	}

	storeRes, err := s.store.BatchWrite(uid, storeReq)
	if err != nil {
		if errors.Is(err, store.ErrBadInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	resp := &pb.BatchShortenResponse{
		Items: make([]*pb.BatchShortenResponseItem, len(request.Items)),
	}
	for i, rec := range storeRes {
		resp.Items[i] = &pb.BatchShortenResponseItem{
			CorrelationId: rec.CorrelationID,
			ShortUrl:      rec.ShortURL,
		}
	}

	return resp, nil
}

func (s *ShortenerService) Expand(ctx context.Context, request *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	u, err := s.store.ReadURL(request.GetId())
	if err != nil {
		if errors.Is(err, store.ErrDeleted) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.ExpandResponse{
		OriginalUrl: u,
	}, nil
}

func (s *ShortenerService) BatchRemove(ctx context.Context, request *pb.BatchRemoveRequest) (*pb.BatchRemoveResponse, error) {
	uid := user.ReadUID(ctx)

	if err := s.store.BatchRemove(uid, request.GetIds()...); err != nil {
		if errors.Is(err, store.ErrBadInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BatchRemoveResponse{}, nil
}

func (s *ShortenerService) UserData(ctx context.Context, request *pb.UserDataRequest) (*pb.UserDataResponse, error) {
	uid := user.ReadUID(ctx)
	rows := s.store.ReadUserData(uid)

	resp := &pb.UserDataResponse{
		Items: make([]*pb.UserDataResponseItem, len(rows)),
	}

	for i, row := range rows {
		resp.Items[i] = &pb.UserDataResponseItem{
			ShortUrl:    row.ShortURL,
			OriginalUrl: row.OriginalURL,
		}
	}

	return resp, nil
}
