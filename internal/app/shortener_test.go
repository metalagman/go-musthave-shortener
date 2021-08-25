package app

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMemoryShortenerService_Read(t *testing.T) {
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
				Mutex:   sync.Mutex{},
				addr:    "localhost:8080",
				base:    10,
				counter: tt.fields.counter,
				urls:    tt.fields.urls,
			}
			got, err := svc.Read(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryShortenerService_Write(t *testing.T) {
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
				Mutex:   sync.Mutex{},
				addr:    "localhost:8080",
				base:    10,
				counter: tt.fields.counter,
				urls:    tt.fields.urls,
			}
			got, err := svc.Write(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMemoryShortenerService(t *testing.T) {
	svc := NewMemoryShortenerService("localhost:8080")
	assert.NotNil(t, svc)
	assert.Equal(t, svc.addr, "localhost:8080")
}
