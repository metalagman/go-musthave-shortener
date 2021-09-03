package app

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type ShortenerService interface {
	WriteURL(url string) (string, error)
	ReadURL(id string) (string, error)
}

type MemoryShortenerService struct {
	sync.Mutex
	listenAddr string
	baseURL    string
	counter    uint64
	base       int
	urls       map[uint64]string
}

func NewMemoryShortenerService(listenAddr string, baseURL string) *MemoryShortenerService {
	return &MemoryShortenerService{
		Mutex:      sync.Mutex{},
		counter:    30,
		listenAddr: listenAddr,
		baseURL:    baseURL,
		base:       36,
		urls:       make(map[uint64]string),
	}
}

func (svc *MemoryShortenerService) WriteURL(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("empty url")
	}

	svc.Lock()
	defer svc.Unlock()

	svc.counter++
	svc.urls[svc.counter] = url
	id := strconv.FormatUint(svc.counter, svc.base)

	return fmt.Sprintf("%s/%s", svc.baseURL, id), nil
}

func (svc *MemoryShortenerService) ReadURL(id string) (string, error) {
	intID, err := strconv.ParseUint(id, svc.base, 64)
	if err != nil {
		return "", err
	}

	svc.Lock()
	defer svc.Unlock()

	if val, ok := svc.urls[intID]; ok {
		return val, nil
	}

	return "", errors.New("not found")
}
