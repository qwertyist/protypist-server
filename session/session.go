package session

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
)

type Session struct {
	id       string
	password string
	conns    []*websocket.Conn
	buf      []byte
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
	sess.id = id
	sess.password = password
	s.ss[id] = sess
	return id
}

func (s *Service) GetSession(id string) []*websocket.Conn {
	return s.ss[id].conns
}

func (s *Service) JoinSession(id string, conn *websocket.Conn) {
	s.ss[id].conns = append(s.ss[id].conns, conn)
}

func (s *Service) WriteBuf(id string, buf []byte) error {
	s.ss[id].buf = buf
	return nil
}

func (s *Service) ReadBuf(id string) []byte {
	return s.ss[id].buf
}
