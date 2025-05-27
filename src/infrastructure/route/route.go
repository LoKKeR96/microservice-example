package route

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/controller"
	"github.com/lokker96/microservice_example/infrastructure/graph"
	resolvers "github.com/lokker96/microservice_example/infrastructure/graph/resolver"
	"github.com/vektah/gqlparser/v2/ast"

	middleware "github.com/lokker96/microservice_example/infrastructure/http"
)

func SetupRoutes(c container.Container) *echo.Echo {
	e := echo.New()

	// Setting up GraphQL config
	graphConfig := graph.Config{
		Resolvers: &resolvers.Resolver{C: c},
	}

	// Schema Directive - @example
	// graphConfig.Directives.example = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	// 	return next(ctx)
	// }

	// Setting up GraphQL handlers
	graphqlHandler := handler.New(graph.NewExecutableSchema(graphConfig))

	graphqlHandler.AddTransport(transport.Options{})
	graphqlHandler.AddTransport(transport.GET{})
	graphqlHandler.AddTransport(transport.POST{})

	graphqlHandler.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	graphqlHandler.Use(extension.Introspection{})
	graphqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	playgroundHandler := playground.Handler("GraphQL", "/query")

	e.POST("/query", func(ctx echo.Context) error {
		graphqlHandler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	},
	// Need to finish this off for authentication in GraphQL
	// middleware.AuthenticationMiddleware,
	)

	e.GET("/playground", func(ctx echo.Context) error {
		playgroundHandler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})

	// Setting up REST APIs handlers
	// adding custom validator for modularity and scalability
	// separating this logic makes it easier to make validation changes in the future
	e.Validator = &CustomValidator{validator.New()}

	createMessageController := controller.NewCreateMessage()
	deleteMessageController := controller.NewDeleteMessage()
	getAllMessages := controller.NewGetMessagesByFilter()
	getMessageController := controller.NewGetMessage()
	updateMessageController := controller.NewUpdateMessage()

	userAuthenticationController := controller.NewUserAuthentication()

	// REST APIs Routes definitions

	e.POST(
		"/message/create",
		func(ctx echo.Context) error {
			return createMessageController.Create(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.PUT(
		"/message/:uuid/update",
		func(ctx echo.Context) error {
			return updateMessageController.Update(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.DELETE(
		"/message/:uuid/delete",
		func(ctx echo.Context) error {
			return deleteMessageController.Delete(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.POST(
		"/message/search",
		func(ctx echo.Context) error {
			return getAllMessages.Filter(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.GET(
		"/message/:uuid/get",
		func(ctx echo.Context) error {
			return getMessageController.Get(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.POST("/user/login",
		func(ctx echo.Context) error {
			return userAuthenticationController.Authenticate(ctx, c)
		},
	)

	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
