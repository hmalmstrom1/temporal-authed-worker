package main

import (
	"github.com/example/temporal-custom-server/plugins"
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/temporal"
)

func main() {
	// Create a new server instance
	logger := log.NewCLILogger()

	s, err := temporal.NewServer(
		temporal.ForServices(temporal.DefaultServices),
		temporal.WithAuthorizer(plugins.NewSimpleTokenAuthorizer(logger)),
		temporal.WithClaimMapper(func(cfg *config.Config) authorization.ClaimMapper {
			return plugins.NewSimpleTokenClaimMapper(logger)
		}),
	)
	if err != nil {
		panic(err)
	}

	// Start the server
	if err := s.Start(); err != nil {
		panic(err)
	}
}
