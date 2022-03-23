package basic

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/service/store"
	storemock "shortener/internal/app/service/store/mock"
	"testing"
)

func TestPingHandler(t *testing.T) {
	type args struct {
		store  store.HealthChecker
		path   string
		method string
	}
	type want struct {
		code int
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := storemock.NewMockHealthChecker(ctrl)
	s.EXPECT().HealthCheck().Return(nil)
	s.EXPECT().HealthCheck().Return(errors.New("health check error"))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"ping ok",
			args{
				store:  s,
				method: http.MethodGet,
				path:   "/ping",
			},
			want{
				code: http.StatusOK,
			},
		},
		{
			"ping error",
			args{
				store:  s,
				method: http.MethodGet,
				path:   "/ping",
			},
			want{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := PingHandler(s)
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
			_ = res.Body.Close()
		})
	}
}
