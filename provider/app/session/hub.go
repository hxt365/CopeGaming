package session

import "sync"

type Hub struct {
	sessions map[string]*Session
	rwMutex  sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[string]*Session),
		rwMutex:  sync.RWMutex{},
	}
}

func (h *Hub) AddSession(s *Session) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	h.sessions[s.playerID] = s
}

func (h *Hub) RemoveSession(playerID string) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	if _, ok := h.sessions[playerID]; ok {
		delete(h.sessions, playerID)
	}
}

func (h *Hub) GetSession(playerID string) *Session {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	if s, ok := h.sessions[playerID]; ok {
		return s
	}

	return nil
}
