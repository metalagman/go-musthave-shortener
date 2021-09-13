package shortener

import (
	"errors"
	"testing"
)

func Test_ValidateUrl(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "empty url",
			args: args{
				str: "",
			},
			wantErr: ErrEmptyInput,
		},
		{
			name: "some word",
			args: args{
				str: "localhost",
			},
			wantErr: ErrBadInput,
		},
		{
			name: "proper url",
			args: args{
				str: "https://ya.ru",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateURL(tt.args.str); !errors.Is(got, tt.wantErr) {
				t.Errorf("ValidateURL() = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}
