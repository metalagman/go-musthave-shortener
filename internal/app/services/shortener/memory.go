package shortener

import (
	"errors"
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
	urls       map[uint64]string
}

func NewMemoryStore(listenAddr string, baseURL string) *MemoryStore {
	return &MemoryStore{
		Mutex:      sync.Mutex{},
		counter:    30,
		listenAddr: listenAddr,
		baseURL:    baseURL,
		base:       36,
		urls:       make(map[uint64]string),
	}
}

func (store *MemoryStore) WriteURL(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("empty url")
	}

	store.Lock()
	defer store.Unlock()

	store.counter++
	store.urls[store.counter] = url
	id := strconv.FormatUint(store.counter, store.base)

	return fmt.Sprintf("%s/%s", store.baseURL, id), nil
}

func (store *MemoryStore) ReadURL(id string) (string, error) {
	intID, err := strconv.ParseUint(id, store.base, 64)
	if err != nil {
		return "", err
	}

	store.Lock()
	defer store.Unlock()

	if val, ok := store.urls[intID]; ok {
		return val, nil
	}

	return "", errors.New("not found")
}
