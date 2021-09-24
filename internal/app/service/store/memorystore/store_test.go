package memorystore

import (
	"testing"
)

func TestStore_ReadURL(t *testing.T) {
	type fields struct {
		counter uint64
		db      db
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
				db: db{
					1: {
						"1",
						"https://example.org",
						"http://localhost/1",
						"test",
					},
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
				db: db{
					1: {
						"1",
						"https://example.org",
						"http://localhost/1",
						"test",
					},
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
				db: db{
					1: {
						"1",
						"https://example.org",
						"http://localhost/1",
						"test",
					},
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
			store := &Store{
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

func TestStore_WriteURL(t *testing.T) {
	type fields struct {
		counter uint64
		db      db
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
				db:      db{},
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
				db:      db{},
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
			store := &Store{
				listenAddr: "localhost:8080",
				baseURL:    "http://localhost:8080",
				base:       10,
				counter:    tt.fields.counter,
				db:         tt.fields.db,
			}
			got, err := store.WriteUserURL(tt.args.url, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteUserURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteUserURL() got = %v, wantErr %v", got, tt.want)
			}
		})
	}
}
