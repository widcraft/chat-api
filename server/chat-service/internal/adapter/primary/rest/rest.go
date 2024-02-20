package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/widcraft/chat-service/internal/adapter/primary/rest/handler/chat"
	"github.com/widcraft/chat-service/internal/port/primary"
	"github.com/widcraft/chat-service/internal/port/secondary"
)

type Rest struct {
	logger  secondary.Logger
	server  *http.Server
	chatApp primary.MessageService
}

func New(logger secondary.Logger, chatApp primary.MessageService) *Rest {
	router := gin.Default()
	group := router.Group("/api/v1")

	chat.New(logger, chatApp).Register(group)

	return &Rest{
		logger:  logger,
		chatApp: chatApp,
		server: &http.Server{
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (rest *Rest) Run(port string) {
	rest.server.Addr = ":" + port
	err := rest.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		rest.logger.Errorf("websocket server error: %s", err)
	}
}

func (rest *Rest) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return rest.server.Shutdown(ctx)
}
