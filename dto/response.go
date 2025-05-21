package dto

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Response[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PagedResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
	Meta    Meta   `json:"meta"`
}

func CreateResponseError(message string) Response[string] {
	return Response[string]{
		Code:    "99",
		Message: message,
		Data:    "",
	}
}

func CreateResponseSuccess[T any](data T) Response[T] {
	return Response[T]{
		Code:    "00",
		Message: "success",
		Data:    data,
	}
}

func CreateResponseErrorData(message string, data map[string]string) Response[map[string]string] {
	return Response[map[string]string]{
		Code:    "99",
		Message: message,
		Data:    data,
	}
}

func PaginateAndRespond[T any](ctx *fiber.Ctx, data []T, total int) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	return ctx.JSON(PagedResponse[T]{
		Code:    "00",
		Message: "success",
		Data:    data,
		Meta: Meta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}
