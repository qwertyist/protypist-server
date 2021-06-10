package session

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
)

type Session struct {
	ID       string
	Password string
	Conns    []*websocket.Conn
	Buf      []byte
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

func (s *Service) GetSessionClients(id string) []*websocket.Conn {
	return s.ss[id].Conns
}

func (s *Service) JoinSession(id string, conn *websocket.Conn) {
	s.ss[id].Conns = append(s.ss[id].Conns, conn)
}

func (s *Service) LeaveSession(id string, conn *websocket.Conn) {
	a := s.ss[id].Conns
	for i, c := range s.ss[id].Conns {
		if c == conn {
			copy(a[i:], s.ss[id].Conns[i+1:])
			a[len(a)-1] = nil // or the zero value of T
			a = a[:len(a)-1]
		}
	}
	s.ss[id].Conns = a
}

func (s *Service) WriteBuf(id string, buf []byte) error {
	s.ss[id].Buf = buf
	return nil
}

func (s *Service) ReadBuf(id string) []byte {
	return s.ss[id].Buf
}
