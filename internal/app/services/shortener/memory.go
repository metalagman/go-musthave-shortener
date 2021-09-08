package shortener

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

var _ Store = (*MemoryStore)(nil)

type MemoryStore struct {
	sync.Mutex
	listenAddr string
	baseURL    string
	dbFilePath string
	counter    uint64
	base       int
	db         MemoryDB
}

type MemoryDB map[uint64]string

func NewMemoryStore(listenAddr string, baseURL string, dbFilePath string) *MemoryStore {
	return &MemoryStore{
		Mutex:      sync.Mutex{},
		counter:    30,
		listenAddr: listenAddr,
		baseURL:    baseURL,
		dbFilePath: dbFilePath,
		base:       36,
		db:         make(MemoryDB),
	}
}

func (store *MemoryStore) WriteURL(url string) (string, error) {
	if err := ValidateURL(url); err != nil {
		return "", err
	}

	store.Lock()
	defer store.Unlock()

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

	store.Lock()
	defer store.Unlock()

	if val, ok := store.db[intID]; ok {
		return val, nil
	}

	return "", ErrNotFound
}

// ReadDB from file
func (store *MemoryStore) ReadDB() error {
	file, err := os.OpenFile(store.dbFilePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error reading db at %q: %w", store.dbFilePath, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := gob.NewDecoder(file)

	store.Lock()
	defer store.Unlock()

	err = decoder.Decode(&store.db)
	if err != nil && err != io.EOF {
		return fmt.Errorf("decode error: %w", err)
	}

	return nil
}

// WriteDB db to file
func (store *MemoryStore) WriteDB() error {
	file, err := os.OpenFile(store.dbFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error writing db at %q: %w", store.dbFilePath, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := gob.NewEncoder(file)

	store.Lock()
	defer store.Unlock()

	err = encoder.Encode(&store.db)
	if err != nil && err != io.EOF {
		return fmt.Errorf("encode error: %w", err)
	}

	return nil
}
