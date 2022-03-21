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

func TestReadHandler(t *testing.T) {
	type args struct {
		store store.Store
		path  string
	}
	type want struct {
		code        int
		redirectURL string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := storemock.NewMockStore(ctrl)
	s.EXPECT().ReadURL("test1").Return("https://example.org", nil)
	s.EXPECT().ReadURL("").Return("", errors.New("empty id"))
	s.EXPECT().ReadURL("missing").Return("", errors.New("missing id"))
	s.EXPECT().ReadURL("deleted").Return("", store.ErrDeleted)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"read ok",
			args{
				store: s,
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
				store: s,
				path:  "/",
			},
			want{
				code: http.StatusBadRequest,
			},
		},
		{
			"read missing",
			args{
				store: s,
				path:  "/missing",
			},
			want{
				code: http.StatusBadRequest,
			},
		},
		{
			"read deleted",
			args{
				store: s,
				path:  "/deleted",
			},
			want{
				code: http.StatusGone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", tt.args.path, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := ReadHandler(s)
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
