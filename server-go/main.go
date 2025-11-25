package main

import (
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/temporal"

	"github.com/example/temporal-custom-server/authorizer"
	"github.com/example/temporal-custom-server/claims"
)

func main() {
	// Create a new server instance
	logger := log.NewCLILogger()

	s, err := temporal.NewServer(
		temporal.ForServices(temporal.DefaultServices),
		temporal.WithAuthorizer(authorizer.NewJwtAuthorizer()),
		temporal.WithClaimMapper(func(cfg *config.Config) authorization.ClaimMapper {
			return claims.NewTLSClaimMapper(logger)
		}),
	)
	if err != nil {
		panic(err)
	}

	// Start the server
	if err := s.Start(); err != nil {
		panic(err)
	}

	// Block forever (or handle signals properly)
	select {}
}
