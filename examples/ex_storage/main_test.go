package ex_storage

import (
	"log"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/ybru-tech/georm/examples/testutil"
)

var (
	db      *gorm.DB
	storage *Storage
)

func TestMain(m *testing.M) {
	conn, closer := testutil.InitTempDB()

	db = conn
	storage = NewStorage(conn)

	if err := storage.MigrationTables(); err != nil {
		panic(err)
	}

	code := m.Run()

	if err := storage.DropTables(); err != nil {
		log.Print(err) // если тут будет паника, не сработает closer
	}

	closer()

	os.Exit(code)
}
