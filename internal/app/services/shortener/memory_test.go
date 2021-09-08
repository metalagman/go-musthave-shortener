package shortener

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMemoryStore_ReadURL(t *testing.T) {
	type fields struct {
		counter uint64
		db      MemoryDB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"read existing",
			fields{
				db: MemoryDB{
					1: "https://example.org",
				},
			},
			args{
				id: "1",
			},
			"https://example.org",
			false,
		},
		{
			"read empty",
			fields{
				db: MemoryDB{
					1: "https://example.org",
				},
			},
			args{
				id: "",
			},
			"",
			true,
		},
		{
			"read missing",
			fields{
				db: MemoryDB{
					1: "https://example.org",
				},
			},
			args{
				id: "2",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemoryStore{
				Mutex:      sync.Mutex{},
				listenAddr: "localhost:8080",
				base:       10,
				counter:    tt.fields.counter,
				db:         tt.fields.db,
			}
			got, err := store.ReadURL(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadURL() got = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_WriteURL(t *testing.T) {
	type fields struct {
		counter uint64
		urls    MemoryDB
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"write url",
			fields{
				counter: 0,
				urls:    MemoryDB{},
			},
			args{
				url: "https://example.org",
			},
			"http://localhost:8080/1",
			false,
		},
		{
			"write empty url",
			fields{
				counter: 0,
				urls:    MemoryDB{},
			},
			args{
				url: "",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemoryStore{
				Mutex:      sync.Mutex{},
				listenAddr: "localhost:8080",
				baseURL:    "http://localhost:8080",
				base:       10,
				counter:    tt.fields.counter,
				db:         tt.fields.urls,
			}
			got, err := store.WriteURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteURL() got = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func TestNewMemoryShortenerService(t *testing.T) {
	store := NewMemoryStore("localhost:8080", "http://localhost:8080", "urls.gob")
	assert.NotNil(t, store)
	assert.Equal(t, store.listenAddr, "localhost:8080")
}
