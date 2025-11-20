package plugins

import (
	"context"
	"errors"
	"strings"

	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
)

type (
	// SimpleTokenAuthorizer implements authorization.Authorizer
	SimpleTokenAuthorizer struct {
		logger log.Logger
	}

	// SimpleTokenClaimMapper implements authorization.ClaimMapper
	SimpleTokenClaimMapper struct {
		logger log.Logger
	}
)

func NewSimpleTokenAuthorizer(logger log.Logger) *SimpleTokenAuthorizer {
	return &SimpleTokenAuthorizer{logger: logger}
}

func (a *SimpleTokenAuthorizer) Authorize(ctx context.Context, caller *authorization.Claims, target *authorization.CallTarget) (authorization.Result, error) {
	// For this example, we'll trust the claims mapped by the ClaimMapper.
	// In a real scenario, you might check specific permissions here if not covered by RBAC.
	// If the caller has claims (meaning the token was valid), we allow.
	if caller != nil && caller.Subject != "" {
		return authorization.Result{Decision: authorization.DecisionAllow}, nil
	}
	// If no claims, deny.
	return authorization.Result{Decision: authorization.DecisionDeny, Reason: "Unauthorized: No valid claims"}, nil
}

func NewSimpleTokenClaimMapper(logger log.Logger) *SimpleTokenClaimMapper {
	return &SimpleTokenClaimMapper{logger: logger}
}

func (m *SimpleTokenClaimMapper) GetClaims(authInfo *authorization.AuthInfo) (*authorization.Claims, error) {
	// Extract the token from the Authorization header
	authHeader := authInfo.AuthToken
	if authHeader == "" {
		return nil, nil // No auth token provided
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("invalid authorization header format")
	}
	token := parts[1]

	// INSECURE: For demonstration, we are NOT validating the JWT signature against an IDP.
	// In production, you MUST verify the token signature using your OIDC provider's public keys.
	// Here we just check if the token is "test-token" or pretend to parse it.
	
	// Let's assume the token is valid if it's not empty for this example, 
	// or strictly check for our mock token.
	// Real implementation would use a JWT library to parse and validate.
	
	claims := &authorization.Claims{
		Subject: "test-worker", // This would come from the 'sub' claim
		Namespaces: map[string]authorization.Role{
			"default": authorization.RoleWriter, // Grant Writer role on 'default' namespace
		},
	}

	// Log for debugging
	m.logger.Info("Mapped claims", tag.NewStringTag("subject", claims.Subject), tag.NewStringTag("token_snippet", token[:5]+"..."))

	return claims, nil
}
