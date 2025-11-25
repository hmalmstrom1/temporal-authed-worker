package claims

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
)

const (
	adminCN      = "temporal"
	clientDomain = "client"
)

type tlsClaimMapper struct {
	logger log.Logger
}

func NewTLSClaimMapper(logger log.Logger) authorization.ClaimMapper {
	return &tlsClaimMapper{logger: logger}
}

func (a *tlsClaimMapper) GetClaims(authInfo *authorization.AuthInfo) (*authorization.Claims, error) {
	claims := authorization.Claims{}

	// Check for Bearer token
	if authInfo.AuthToken != "" {
		parts := strings.Split(authInfo.AuthToken, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString := parts[1]

			// Parse the token (without verification for now as we don't have the key configured here yet,
			// but in production you MUST verify the signature)
			// TODO: Inject public key or JWKS provider for verification
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				a.logger.Error("Failed to parse token", tag.Error(err))
				return nil, err
			}

			if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
				// Map 'sub' to Subject
				if sub, ok := mapClaims["sub"].(string); ok {
					claims.Subject = sub
				}

				// Map 'scp' or custom claims to Roles
				// Example: if scp contains "temporal:writer", grant writer role
				if scp, ok := mapClaims["scp"].([]interface{}); ok {
					for _, s := range scp {
						if str, ok := s.(string); ok && str == "temporal:writer" {
							if claims.Namespaces == nil {
								claims.Namespaces = make(map[string]authorization.Role)
							}
							claims.Namespaces["default"] = authorization.RoleWriter
						}
					}
				}

				// For this specific user request, we grant writer on default if token is present
				if claims.Namespaces == nil {
					claims.Namespaces = make(map[string]authorization.Role)
				}
				claims.Namespaces["default"] = authorization.RoleWriter

				a.logger.Info("Mapped claims from token", tag.NewStringTag("subject", claims.Subject), tag.NewStringTag("namespaces", "default:writer"))
				return &claims, nil
			}
		}
	}

	if authInfo.TLSSubject != nil {
		cn := authInfo.TLSSubject.CommonName
		if cn == adminCN {
			claims.System = authorization.RoleAdmin
		} else if ns, domain, ok := strings.Cut(cn, "."); ok && domain == clientDomain {
			claims.Namespaces = map[string]authorization.Role{ns: authorization.RoleWriter}
		}
		a.logger.Info("Mapped claims from TLS", tag.NewStringTag("cn", cn))
	}
	return &claims, nil
}
