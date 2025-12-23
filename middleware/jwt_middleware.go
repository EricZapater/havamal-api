// middleware/jwt.go
package middleware

import (
	"havamal-api/config"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthUser struct {
	ID         string
	TenantID   string
	Username   string
	Email      string
	CustomerID string
	IsAdmin    bool
}

func SetupJWT(cfg config.Config) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "turniq",
		Key:         []byte(cfg.Auth.Secret),
		Timeout:     cfg.Auth.TTL,
		MaxRefresh:  cfg.Auth.TTL,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*AuthUser); ok {
				return jwt.MapClaims{
					"id":          v.ID,
					"username":    v.Username,
					"email":       v.Email,
					"customer_id": v.CustomerID,
					"is_admin":    v.IsAdmin,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &AuthUser{
				ID:         getStringClaim(claims, "id"),
				Username:   getStringClaim(claims, "username"),
				Email:      getStringClaim(claims, "email"),
				CustomerID: getStringClaim(claims, "customer_id"),
				IsAdmin:    getBoolClaim(claims, "is_admin"),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if user, ok := data.(*AuthUser); ok {
				return user.ID != ""
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if v, ok := claims[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getBoolClaim(claims jwt.MapClaims, key string) bool {
	if v, ok := claims[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// Helpers

func GetUser(c *gin.Context) *AuthUser {
	if v, exists := c.Get("id"); exists {
		if user, ok := v.(*AuthUser); ok {
			return user
		}
	}
	return nil
}

func GetUserID(c *gin.Context) string {
	if user := GetUser(c); user != nil {
		return user.ID
	}
	return ""
}

func GetCustomerID(c *gin.Context) string {
	if user := GetUser(c); user != nil {
		return user.CustomerID
	}
	return ""
}

func IsAdmin(c *gin.Context) bool {
	if user := GetUser(c); user != nil {
		return user.IsAdmin
	}
	return false
}