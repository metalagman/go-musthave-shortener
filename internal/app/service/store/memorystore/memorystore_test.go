package memorystore

import (
	"log"
	"time"
)

func Example() {
	s := NewStore(
		WithBaseURL("http://localhost:8080"),
		WithListenAddr("localhost:8080"),
		WithFilePath("/tmp/urls.data"),
		WithFlushInterval(time.Second),
	)

	if err := s.Start(); err != nil {
		log.Println(err)
		return
	}

	if url, err := s.WriteURL("http://somelongurl.test/foo/bar", "user1"); err != nil {
		log.Println(err)
		return
	} else {
		log.Println(url)
	}

	time.Sleep(2 * time.Second)

	if err := s.Stop(); err != nil {
		log.Println(err)
		return
	}

	// Output: "http://localhost:8080/1
}
