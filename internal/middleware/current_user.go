package middleware

import "github.com/gofiber/fiber/v2"

type User struct {
	ID   int
	Role string
}

func GetUser(c *fiber.Ctx) *User {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return nil
	}

	role, ok := c.Locals("role").(string)
	if !ok {
		return nil
	}

	return &User{
		ID:   userID,
		Role: role,
	}
}
