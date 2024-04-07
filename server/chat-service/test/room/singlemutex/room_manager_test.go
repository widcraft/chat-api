/*
 * Performance comparison between "Room Manager of Production" and "Room Manager using single Mutex"
 * Result: Production uses more memory, but it's 6x faster when dealing with Broadcasting
 */
package singlemutex

import (
	"chat-service/internal/application/dto"
	"chat-service/internal/domain/service"
	"chat-service/internal/port/driver"
	"fmt"
	"sync"
	"testing"
	"time"
)

type MockClient struct {
	roomId string
	userId string
}

func (m *MockClient) RoomId() string {
	return m.roomId
}

func (m *MockClient) UserId() string {
	return m.userId
}

func (m *MockClient) Send(message dto.Message) error {
	time.Sleep(time.Millisecond)
	return nil
}

var (
	benchRoomCount        = 100
	benchParticipantCount = 100
	benchBroadcastCount   = 10
)

func BenchmarkRoomManager(b *testing.B) {
	// given
	var clients []*MockClient
	roomManager := service.NewRoomManager()

	for i := range benchRoomCount {
		for j := range benchParticipantCount {
			clients = append(clients, &MockClient{
				roomId: fmt.Sprintf("room-bench-%d", i),
				userId: fmt.Sprintf("user-bench-%d", j),
			})
		}
	}

	// when
	var wg sync.WaitGroup

	for _, client := range clients {
		wg.Add(1)

		go func(c driver.MessengerClient) {
			defer wg.Done()

			roomManager.Join(c)
			defer roomManager.Leave(c)

			for range benchBroadcastCount {
				_ = roomManager.Broadcast(dto.Message{RoomId: c.RoomId()})
			}
		}(client)
	}

	wg.Wait()
}

func BenchmarkSingleMutexRoomManager(b *testing.B) {
	// given
	var clients []*MockClient
	roomManager := NewRoomManager()

	for i := range benchRoomCount {
		for j := range benchParticipantCount {
			clients = append(clients, &MockClient{
				roomId: fmt.Sprintf("room-bench-%d", i),
				userId: fmt.Sprintf("user-bench-%d", j),
			})
		}
	}

	// when
	var wg sync.WaitGroup

	for _, client := range clients {
		wg.Add(1)

		go func(c driver.MessengerClient) {
			defer wg.Done()

			roomManager.Join(c)
			defer roomManager.Leave(c)

			for range benchBroadcastCount {
				_ = roomManager.Broadcast(dto.Message{RoomId: c.RoomId()})
			}
		}(client)
	}

	wg.Wait()
}