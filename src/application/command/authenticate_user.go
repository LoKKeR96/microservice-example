package command

import (
	"context"

	"github.com/lokker96/microservice_example/domain/service"
	"github.com/palantir/stacktrace"
)

type CreateUserTokenCommand struct {
	ctx                       context.Context
	userAuthenticationService service.UserAuthenticationService
}

func NewCreateUserTokenCommand(
	ctx context.Context,
	userAuthenticationService service.UserAuthenticationService) CreateUserTokenCommand {
	return CreateUserTokenCommand{
		ctx:                       ctx,
		userAuthenticationService: userAuthenticationService,
	}
}

func (cs *CreateUserTokenCommand) Do(username string, password string) (string, error) {

	// calling domain message creator service
	token, err := cs.userAuthenticationService.Authenticate(username, password)
	if err != nil {
		return "", stacktrace.Propagate(err, "error on authentication user")
	}

	return token, nil
}
