package db

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBName = "testdb"
	testDBUser = "root"
	testDBPass = "root"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	c, host, port := RunDBContainer(ctx)
	defer func() {
		_ = c.Terminate(ctx)
	}()

	var err error
	db, err = StartDBStore(StartUpOptions{
		DBHost:         host,
		DBPort:         port,
		DBName:         testDBName,
		DBUsername:     testDBUser,
		DBPassword:     testDBPass,
		SkipMigrations: false,
	})
	if err != nil {
		log.Fatalf("Failed to start the database: %s", err)
	}
	code := m.Run()
	os.Exit(code)
}

func RunDBContainer(ctx context.Context) (dbC testcontainers.Container, host string, port int) {
	basePort, err := nat.NewPort("tcp", "3306")
	if err != nil {
		log.Fatal(err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "mariadb:10.6",
		ExposedPorts: []string{"3306"},
		Env:          map[string]string{"MARIADB_ROOT_PASSWORD": testDBPass, "MARIADB_DATABASE": testDBName},
		WaitingFor: wait.ForAll(
			wait.ForLog("mariadbd: ready for connections."),
			wait.ForLog("port: 3306"),
		),
	}

	dbC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	host, err = dbC.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	natPort, err := dbC.MappedPort(ctx, basePort)
	if err != nil {
		log.Fatalf("Could not get test container port: %s", err)
	}

	port, err = strconv.Atoi(string(natPort.Port()))
	if err != nil {
		log.Fatalf("Could not parse test container port: %s", err)
	}

	return
}
