package json

import (
	"errors"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteHandler(t *testing.T) {
	type args struct {
		store       store.Store
		contentType string
		body        string
	}
	type want struct {
		code int
		body string
	}

	s := &store.Mock{}
	s.On("WriteURL", "https://example.org").Return("http://localhost/bar", nil)
	s.On("WriteURL", "").Return("", errors.New("bad url"))
	s.On("WriteURL", "bad").Return("", errors.New("bad url"))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"write ok",
			args{
				store:       s,
				contentType: "application/json",
				body:        "{\"url\":\"https://example.org\"}",
			},
			want{
				code: http.StatusCreated,
				body: "{\"result\":\"http://localhost/bar\"}",
			},
		},
		{
			"write empty",
			args{
				store:       s,
				contentType: "application/json",
				body:        "",
			},
			want{
				code: http.StatusBadRequest,
				body: "json read error: unexpected end of JSON input\n",
			},
		},
		{
			"write bad",
			args{
				store:       s,
				contentType: "application/json",
				body:        "{\"url\":\"bad\"}",
			},
			want{
				code: http.StatusBadRequest,
				body: "bad url\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.args.body))
			request.Header.Set("Content-Type", tt.args.contentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := WriteHandler(s)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			resBody, _ := ioutil.ReadAll(res.Body)
			assert.Equal(
				t,
				tt.want.code,
				res.StatusCode,
				"Expected status code %d, got %d",
				tt.want.code,
				w.Code,
			)
			assert.Equal(
				t,
				tt.want.body,
				string(resBody),
				"Expected body %s, got %s",
				tt.want.body,
				string(resBody),
			)
			_ = res.Body.Close()
		})
	}
}
