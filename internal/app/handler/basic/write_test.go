package basic

import (
	"context"
	"errors"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler"
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
		store store.Store
		body  string
	}
	type want struct {
		code int
		body string
	}

	s := &store.Mock{}
	s.On("WriteUserURL", "https://example.org", "").Return("http://localhost/bar", nil)
	s.On("WriteUserURL", "https://example.org", "test").Return("http://localhost/bar", nil)
	s.On("WriteUserURL", "", "test").Return("", errors.New("bad url"))
	s.On("WriteUserURL", "bad", "test").Return("", errors.New("bad url"))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"write ok",
			args{
				store: s,
				body:  "https://example.org",
			},
			want{
				code: http.StatusCreated,
				body: "http://localhost/bar",
			},
		},
		{
			"write empty",
			args{
				store: s,
				body:  "",
			},
			want{
				code: http.StatusBadRequest,
				body: "bad url\n",
			},
		},
		{
			"write bad",
			args{
				store: s,
				body:  "bad",
			},
			want{
				code: http.StatusBadRequest,
				body: "bad url\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.args.body))
			request = request.WithContext(context.WithValue(request.Context(), handler.ContextKeyUID{}, "test"))
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
