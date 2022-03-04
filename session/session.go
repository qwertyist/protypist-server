package session

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
)

type Client struct {
	Conn *websocket.Conn
	Name string
}

type Session struct {
	ID         string
	Password   string
	MaxClients int
	Clients    []*Client
	Buf        []byte
}

type Service struct {
	db *bolt.DB
	ss map[string]*Session
}

func NewService(db *bolt.DB) *Service {
	s := &Service{}
	s.db = db
	s.ss = make(map[string]*Session)
	return s
}

func (s *Service) CreateSession(password string) string {
	sess := &Session{}
	id := shortuuid.New()
	sess.ID = id
	sess.Password = password
	s.ss[id] = sess
	return id
}

func (s *Service) GetSession(id string) *Session {
	return s.ss[id]
}

func (s *Service) GetSessionClients(id string) []*Client {
	return s.ss[id].Clients
}

func (s *Service) JoinSession(id string, client *Client) {
	s.ss[id].Clients = append(s.ss[id].Clients, client)
}

func (s *Service) LeaveSession(id string, client *Client) {
	a := s.ss[id].Clients
	for i, c := range s.ss[id].Clients {
		if c == client {
			copy(a[i:], a[i+1:])
			a[len(a)-1] = nil
			a = a[:len(a)-1]
		}
	}
	s.ss[id].Clients = a
}

func (s *Service) WriteBuf(id string, buf []byte) error {
	s.ss[id].Buf = buf
	return nil
}

func (s *Service) ReadBuf(id string) []byte {
	return s.ss[id].Buf
}
