package sqlstore

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func Test_HealthCheck(t *testing.T) {
	mdb, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer func() {
		_ = mdb.Close()
	}()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ping ok",
			wantErr: false,
		},
		{
			name:    "ping fail",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Store{
				db: mdb,
			}

			if tt.wantErr {
				mock.ExpectPing().WillReturnError(errors.New("ping error"))
			} else {
				mock.ExpectPing()
			}

			err := r.HealthCheck()
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
