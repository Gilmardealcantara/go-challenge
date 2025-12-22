//go:build integration

package integration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
	"github.com/mytheresa/go-hiring-challenge/app/api/server"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	container testcontainers.Container
	db        *gorm.DB
	ctx       = context.Background()
)

// IntegrationTestSuite provides database setup for integration tests
type IntegrationTestSuite struct {
	suite.Suite
	DB  *gorm.DB
	mux *http.ServeMux
}

// SetupSuite runs once before all tests in the suite
func (s *IntegrationTestSuite) SetupSuite() {
	var err error
	container, db, err = initializeContainer()
	s.Require().NoError(err)
	s.DB = db

	s.mux = server.Setup(s.DB)
}

// TearDownSuite runs once after all tests in the suite
func (s *IntegrationTestSuite) TearDownSuite() {
	if container != nil {
		if err := container.Terminate(ctx); err != nil {
			s.T().Logf("failed to terminate container: %v", err)
		}
	}
}

// SetupTest runs before each test
func (s *IntegrationTestSuite) SetupTest() {
	err := resetDatabase()
	s.Require().NoError(err)
}

// initializeContainer creates and starts a PostgreSQL container
func initializeContainer() (testcontainers.Container, *gorm.DB, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:17.5",
		// ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		).WithStartupTimeoutDefault(90 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	// Get container port
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		container.Terminate(ctx)
		return nil, nil, err
	}

	// Connect to database with retry
	dsn := fmt.Sprintf("postgres://testuser:testpass@localhost:%s/testdb?sslmode=disable", port.Port())
	var database *gorm.DB
	var lastErr error

	for i := 0; i < 10; i++ {
		database, lastErr = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Discard,
		})
		if lastErr == nil {
			break
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}

	if lastErr != nil {
		container.Terminate(ctx)
		return nil, nil, lastErr
	}

	// Run seed command to initialize database
	if err := runSeed(port.Port()); err != nil {
		container.Terminate(ctx)
		return nil, nil, err
	}

	return container, database, nil
}

// resetDatabase runs the seed command to reset database state
func resetDatabase() error {
	if container == nil {
		return fmt.Errorf("container not initialized")
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return err
	}

	return runSeed(port.Port())
}

// runSeed executes the seed command with the container's database connection
func runSeed(port string) error {
	cmd := exec.Command("go", "run", "cmd/seed/main.go")
	cmd.Dir = "../../"

	// Set environment variables for the seed command
	env := os.Environ()
	env = append(env, "POSTGRES_USER=testuser")
	env = append(env, "POSTGRES_PASSWORD=testpass")
	env = append(env, "POSTGRES_DB=testdb")
	env = append(env, fmt.Sprintf("POSTGRES_PORT=%s", port))
	env = append(env, "POSTGRES_SQL_DIR=sql")

	cmd.Env = env

	// Run the seed command
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("seed command failed: %w\noutput: %s", err, string(output))
	}

	return nil
}
