// A generic base handler can that we can use for CRUD operations
//
// To use simple crud can see user handler and year handler

package handler

import (
	"context"
	"strconv"

	"github.com/minisource/template_go/config"
	"github.com/gofiber/fiber/v2"
	"github.com/minisource/go-common/filter"
	"github.com/minisource/go-common/http/helper"
	"github.com/minisource/go-common/logging"
)

var logger = logging.NewLogger(&config.GetConfig().Logger)

// Create an entity
// TRequest: Http request body
// TUInput: Usecase method input that mapped from TRequest with TUInput := mapper(TRequest)
// TUOutput: Usecase function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// requestMapper: this function map endpoint input to usecase input
// responseMapper: this function map usecase output to endpoint output
// usecaseCreate: usecase Create method
func Create[TRequest any, TUInput any, TUOutput any, TResponse any](
	c *fiber.Ctx,
	requestMapper func(req TRequest) TUInput,
	responseMapper func(req TUOutput) TResponse,
	usecaseCreate func(ctx context.Context, req TUInput) (TUOutput, error),
) error {
	request := new(TRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
	}

	usecaseInput := requestMapper(*request)

	usecaseResult, err := usecaseCreate(c.Context(), usecaseInput)
	if err != nil {
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	response := responseMapper(usecaseResult)

	return c.Status(fiber.StatusCreated).JSON(
		helper.GenerateBaseResponse(response, true, 0),
	)
}

// Update an entity
// TRequest: Http request body
// TUInput: Use case method input that mapped from TRequest with TUInput := mapper(TRequest)
// TUOutput: Use case function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// requestMapper: this function map endpoint input to usecase input
// responseMapper: this function map usecase output to endpoint output
// usecaseUpdate: usecase Update method
func Update[TRequest any, TUInput any, TUOutput any, TResponse any](
	c *fiber.Ctx,
	requestMapper func(req TRequest) TUInput,
	responseMapper func(req TUOutput) TResponse,
	usecaseUpdate func(ctx context.Context, id int, req TUInput) (TUOutput, error),
) error {
	// Bind path param
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			helper.GenerateBaseResponse(nil, false, helper.ValidationError),
		)
	}

	// Bind request body
	request := new(TRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
	}

	// Map to usecase input
	usecaseInput := requestMapper(*request)

	// Call usecase
	usecaseResult, err := usecaseUpdate(c.Context(), id, usecaseInput)
	if err != nil {
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	// Map and return response
	response := responseMapper(usecaseResult)
	return c.Status(fiber.StatusOK).JSON(
		helper.GenerateBaseResponse(response, true, 0),
	)
}


func Delete(c *fiber.Ctx, usecaseDelete func(ctx context.Context, id int) error) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			helper.GenerateBaseResponse(nil, false, helper.ValidationError),
		)
	}

	err = usecaseDelete(c.Context(), id)
	if err != nil {
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.GenerateBaseResponse(nil, true, 0),
	)
}


// Get an entity
// TUOutput: Usecase function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// responseMapper: this function map usecase output to endpoint output
// usecaseGet: usecase Get method
func GetById[TUOutput any, TResponse any](
	c *fiber.Ctx,
	responseMapper func(req TUOutput) TResponse,
	usecaseGet func(ctx context.Context, id int) (TUOutput, error),
) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			helper.GenerateBaseResponse(nil, false, helper.ValidationError),
		)
	}

	usecaseResult, err := usecaseGet(c.Context(), id)
	if err != nil {
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	response := responseMapper(usecaseResult)
	return c.Status(fiber.StatusOK).JSON(
		helper.GenerateBaseResponse(response, true, 0),
	)
}

// Get entities by filter
// TUOutput: Usecase function output
// TResponse: Http response body that mapped from TUOutput with TResponse := mapper(TUOutput)
// responseMapper: this function map usecase output to endpoint output
// usecaseList: usecase GetByFilter method
func GetByFilter[TUOutput any, TResponse any](
	c *fiber.Ctx,
	responseMapper func(req TUOutput) TResponse,
	usecaseList func(ctx context.Context, req filter.PaginationInputWithFilter) (*filter.PagedList[TUOutput], error),
) error {
	req := new(filter.PaginationInputWithFilter)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
	}

	usecaseResult, err := usecaseList(c.Context(), *req)
	if err != nil {
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	response := filter.PagedList[TResponse]{
		PageNumber:      usecaseResult.PageNumber,
		PageSize:        usecaseResult.PageSize,
		TotalRows:       usecaseResult.TotalRows,
		TotalPages:      usecaseResult.TotalPages,
		HasPreviousPage: usecaseResult.HasPreviousPage,
		HasNextPage:     usecaseResult.HasNextPage,
	}

	items := []TResponse{}
	for _, item := range *usecaseResult.Items {
		items = append(items, responseMapper(item))
	}
	response.Items = &items

	return c.Status(fiber.StatusOK).JSON(
		helper.GenerateBaseResponse(response, true, 0),
	)
}
