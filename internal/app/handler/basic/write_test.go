package basic

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/handler"
	"shortener/internal/app/service/store"
	storemock "shortener/internal/app/service/store/mock"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := storemock.NewMockStore(ctrl)
	s.EXPECT().WriteURL("https://example.org", "test").Return("http://localhost/bar", nil)
	s.EXPECT().WriteURL("", "test").Return("", errors.New("bad url"))
	s.EXPECT().WriteURL("bad", "test").Return("", errors.New("bad url"))
	s.EXPECT().WriteURL("https://example.org/conflict", "test").Return("", &store.ConflictError{
		ExistingURL: "https://example.org/non-conflict",
	})

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
		{
			"conflict",
			args{
				store: s,
				body:  "https://example.org/conflict",
			},
			want{
				code: http.StatusConflict,
				body: "https://example.org/non-conflict",
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
