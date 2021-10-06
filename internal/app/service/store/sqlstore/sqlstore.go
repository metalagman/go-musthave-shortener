package sqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"runtime"
	"shortener/internal/app/logger"
	"sync"
)

type Store struct {
	baseURL   string
	base      int
	workerNum int
	db        *sql.DB
	log       logger.Logger

	jobs chan Job
}

var ErrAlreadyStarted = errors.New("already started")
var ErrNotStarted = errors.New("not started")

// New constructor
func New(db *sql.DB, opts ...Option) (*Store, error) {
	const (
		defaultBase = 36
	)

	s := &Store{
		base:      defaultBase,
		workerNum: runtime.GOMAXPROCS(0),
		db:        db,
		log:       logger.Global().Component("Store"),
	}

	if err := s.createTables(); err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

type Option func(*Store)

func WithBaseURL(url string) Option {
	return func(s *Store) {
		s.baseURL = url
	}
}

// Start db connection
func (s *Store) Start() error {
	if s.jobs != nil {
		return ErrAlreadyStarted
	}
	if err := s.createTables(); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	s.jobs = make(chan Job)
	wg := &sync.WaitGroup{}
	wg.Add(s.workerNum)

	for i := 0; i < s.workerNum; i++ {
		go func(id int, in <-chan Job) {
			defer wg.Done()
			s.log.Debug().Msgf("Worker [%d] started", id)
			for job := range in {
				s.log.Debug().Msgf("Worker [%d] doing job", id)
				job()
			}
			s.log.Debug().Msgf("Worker [%d] finished", id)
		}(i, s.jobs)
	}

	go func() {
		wg.Wait()
		s.log.Debug().Msgf("All workers finished")
	}()

	return nil
}

// Stop store db connection
func (s *Store) Stop() error {
	if s.jobs == nil {
		return ErrNotStarted
	}
	close(s.jobs)
	s.jobs = nil
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("db close: %w", err)
	}
	return nil
}
