package server

import (
	"sync"

	"github.com/teryble09/go_grpc_chat/proto"
)

type ConnStorage struct {
	conn map[string]proto.Chat_StreamServer
	sync.Mutex
}

func NewConnStorage() ConnStorage {
	return ConnStorage{conn: make(map[string]proto.Chat_StreamServer, 100)}
}

func (c *ConnStorage) SendMessageToActiveUsers(mes *proto.Message) {
	c.Lock()
	wg := sync.WaitGroup{}
	for username, cnn := range c.conn {
		if username == mes.GetUsername() {
			continue
		}
		wg.Add(1)
		go func() {
			err := cnn.Send(mes)
			if err != nil {
				delete(c.conn, username)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	c.Unlock()
}

func (c *ConnStorage) RegisterNewUser(username string, cnn proto.Chat_StreamServer) {
	c.Lock()
	c.conn[username] = cnn
	c.Unlock()
}
