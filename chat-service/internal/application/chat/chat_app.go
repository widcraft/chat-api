package chat

import (
	log "github.com/sirupsen/logrus"
	"github.com/widcraft/chat-service/internal/domain/dto"
	"github.com/widcraft/chat-service/internal/port"
)

type workerPool interface {
	RegisterJob(func())
}

type ChatApp struct {
	logger     *log.Logger
	roomManger *roomManager
	repo       port.ChatRepository
}

func NewChatApp(logger *log.Logger, pool workerPool, repo port.ChatRepository) *ChatApp {
	return &ChatApp{
		logger:     logger,
		roomManger: NewRoomManager(),
		repo:       repo,
	}
}

func (app *ChatApp) Connect(roomIdx uint, client port.ChatClient) error {
	// TODO: check if user is valid from user service
	app.roomManger.add(roomIdx, client)
	return nil
}

func (app *ChatApp) Disconnect(roomIdx uint, client port.ChatClient) error {
	return app.roomManger.quit(roomIdx, client)
}

func (app *ChatApp) SendMessge(message dto.MessageDto) error {
	return app.SendMessge(message)
}

func (app *ChatApp) GetMessages(roomIdx uint) ([]dto.MessageDto, error) {
	return nil, nil
}