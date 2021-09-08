package basic

import (
	"errors"
	"github.com/russianlagman/go-musthave-shortener/internal/app/services/shortener"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadHandler(t *testing.T) {
	type args struct {
		store shortener.Store
		path  string
	}
	type want struct {
		code        int
		redirectURL string
	}

	store := &shortener.StoreMock{}
	store.On("ReadURL", "test1").Return("https://example.org", nil)
	store.On("ReadURL", "").Return("", errors.New("empty id"))
	store.On("ReadURL", "missing").Return("", errors.New("missing id"))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"read ok",
			args{
				store: store,
				path:  "/test1",
			},
			want{
				code:        http.StatusTemporaryRedirect,
				redirectURL: "https://example.org",
			},
		},
		{
			"read empty",
			args{
				store: store,
				path:  "/",
			},
			want{
				code: http.StatusBadRequest,
			},
		},
		{
			"read missing",
			args{
				store: store,
				path:  "/",
			},
			want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", tt.args.path, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := ReadHandler(store)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(
				t,
				tt.want.code,
				res.StatusCode,
				"Expected status code %d, got %d",
				tt.want.code,
				w.Code,
			)
			if len(tt.want.redirectURL) > 0 {
				assert.Equal(
					t,
					tt.want.redirectURL,
					res.Header.Get("Location"),
					"Expected location header %s, got %s",
					tt.want.redirectURL,
					res.Header.Get("Location"),
				)
			}
			_ = res.Body.Close()
		})
	}
}

func TestWriteHandler(t *testing.T) {
	type args struct {
		store shortener.Store
		body  string
	}
	type want struct {
		code int
		body string
	}

	store := &shortener.StoreMock{}
	store.On("WriteURL", "https://example.org").Return("http://localhost/bar", nil)
	store.On("WriteURL", "").Return("", errors.New("bad url"))
	store.On("WriteURL", "bad").Return("", errors.New("bad url"))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"write ok",
			args{
				store: store,
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
				store: store,
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
				store: store,
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
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := WriteHandler(store)
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

func Test_IsUrl(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty url",
			args: args{
				str: "",
			},
			want: false,
		},
		{
			name: "some word",
			args: args{
				str: "localhost",
			},
			want: false,
		},
		{
			name: "proper url",
			args: args{
				str: "https://ya.ru",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURL(tt.args.str); got != tt.want {
				t.Errorf("isUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
