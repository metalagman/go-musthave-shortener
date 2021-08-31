package json

import (
	"errors"
	"github.com/russianlagman/go-musthave-shortener/internal/app"
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

func TestWriteHandler(t *testing.T) {
	type args struct {
		svc  app.ShortenerService
		body string
	}
	type want struct {
		code int
		body string
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
				body: "{\"url\":\"https://example.org\"}",
			},
			want{
				code: http.StatusCreated,
				body: "{\"result\":\"http://localhost/bar\"}",
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
				body: "json read error: unexpected end of JSON input\n",
			},
		},
		{
			"write bad",
			args{
				svc:  svcMock,
				body: "{\"url\":\"bad\"}",
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
