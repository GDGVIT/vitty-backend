package pkg

import (
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AdminSvc struct {
	Models []ModelInterface
	WebApp *echo.Echo
}

func NewAdminSvc() *AdminSvc {
	svc := &AdminSvc{}
	svc.init()
	return svc
}

func (s *AdminSvc) init() {
	s.WebApp = echo.New()
	// Logger
	s.WebApp.Use(middleware.Logger())
	s.WebApp.Use(s.AddSVCToEchoContext)

	_, currentFile, _, _ := runtime.Caller(0) // Get the path of the current file
	currentDir := filepath.Dir(currentFile)
	templatesPath := filepath.Join(currentDir, "../templates")

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob(templatesPath + "/*.html")),
	}
	s.WebApp.Renderer = renderer
	s.WebApp.GET("", func(c echo.Context) error {
		fmt.Println("Hello World")
		return c.JSON(200, map[string]interface{}{
			"message": "Hello World",
		})
	})
	AdminHandler(s.WebApp)
}

func (s *AdminSvc) Register(model ModelInterface) {
	s.Models = append(s.Models, model)
}

func (s *AdminSvc) GetModels() []ModelInterface {
	return s.Models
}

func (s *AdminSvc) GetModel(modelName string) ModelInterface {
	for _, model := range s.Models {
		if reflect.TypeOf(model).Elem().Name() == modelName {
			return model
		}
	}
	return nil
}

func GetModelName(model ModelInterface) string {
	return reflect.TypeOf(model).Elem().Name()
}

func GetFields(model ModelInterface) map[string]string {
	fields := make(map[string]string)
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() == reflect.Struct {
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			fields[field.Name] = field.Type.Name()
		}
	}

	return fields
}

func GetPKValue(item interface{}, pkField string) interface{} {
	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	fieldValue := itemValue.FieldByName(pkField)
	if !fieldValue.IsValid() {
		// Handle the case where the field does not exist
		fmt.Println("Field does not exist")
		return nil
	}

	return fieldValue.Interface()
}
