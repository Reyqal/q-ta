package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

// JWTProtected mengembalikan middleware JWT yang memvalidasi token dari header Authorization.
func JWTProtected(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau sudah kedaluwarsa",
				"data":    nil,
			})
		},
		ContextKey: "jwt",
	})
}

// ExtractUserID mengambil user ID dari JWT claims.
func ExtractUserID(c *fiber.Ctx) uint {
	token := c.Locals("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	return uint(userID)
}

// ExtractUserRole mengambil role dari JWT claims.
func ExtractUserRole(c *fiber.Ctx) string {
	token := c.Locals("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	return role
}

// ExtractClaims mengambil user ID dan role sekaligus dari JWT claims.
func ExtractClaims(c *fiber.Ctx) (uint, string) {
	return ExtractUserID(c), ExtractUserRole(c)
}
