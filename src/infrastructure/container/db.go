package container

import (
	"github.com/lokker96/microservice_example/domain/entity"

	"github.com/palantir/stacktrace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConnection(dsn string) (*gorm.DB, error) {
	EntityTypes := []interface{}{
		&entity.Message{},
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, stacktrace.Propagate(err, "error on opening db connection")
	}

	// tables are dropped just for demonstration purposes, dropping tables in production would cause data to be lost.
	for _, entityType := range EntityTypes {
		if err := db.Migrator().DropTable(entityType); err != nil {
			return nil, stacktrace.Propagate(err, "error on dropping table")
		}
	}

	// iterate through entities and use gorm tags to build tables
	for _, entityType := range EntityTypes {
		if err := db.AutoMigrate(entityType); err != nil {
			return nil, stacktrace.Propagate(err, "error on auto migrating table")
		}
	}

	return db, nil
}
