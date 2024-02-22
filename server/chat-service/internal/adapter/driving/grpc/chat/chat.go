package chat

import (
	"errors"

	"github.com/widcraft/chat-service/internal/adapter/driving/grpc/chat/pb"
	"github.com/widcraft/chat-service/internal/application/dto"
	"github.com/widcraft/chat-service/internal/port/driven"
	"github.com/widcraft/chat-service/internal/port/driving"
)

type Server struct {
	pb.UnimplementedChatServer
	logger driven.Logger
	app    driving.MessageService
}

func New(logger driven.Logger, app driving.MessageService) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Connect(stream pb.Chat_ConnectServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	joinReq, ok := req.GetType().(*pb.ChatReq_Join)
	if !ok {
		return errors.New("should join before other request")
	}

	return s.handleConnetcion(stream, joinReq.Join)
}

func (s *Server) handleConnetcion(stream pb.Chat_ConnectServer, joinReq *pb.JoinReq) error {
	client := &client{
		stream:  stream,
		roomIdx: uint(joinReq.UserIdx),
		userIdx: uint(joinReq.RoomIdx),
	}

	s.app.Join(client)
	defer s.app.Leave(client)

	return s.handleMessage(client)
}

func (s *Server) handleMessage(client *client) error {
	for {
		req, err := client.stream.Recv()
		if err != nil {
			return err
		}

		switch typedReq := req.GetType().(type) {
		case *pb.ChatReq_Message:
			s.sendMessage(client, typedReq.Message)
		default:
			return errors.New("wrong request type")
		}
	}
}

func (s *Server) sendMessage(client *client, payload *pb.MessageReq) {
	err := s.app.SendMessage(&dto.MessageDto{
		RoomIdx:  client.roomIdx,
		UserIdx:  client.userIdx,
		Name:     client.name,
		ImageUrl: client.imageUrl,
		Message:  payload.GetMessage(),
	})
	if err != nil {
		s.logger.Errorf("send message failed: %s", err)
	}
}