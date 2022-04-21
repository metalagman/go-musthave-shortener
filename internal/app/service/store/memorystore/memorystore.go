package memorystore

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"shortener/internal/app/service/store"
	"sync"
	"time"
)

var ErrAlreadyStarted = errors.New("already started")
var ErrNotStarted = errors.New("not started")

type Store struct {
	mu              sync.RWMutex
	listenAddr      string
	baseURL         string
	counter         uint64
	base            int
	db              db
	dbFilePath      string
	dbFlushInterval time.Duration
	dbFlushCh       chan struct{}
	dbFlushTicker   *time.Ticker
}

func (s *Store) BatchWrite(uid string, in []store.Record) ([]store.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) BatchRemove(uid string, ids ...string) error {
	//TODO implement me
	panic("implement me")
}

type db map[uint64]dbRow
type dbRow struct {
	ID          string
	OriginalURL string
	ShortURL    string
	UID         string
}

func NewStore(opts ...StoreOption) *Store {
	const (
		defaultBase          = 36
		defaultFlushInterval = time.Second * 5
	)
	s := &Store{
		base:            defaultBase,
		dbFlushInterval: defaultFlushInterval,
		db:              make(db),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type StoreOption func(*Store)

func WithListenAddr(v string) StoreOption {
	return func(s *Store) {
		s.listenAddr = v
	}
}

func WithBaseURL(v string) StoreOption {
	return func(s *Store) {
		s.baseURL = v
	}
}

func WithFilePath(v string) StoreOption {
	return func(s *Store) {
		s.dbFilePath = v
	}
}

func WithFlushInterval(v time.Duration) StoreOption {
	return func(s *Store) {
		s.dbFlushInterval = v
	}
}

// Start periodic db flushing
func (s *Store) Start() error {
	start := make(chan struct{})
	defer close(start)

	// lock store to init ticker
	s.mu.Lock()
	defer s.mu.Unlock()

	// if already started
	if s.dbFlushTicker != nil {
		return ErrAlreadyStarted
	}

	if err := s.readDB(false); err != nil {
		return fmt.Errorf("serve error: %w", err)
	}

	// init new ticker and worker signal channel
	s.dbFlushCh = make(chan struct{})
	s.dbFlushTicker = time.NewTicker(s.dbFlushInterval)

	go func() {
		<-start
		for {
			select {
			case <-s.dbFlushCh:
				return
			case <-s.dbFlushTicker.C:
				//log.Print("timer writing db")
				s.mu.RLock()
				_ = s.writeDB(false)
				s.mu.RUnlock()
			}
		}
	}()

	return nil
}

// Stop periodic db flushing
func (s *Store) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// nothing to stop
	if s.dbFlushTicker == nil {
		return ErrNotStarted
	}

	log.Print("shutting down db")
	// stop goroutine worker
	close(s.dbFlushCh)

	// stop and reset ticker
	s.dbFlushTicker.Stop()

	// write db to file
	_ = s.writeDB(false)
	return nil
}

// readDB from file
func (s *Store) readDB(doLock bool) error {
	file, err := os.OpenFile(s.dbFilePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error reading db at %q: %w", s.dbFilePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := gob.NewDecoder(file)

	if doLock {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	err = decoder.Decode(&s.db)
	if err != nil && err != io.EOF {
		return fmt.Errorf("decode error: %w", err)
	}
	s.counter = uint64(len(s.db))
	log.Printf("db records loaded: %d", s.counter)

	return nil
}

// writeDB db to file
func (s *Store) writeDB(doLock bool) error {
	file, err := os.OpenFile(s.dbFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error writing db at %q: %w", s.dbFilePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := gob.NewEncoder(file)

	if doLock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	err = encoder.Encode(&s.db)
	if err != nil && err != io.EOF {
		return fmt.Errorf("encode error: %w", err)
	}

	return nil
}
