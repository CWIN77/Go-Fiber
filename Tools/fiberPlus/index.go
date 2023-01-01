package fiberPlus

import (
	"github.com/gofiber/fiber/v2"
)

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func GetParams(c *fiber.Ctx, paramList []string) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	err := c.BodyParser(&p)
	for _, value := range paramList {
		if p[value] == "" || p[value] == nil {
			errMessage := "Require parameter : "
			for _, value := range paramList {
				errMessage += value + ", "
			}
			err = &CustomError{Message: errMessage}
		}
	}
	return p, err
}
