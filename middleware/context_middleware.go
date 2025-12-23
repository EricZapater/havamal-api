package middleware

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextMiddleware injects "is_admin" and "customer_id" from Gin context (JWT claims)
// into the standard Request context, so services can access them via ctx.Value()
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user == nil {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		// Inject is_admin
		ctx = context.WithValue(ctx, "is_admin", user.IsAdmin)

		// Inject customer_id (UUID)
		if user.CustomerID != "" {
			parsedID, err := uuid.Parse(user.CustomerID)
			if err == nil {
				ctx = context.WithValue(ctx, "customer_id", parsedID)
			}
		}

		// Update request with the new context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// UserFromContext retrieves the user ID from the context. NOT USED by the services directly,
// but useful helper if needed. Services generally look for "is_admin" and "customer_id".
func GetCustomerIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value("customer_id")
	if id, ok := val.(uuid.UUID); ok {
		return id, nil
	}
	return uuid.Nil, errors.New("customer_id not found in context")
}

func GetIsAdminFromCtx(ctx context.Context) (bool, error) {
	val := ctx.Value("is_admin")
	if isAdmin, ok := val.(bool); ok {
		return isAdmin, nil
	}
	return false, errors.New("is_admin not found in context")
}
