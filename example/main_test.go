package example

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	conn, closer := initTempDB()
	db = conn

	code := m.Run()

	closer()

	os.Exit(code)
}

// InitTempDB - запускает docker с postgres, подключается к бд и возвращает репозиторий, соединенный с новой базой данных
func initTempDB() (db *gorm.DB, closer func()) {
	var (
		user     = "postgres"
		password = "secret"
		dbName   = "eda_sandbox_tests"
	)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs itgit
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgis/postgis",
		Tag:        "16-master",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
			"TZ=Europe/Moscow",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	pgHost, port := getHostPort(resource, "5432/tcp")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", pgHost, port, user, dbName, password)

	err = pool.Retry(func() error {
		db, err = gorm.Open(postgres.Open(dsn))
		if err != nil {
			return fmt.Errorf("could not construct pool: %s", err)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	return db, func() {
		err = pool.Purge(resource)
		if err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func getHostPort(resource *dockertest.Resource, id string) (host, port string) {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		hostPortParts := strings.Split(resource.GetHostPort(id), ":")

		return hostPortParts[0], hostPortParts[1]
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	return u.Hostname(), resource.GetPort(id)
}
