package api

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

func TestBatchDeleteHandler(t *testing.T) {
	type args struct {
		user string
		body string
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
			"delete ok",
			args{
				user: "test1",
				body: `["ok1", "ok2"]`,
			},
			want{
				code: http.StatusAccepted,
			},
		},
		{
			"delete bad input",
			args{
				user: "test1",
				body: `["bad", "input"]`,
			},
			want{
				code: http.StatusBadRequest,
				body: `{"error":"bad input"}`,
			},
		},
		{
			"delete 500",
			args{
				user: "test1",
				body: `["error", "500"]`,
			},
			want{
				code: http.StatusInternalServerError,
				body: `{"error":"internal"}`,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := storemock.NewMockBatchRemover(ctrl)
	s.EXPECT().BatchRemove("test1", "ok1", "ok2").Return(nil)
	s.EXPECT().BatchRemove("test1", "bad", "input").Return(store.ErrBadInput)
	s.EXPECT().BatchRemove("test1", "error", "500").Return(errors.New("internal"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(tt.args.body))
			request = request.WithContext(context.WithValue(request.Context(), handler.ContextKeyUID{}, tt.args.user))
			request.Header.Set("Content-Type", "application/json")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := BatchRemoveHandler(s)
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
