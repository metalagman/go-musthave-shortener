package shortener

import (
	"fmt"
	"strconv"
	"sync"
)

var _ Store = (*MemoryStore)(nil)

type MemoryStore struct {
	sync.Mutex
	listenAddr string
	baseURL    string
	counter    uint64
	base       int
	db         MemoryDB
}

type MemoryDB map[uint64]string

func NewMemoryStore(listenAddr string, baseURL string) *MemoryStore {
	return &MemoryStore{
		Mutex:      sync.Mutex{},
		counter:    30,
		listenAddr: listenAddr,
		baseURL:    baseURL,
		base:       36,
		db:         make(MemoryDB),
	}
}

func (store *MemoryStore) SetDB(db MemoryDB) {
	store.Lock()
	defer store.Unlock()
	store.db = db
}

func (store *MemoryStore) GetDB() MemoryDB {
	store.Lock()
	defer store.Unlock()
	return store.db
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
