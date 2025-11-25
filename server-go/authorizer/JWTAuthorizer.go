package authorizer

import (
	"context"
	"strings"

	"go.temporal.io/server/common/authorization"
)

type jwtAuthorizer struct{}

func NewJwtAuthorizer() authorization.Authorizer {
	return &jwtAuthorizer{}
}

var decisionAllow = authorization.Result{Decision: authorization.DecisionAllow}
var decisionDeny = authorization.Result{Decision: authorization.DecisionDeny}

func (a *jwtAuthorizer) Authorize(_ context.Context, claims *authorization.Claims, target *authorization.CallTarget) (authorization.Result, error) {

	if target.Namespace == "temporal-system" {
			// Change this before using in a prod env. this is allowing all temporal-systems namespace use to be un-authenticated.
		return decisionAllow, nil
	}

	// Allow Admin and writers
	if claims != nil && claims.System&(authorization.RoleAdmin|authorization.RoleWriter) != 0 {
		return decisionAllow, nil
	}

	// Allow all calls except UpdateNamespace through when claim mapper isn't invoked
	// Claim mapper is skipped unless TLS is configured or an auth token is passed
	if claims == nil && !strings.Contains(target.APIName, "UpdateNamespace") {
		return decisionAllow, nil
	}

	// For other namespaces, deny "UpdateNamespace" API unless the caller has a writer role in it
	if strings.Contains(target.APIName, "UpdateNamespace") {
		if claims != nil && claims.Namespaces[target.Namespace]&authorization.RoleWriter != 0 {
			return decisionAllow, nil
		} else {
			return decisionDeny, nil
		}
	}

	// Allow all other requests
	return decisionAllow, nil
}

