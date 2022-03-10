package api

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/handler"
	"shortener/internal/app/service/store"
	storemock "shortener/internal/app/service/store/mock"
	"testing"
)

func TestUserDataHandler(t *testing.T) {
	type args struct {
		user      string
		storeData []store.Record
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"read ok",
			args{
				storeData: []store.Record{
					{
						ShortURL:    "http://short",
						OriginalURL: "http://long",
					},
				},
			},
			want{
				code: http.StatusOK,
				body: `[{"short_url":"http://short","original_url":"http://long"}]`,
			},
		},
		{
			"read no data",
			args{},
			want{
				code: http.StatusNoContent,
				body: `[]`,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storemock.NewMockStore(ctrl)
			s.EXPECT().ReadUserData(tt.args.user).Return(tt.args.storeData)

			request := httptest.NewRequest("GET", "/api/user/urls", nil)
			request = request.WithContext(context.WithValue(request.Context(), handler.ContextKeyUID{}, tt.args.user))
			request.Header.Set("Content-Type", "application/json")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := UserDataHandler(s)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			resBody, _ := ioutil.ReadAll(res.Body)
			assert.Equal(
				t,
				tt.want.code,
				res.StatusCode,
				"Expected status code %d, got %d\nBody was: %s",
				tt.want.code,
				w.Code,
				resBody,
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
