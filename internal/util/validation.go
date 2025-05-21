package util

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaginationQuery struct {
	Page  string `query:"page" validate:"required,numeric"`
	Limit string `query:"limit" validate:"required,numeric"`
}

type ImageSaveOptions struct {
	FieldName string
	BasePath  string
	PublicURL string
	MaxSizeMB int64
}

func Validate[T any](data T) map[string]string {
	err := validator.New().Struct(data)
	res := map[string]string{}

	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			res[v.StructField()] = TranslateTag(v)
		}
	}

	return res
}

func TranslateTag(fd validator.FieldError) string {
	switch fd.ActualTag() {
	case "required":
		return fmt.Sprintf("field %s wajib diisi", fd.StructField())
	case "min":
		return fmt.Sprintf("field %s size minimal %s", fd.StructField(), fd.Param())
	case "unique":
		return fmt.Sprintf("field %s harus unique", fd.StructField())

	default:
		return "validasi gagal"
	}
}

func SafePaginationParams(ctx *fiber.Ctx) (int, int, map[string]string) {
	pageStr := ctx.Query("page", "1")
	limitStr := ctx.Query("limit", "10")

	// Validasi hanya angka, tolak jika ada karakter aneh
	isNumeric := regexp.MustCompile(`^\d+$`).MatchString
	errors := map[string]string{}

	if !isNumeric(pageStr) {
		errors["page"] = "parameter page harus berupa angka "
	}
	if !isNumeric(limitStr) {
		errors["limit"] = "parameter limit harus berupa angka "
	}

	if len(errors) > 0 {
		return 0, 0, errors
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	return page, limit, nil
}

func ProcessAndSaveImageFile(ctx *fiber.Ctx, opts ImageSaveOptions) (string, error) {
	file, err := ctx.FormFile(opts.FieldName)
	if err != nil {
		return "", errors.New("file tidak ditemukan")
	}

	// Validasi ukuran
	if file.Size > opts.MaxSizeMB*1024*1024 {
		return "", fmt.Errorf("ukuran file maksimal %dMB", opts.MaxSizeMB)
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExt(ext) {
		return "", errors.New("format file tidak valid (hanya jpg, jpeg, png, webp)")
	}

	// Simpan file
	filename := uuid.NewString() + ext
	savePath := filepath.Join(opts.BasePath, filename)
	if err := ctx.SaveFile(file, savePath); err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}

	publicURL := path.Join(opts.PublicURL, filename)
	return publicURL, nil
}

func isValidImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	default:
		return false
	}
}
