package memorystore

import (
	"fmt"
	"time"
)

func Example() {
	s := NewStore(
		WithBaseURL("http://localhost:8080"),
		WithListenAddr("localhost:8080"),
		WithFilePath("test_urls.data"),
		WithFlushInterval(time.Second),
	)

	if err := s.Start(); err != nil {
		fmt.Println(err)
		return
	}

	if url, err := s.WriteURL("http://somelongurl.test/foo/bar", "user1"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(url)
	}

	time.Sleep(2 * time.Second)

	if err := s.Stop(); err != nil {
		fmt.Println(err)
		return
	}

	// Output: http://localhost:8080/1
}
