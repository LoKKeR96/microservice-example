package resolvers

import (
	"github.com/lokker96/microservice_example/infrastructure/container"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	C container.Container
}
