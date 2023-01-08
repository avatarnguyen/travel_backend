package helpers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func MatchUserTypeToUuid(c *fiber.Ctx, userId string) (err error) {
	userType := c.GetRespHeader("user_type")
	uid := c.GetRespHeader("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}
func CheckUserType(c *fiber.Ctx, role string) (err error) {
	userType := c.GetRespHeader("user_type")

	err = nil
	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	return err
}
