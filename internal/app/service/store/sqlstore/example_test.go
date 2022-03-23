package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
)

func Example() {
	db := getDb()
	defer db.Close()

	s, err := New(db, WithBaseURL("http://localhost:8080"))
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = s.Start(); err != nil {
		fmt.Println(err)
		return
	}

	if url, err := s.WriteURL("http://somelongurl.test/foo/bar", "user1"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(url)
	}

	if err := s.Stop(); err != nil {
		fmt.Println(err)
		return
	}

	// Output: http://localhost:8080/1
}

func getDb() *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery("INSERT INTO").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(1),
	)
	mock.ExpectClose()
	return db
}
