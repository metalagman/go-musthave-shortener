package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var shortener *ShortenerService

type ShortenerService struct {
	sync.Mutex
	addr    string
	counter uint64
	base    int
	urls    map[uint64]string
}

func NewShortenerService(addr string) *ShortenerService {
	return &ShortenerService{
		Mutex:   sync.Mutex{},
		counter: 30,
		addr:    addr,
		base:    36,
		urls:    make(map[uint64]string),
	}
}

func (svc *ShortenerService) write(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("empty url")
	}
	svc.Lock()
	svc.counter++
	svc.urls[svc.counter] = url
	id := strconv.FormatUint(svc.counter, svc.base)
	svc.Unlock()
	return fmt.Sprintf("http://%s/%s", svc.addr, id), nil
}

func (svc *ShortenerService) read(id string) (string, error) {
	intId, err := strconv.ParseUint(id, svc.base, 64)
	if err != nil {
		return "", err
	}
	svc.Lock()
	val, ok := svc.urls[intId]
	svc.Unlock()
	if ok {
		return val, nil
	} else {
		return "", errors.New("not found")
	}
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func Shortener(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		u, err := shortener.read(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u := r.FormValue("url")
		if !IsUrl(u) {
			http.Error(w, "bad url", http.StatusBadRequest)
			return
		}
		id, err := shortener.write(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(201)
		_, _ = w.Write([]byte(id))
	default:
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}

func main() {
	addr := "localhost:8080"
	shortener = NewShortenerService(addr)
	// маршрутизация запросов обработчику
	http.HandleFunc("/", Shortener)

	fmt.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
