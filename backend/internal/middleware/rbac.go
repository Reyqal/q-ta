package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRole mengembalikan middleware yang hanya mengizinkan akses untuk role tertentu.
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := ExtractUserRole(c)

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses ke resource ini",
			"data":    nil,
		})
	}
}
