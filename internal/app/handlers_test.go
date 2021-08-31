package app

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ShortenerServiceMock struct {
	mock.Mock
}

func (m *ShortenerServiceMock) WriteURL(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func (m *ShortenerServiceMock) ReadURL(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func TestReadHandler(t *testing.T) {
	type args struct {
		svc  ShortenerService
		path string
	}
	type want struct {
		code        int
		redirectURL string
	}

	svcMock := &ShortenerServiceMock{}
	svcMock.On("ReadURL", "test1").Return("https://example.org", nil)
	svcMock.On("ReadURL", "").Return("", errors.New("empty id"))
	svcMock.On("ReadURL", "missing").Return("", errors.New("missing id"))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"read ok",
			args{
				svc:  svcMock,
				path: "/test1",
			},
			want{
				code:        http.StatusTemporaryRedirect,
				redirectURL: "https://example.org",
			},
		},
		{
			"read empty",
			args{
				svc:  svcMock,
				path: "/",
			},
			want{
				code: http.StatusBadRequest,
			},
		},
		{
			"read missing",
			args{
				svc:  svcMock,
				path: "/",
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
			h := ReadHandler(svcMock)
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
		svc  ShortenerService
		body string
	}
	type want struct {
		code     int
		response string
		body     string
	}

	svcMock := &ShortenerServiceMock{}
	svcMock.On("WriteURL", "https://example.org").Return("http://localhost/bar", nil)
	svcMock.On("WriteURL", "").Return("", errors.New("bad url"))
	svcMock.On("WriteURL", "bad").Return("", errors.New("bad url"))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"write ok",
			args{
				svc:  svcMock,
				body: "https://example.org",
			},
			want{
				code: http.StatusCreated,
				body: "http://localhost/bar",
			},
		},
		{
			"write empty",
			args{
				svc:  svcMock,
				body: "",
			},
			want{
				code: http.StatusBadRequest,
				body: "bad url\n",
			},
		},
		{
			"write bad",
			args{
				svc:  svcMock,
				body: "bad",
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
			h := WriteHandler(svcMock)
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
