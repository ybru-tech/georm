package examples

import (
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/ybru-tech/georm/examples/testutil"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	conn, closer := testutil.InitTempDB()
	db = conn

	code := m.Run()

	closer()

	os.Exit(code)
}
