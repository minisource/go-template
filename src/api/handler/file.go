package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/minisource/template_go/api/dto"
	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/dependency"
	"github.com/minisource/template_go/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minisource/go-common/http/helper"
	"github.com/minisource/go-common/logging"
)

type FileHandler struct {
	usecase *usecase.FileUsecase
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		usecase: usecase.NewFileUsecase(cfg, dependency.GetFileRepository(cfg)),
	}
}

// CreateFile godoc
// @Summary Create a file
// @Description Create a file
// @Tags Files
// @Accept x-www-form-urlencoded
// @produces json
// @Param file formData dto.UploadFileRequest true "Create a file"
// @Param file formData file true "Create a file"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.FileResponse} "File response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/files/ [post]
// @Security AuthBearer
func (h *FileHandler) Create(c *fiber.Ctx) error {
	upload := dto.UploadFileRequest{}

	// Fiber does not have ShouldBind; use BodyParser or custom binding
	if err := c.BodyParser(&upload); err != nil {
		resp := helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err)
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	// Get the file manually
	file, err := c.FormFile("file")
	if err != nil {
		resp := helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err)
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	req := dto.CreateFileRequest{
		Description: upload.Description,
		MimeType:    file.Header.Get("Content-Type"),
		Directory:   "uploads",
	}

	req.Name, err = saveUploadedFile(file, req.Directory)
	if err != nil {
		resp := helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err)
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(resp)
	}

	res, err := h.usecase.Create(c.Context(), dto.ToCreateFile(req))
	if err != nil {
		resp := helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err)
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(resp)
	}

	return c.Status(fiber.StatusCreated).JSON(helper.GenerateBaseResponse(res, true, helper.Success))
}


// UpdateFile godoc
// @Summary Update a file
// @Description Update a file
// @Tags Files
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Param Request body dto.UpdateFileRequest true "Update a file"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.FileResponse} "File response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/files/{id} [put]
// @Security AuthBearer
func (h *FileHandler) Update(c *fiber.Ctx) error {
	return Update(c, dto.ToUpdateFile, dto.ToFileResponse, h.usecase.Update)
}


// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file
// @Tags Files
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} helper.BaseHttpResponse "response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/files/{id} [delete]
// @Security AuthBearer
func (h *FileHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id == 0 {
		resp := helper.GenerateBaseResponse(nil, false, helper.ValidationError)
		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	file, err := h.usecase.GetById(c.Context(), id)
	if err != nil {
		logger.Error(logging.IO, logging.RemoveFile, err.Error(), nil)
		resp := helper.GenerateBaseResponse(nil, false, helper.NotFoundError)
		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	err = os.Remove(fmt.Sprintf("%s/%s", file.Directory, file.Name))
	if err != nil {
		logger.Error(logging.IO, logging.RemoveFile, err.Error(), nil)
		resp := helper.GenerateBaseResponse(nil, false, helper.InternalError)
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	err = h.usecase.Delete(c.Context(), id)
	if err != nil {
		resp := helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err)
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(helper.GenerateBaseResponse(nil, true, helper.Success))
}


// GetFile godoc
// @Summary Get a file
// @Description Get a file
// @Tags Files
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.FileResponse} "File response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/files/{id} [get]
// @Security AuthBearer
func (h *FileHandler) GetById(c *fiber.Ctx) error {
	return GetById(c, dto.ToFileResponse, h.usecase.GetById)
}

// GetFiles godoc
// @Summary Get Files
// @Description Get Files
// @Tags Files
// @Accept json
// @produces json
// @Param Request body filter.PaginationInputWithFilter true "Request"
// @Success 200 {object} helper.BaseHttpResponse{result=filter.PagedList[dto.FileResponse]} "File response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/files/get-by-filter [post]
// @Security AuthBearer
func (h *FileHandler) GetByFilter(c *fiber.Ctx) error {
	return GetByFilter(c, dto.ToFileResponse, h.usecase.GetByFilter)
}

func saveUploadedFile(file *multipart.FileHeader, directory string) (string, error) {
	// test.txt -> 95239855629856.txt
	randFileName := uuid.New()
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileName := file.Filename
	fileNameArr := strings.Split(fileName, ".")
	fileExt := fileNameArr[len(fileNameArr)-1]
	fileName = fmt.Sprintf("%s.%s", randFileName, fileExt)
	dst := fmt.Sprintf("%s/%s", directory, fileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
