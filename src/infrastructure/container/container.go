package container

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// Define the Container structure
type Container struct {
	Ctx context.Context // Context for the container
	db  *gorm.DB        // Database connection instance
}

// NewContainer function creates and returns a new Container instance
func NewContainer() (*Container, error) {
	ctx := context.Background()

	dbPassword, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error on getting db secrets")
	}

	zone, _ := time.Now().Zone()

	dbDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		string(dbPassword),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		zone,
	)

	dbConnection, err := NewDBConnection(dbDSN)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error on creating new db connection")
	}

	// Return a new Container instance with the database connection set
	return &Container{
		db:  dbConnection,
		Ctx: ctx,
	}, nil
}

func (c *Container) HandleShutdown() {
	// Add here any memory cleanup method
}
