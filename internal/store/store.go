package store

import (
	"sync"

	"github.com/svartvalp/fqw-consensus/internal/log"
)

type Store interface {
	Apply(entry log.Entry)
	Get(key string) string
	GetAll() map[string]string
	Clear()
}

type store struct {
	mut   *sync.RWMutex
	state map[string]string
}

func (s *store) Clear() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.state = make(map[string]string)
}

func (s *store) Apply(entry log.Entry) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.state[entry.Data.Key] = entry.Data.Value
}

func (s *store) Get(key string) string {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.state[key]
}

func (s *store) GetAll() map[string]string {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.state
}

func New() Store {
	return &store{
		mut:   &sync.RWMutex{},
		state: make(map[string]string),
	}
}
