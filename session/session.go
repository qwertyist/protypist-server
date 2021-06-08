package session

import (
	"github.com/boltdb/bolt"
	"github.com/lithammer/shortuuid"
)

type Session struct {
	id       string
	password string
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
