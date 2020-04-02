package main

import "sync"

type Store struct {
	sync.Mutex
	subscriptions map[string]bool
	notified      bool
}

func (s *Store) Subscribe(id string) {
	s.Lock()
	defer s.Unlock()
	s.subscriptions[id] = true
}

func (s *Store) Unsubscribe(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.subscriptions, id)
}

func (s *Store) List() []string {
	s.Lock()
	defer s.Unlock()
	subs := make([]string, 0, len(s.subscriptions))
	for sub := range s.subscriptions {
		subs = append(subs, sub)
	}
	return subs
}

func (s *Store) SetNotified(notified bool) {
	s.Lock()
	defer s.Unlock()
	s.notified = notified
}

func (s *Store) IsNotified() bool {
	s.Lock()
	defer s.Unlock()
	return s.notified
}
