package store

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var _ Store = (*MemoryStore)(nil)

var ErrAlreadyServing = errors.New("already serving")
var ErrNotServing = errors.New("not serving")

type MemoryStore struct {
	mu              sync.RWMutex
	listenAddr      string
	baseURL         string
	counter         uint64
	base            int
	db              MemoryDB
	dbFilePath      string
	dbFlushInterval time.Duration
	dbFlushSignal   chan struct{}
	dbFlushTicker   *time.Ticker
}

type MemoryDB map[uint64]string

func NewMemoryStore(listenAddr string, baseURL string, dbFilePath string, dbFlushInterval time.Duration) *MemoryStore {
	return &MemoryStore{
		counter:         30,
		listenAddr:      listenAddr,
		baseURL:         baseURL,
		dbFilePath:      dbFilePath,
		base:            36,
		db:              make(MemoryDB),
		dbFlushInterval: dbFlushInterval,
	}
}

// Serve periodic db flushing
func (store *MemoryStore) Serve() error {
	start := make(chan struct{})
	defer close(start)

	// lock store to init ticker
	store.mu.Lock()
	defer store.mu.Unlock()

	// if already started
	if store.dbFlushTicker != nil {
		return ErrAlreadyServing
	}

	if err := store.readDB(false); err != nil {
		return fmt.Errorf("serve error: %w", err)
	}

	// init new ticker and worker signal channel
	store.dbFlushSignal = make(chan struct{})
	store.dbFlushTicker = time.NewTicker(store.dbFlushInterval)

	go func() {
		<-start
		for {
			select {
			case <-store.dbFlushSignal:
				return
			case <-store.dbFlushTicker.C:
				log.Print("timer writing db")
				store.mu.RLock()
				_ = store.writeDB(false)
				store.mu.RUnlock()
			}
		}
	}()

	return nil
}

// Shutdown periodic db flushing
func (store *MemoryStore) Shutdown() error {
	store.mu.Lock()
	defer store.mu.Unlock()

	// nothing to stop
	if store.dbFlushTicker == nil {
		return ErrNotServing
	}

	log.Print("shutting down db")
	// stop goroutine worker
	close(store.dbFlushSignal)

	// stop and reset ticker
	store.dbFlushTicker.Stop()
	store.dbFlushTicker = nil

	// write db to file
	_ = store.writeDB(false)
	return nil
}

func (store *MemoryStore) WriteURL(url string) (string, error) {
	if err := validateURL(url); err != nil {
		return "", err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	store.counter++
	store.db[store.counter] = url
	id := strconv.FormatUint(store.counter, store.base)

	return fmt.Sprintf("%s/%s", store.baseURL, id), nil
}

func (store *MemoryStore) ReadURL(id string) (string, error) {
	intID, err := strconv.ParseUint(id, store.base, 64)
	if err != nil {
		return "", fmt.Errorf("invalid id %q: %w", id, ErrBadInput)
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	if val, ok := store.db[intID]; ok {
		return val, nil
	}

	return "", ErrNotFound
}

// readDB from file
func (store *MemoryStore) readDB(doLock bool) error {
	file, err := os.OpenFile(store.dbFilePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error reading db at %q: %w", store.dbFilePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := gob.NewDecoder(file)

	if doLock {
		store.mu.Lock()
		defer store.mu.Unlock()
	}

	err = decoder.Decode(&store.db)
	if err != nil && err != io.EOF {
		return fmt.Errorf("decode error: %w", err)
	}

	return nil
}

// writeDB db to file
func (store *MemoryStore) writeDB(doLock bool) error {
	file, err := os.OpenFile(store.dbFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error writing db at %q: %w", store.dbFilePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := gob.NewEncoder(file)

	if doLock {
		store.mu.RLock()
		defer store.mu.RUnlock()
	}

	err = encoder.Encode(&store.db)
	if err != nil && err != io.EOF {
		return fmt.Errorf("encode error: %w", err)
	}

	return nil
}
