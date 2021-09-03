package app

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMemoryShortenerService_ReadURL(t *testing.T) {
	type fields struct {
		counter uint64
		urls    map[uint64]string
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
				urls: map[uint64]string{
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
				urls: map[uint64]string{
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
				urls: map[uint64]string{
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
			svc := &MemoryShortenerService{
				Mutex:      sync.Mutex{},
				listenAddr: "localhost:8080",
				base:       10,
				counter:    tt.fields.counter,
				urls:       tt.fields.urls,
			}
			got, err := svc.ReadURL(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryShortenerService_WriteURL(t *testing.T) {
	type fields struct {
		counter uint64
		urls    map[uint64]string
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
				urls:    map[uint64]string{},
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
				urls:    map[uint64]string{},
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
			svc := &MemoryShortenerService{
				Mutex:      sync.Mutex{},
				listenAddr: "localhost:8080",
				baseUrl:    "http://localhost:8080",
				base:       10,
				counter:    tt.fields.counter,
				urls:       tt.fields.urls,
			}
			got, err := svc.WriteURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMemoryShortenerService(t *testing.T) {
	svc := NewMemoryShortenerService("localhost:8080", "http://localhost:8080")
	assert.NotNil(t, svc)
	assert.Equal(t, svc.listenAddr, "localhost:8080")
	expInterface := (*ShortenerService)(nil)
	assert.Implementsf(t, expInterface, svc, "Interface %v must be implemented in %v", expInterface, svc)
}
