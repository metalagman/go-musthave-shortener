package app

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type ShortenerService interface {
	Write(url string) (string, error)
	Read(id string) (string, error)
}

type MemoryShortenerService struct {
	sync.Mutex
	addr    string
	counter uint64
	base    int
	urls    map[uint64]string
}

func NewMemoryShortenerService(addr string) *MemoryShortenerService {
	return &MemoryShortenerService{
		Mutex:   sync.Mutex{},
		counter: 30,
		addr:    addr,
		base:    36,
		urls:    make(map[uint64]string),
	}
}

func (svc *MemoryShortenerService) Write(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("empty url")
	}

	svc.Lock()
	defer svc.Unlock()

	svc.counter++
	svc.urls[svc.counter] = url
	id := strconv.FormatUint(svc.counter, svc.base)

	return fmt.Sprintf("http://%s/%s", svc.addr, id), nil
}

func (svc *MemoryShortenerService) Read(id string) (string, error) {
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
