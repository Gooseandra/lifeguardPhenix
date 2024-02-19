package sessions

import (
	"fmt"
	"log"
	"swagger/models"
	"sync"
	"time"

	"github.com/google/uuid"

	"swagger/services"
)

type (
	SessionEntity struct {
		next, prev *SessionEntity
		time       time.Time
		user       models.User
		iD         uuid.UUID
	}

	SessionId = uuid.UUID

	Sessions struct {
		duration    time.Duration
		event       chan struct{}
		log         *services.Log
		mutex       sync.Mutex
		first, last *SessionEntity
		dict        map[uuid.UUID]*SessionEntity
		users       models.UserManager
	}
)

func NewSessions(l *services.Log, s models.UserManager, d time.Duration) Sessions {
	result := Sessions{dict: map[uuid.UUID]*SessionEntity{}, duration: d, log: l, users: s}
	go result.routine()
	return result
}

func (e SessionEntity) ID() SessionId { return e.iD }

func (e SessionEntity) User() models.User { return e.user }

func (s *Sessions) erase() {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	log.Printf("Sessions.erase 1\n")
	for s.first != nil && time.Now().After(s.first.time.Add(s.duration)) {
		log.Printf("Sessions.erase 2: id=%v time=%v\n", s.first.iD, s.first.time)
		delete(s.dict, s.first.iD)
		s.first = s.first.next
		if s.first == nil {
			s.last = nil
		}
	}
}

func (s Sessions) Get(i services.AuthID) (models.User, error) {
	defer s.erase()
	si, e := uuid.Parse(i)
	if e != nil {
		return nil, e
	}
	defer s.mutex.Unlock()
	s.mutex.Lock()
	if r, o := s.dict[si]; o {
		r.time = time.Now()
		return r.user, nil
	}
	return nil, fmt.Errorf("not found")
}

func (s *Sessions) New(n models.UserNickName, p models.UserPassword) (services.AuthID, error) {
	s.erase()
	u, e := s.users.ByName(n)
	if e != nil {
		return "", e
	}
	for {
		s.mutex.Lock()
		// TODO: нужен счетчик
		i := uuid.New()
		if _, o := s.dict[i]; !o {
			entity := &SessionEntity{next: nil, prev: s.last, iD: i, time: time.Now(), user: u}
			log.Printf("Sessions.New id=%v time=%v\n", entity.iD, entity.time)
			s.dict[i] = entity
			if s.last == nil {
				s.first, s.last = entity, entity
			} else {
				s.last, s.last.next = entity, entity
			}
			s.mutex.Unlock()
			return i.String(), nil
		}
		s.mutex.Unlock()
	}
}

func (s *Sessions) routine() {
	for {
		first := s.first
		log.Println("Sessions.routine 1")
		if first == nil {
			timer := time.After(s.duration)
			<-timer
		} else {
			timer := time.After(time.Now().Sub(first.time))
			<-timer
			s.erase()
		}
	}
}
