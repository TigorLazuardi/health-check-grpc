package app

import (
	"context"
	"log"
	"sync"

	"github.com/tigorlazuardi/health-check-grpc/app/hcproto"
)

type Server struct {
	hcproto.UnimplementedHealthCheckServer
	subscribers *sync.Map
}

type sub struct {
	stream hcproto.HealthCheck_SubscribeServer
	finish chan<- struct{}
}

func New() hcproto.HealthCheckServer {
	return &Server{
		subscribers: &sync.Map{},
	}
}

func (s *Server) Subscribe(payload *hcproto.SubPayload, stream hcproto.HealthCheck_SubscribeServer) error {
	log.Println("a subscriber joined with id: ", payload.Id)

	fin := make(chan struct{})

	s.subscribers.Store(payload.Id, sub{stream: stream, finish: fin})

	ctx := stream.Context()

	select {
	case <-fin:
		log.Println("closing stream for subscriber ", payload.Id)
	case <-ctx.Done():
		log.Println("client has disconnected. id: ", payload.Id)
	}
	s.subscribers.Delete(payload.Id)
	return nil
}

func (s *Server) Unsubscribe(ctx context.Context, payload *hcproto.SubPayload) (*hcproto.Ack, error) {
	v, exist := s.subscribers.Load(payload.Id)
	if !exist {
		return &hcproto.Ack{Message: "not subscribed"}, nil
	}
	v.(sub).finish <- struct{}{}
	return &hcproto.Ack{Message: "unsubscribe client: " + payload.Id}, nil
}

func (s *Server) Subscribers() *sync.Map {
	return s.subscribers
}
