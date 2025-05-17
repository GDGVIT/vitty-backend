package pkg

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"reflect"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/middleware"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/cli/commands"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func AdminHandler(app fiber.Router) {
	group := app.Group("/admin")
	group.Use(middleware.JWTAuthMiddleware)
	group.Use(middleware.IsAdminMiddleware)
	group.Post("/empty-classrooms/seed", seedEmptyClassrooms)

	// group.Get("", GetModelsView)
	// group.Get("/:model", GetModelView)
	// group.Get("/:model/create", CreateItemView)
	// group.Post("/:model/create", CreateItem)
	// group.Get("/:model/:id", GetItemView)
	// group.Get("/:model/:id/edit", UpdateItemView)
	// group.Put("/:model/:id", UpdateItem)
	// group.Delete("/:model/:id", DeleteItem)
}

func GetModelsView(c echo.Context) error {
	fmt.Println("list models")
	svc := c.Get("svc").(*AdminSvc)
	models := svc.GetModels()
	// Model names -
	var names []string
	for _, model := range models {
		names = append(names, GetModelName(model))
	}
	return c.Render(http.StatusOK, "list-models.html", map[string]interface{}{
		"model_names": names,
	})
}

func GetModelView(c echo.Context) error {
	fmt.Println("get model")
	svc := c.Get("svc").(*AdminSvc)
	modelName := c.Param("model")
	model := svc.GetModel(modelName)
	fmt.Println("Model", model)
	if model == nil {
		fmt.Println("Model not found")
		return c.Redirect(http.StatusFound, "/")
	}
	items, err := model.GetAll()
	fmt.Println("Items", items)
	var itemList []map[string]interface{}
	for _, item := range items {
		itemList = append(itemList, map[string]interface{}{
			"pk_field":  model.GetPKField(),
			"pk_value":  GetPKValue(item, model.GetPKField()),
			"item_name": item.GetItemName(),
		})
	}
	fmt.Println("item list", itemList)
	if err != nil {
		fmt.Println("Error", err)
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Render(http.StatusOK, "list-items.html", map[string]interface{}{
		"model":  modelName,
		"fields": GetFields(model),
		"items":  itemList,
	})
}

func CreateItemView(c echo.Context) error {
	svc := c.Get("svc").(*AdminSvc)
	modelName := c.Param("model")
	model := svc.GetModel(modelName)
	if model == nil {
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Render(http.StatusOK, "create-model.html", map[string]interface{}{
		"model":  GetModelName(model),
		"fields": GetFields(model),
	})
}

func CreateItem(c echo.Context) error {
	svc := c.Get("svc").(*AdminSvc)
	modelName := c.Param("model")
	model := svc.GetModel(modelName)
	if model == nil {
		return c.Redirect(http.StatusFound, "/")
	}
	model.Create(model)
	return c.Redirect(http.StatusFound, "/"+modelName)
}

func GetItemView(c echo.Context) error {
	fmt.Println("get item")
	svc := c.Get("svc").(*AdminSvc)
	modelName := c.Param("model")
	model := svc.GetModel(modelName)
	if model == nil {
		fmt.Println("Model not found")
		return c.Redirect(http.StatusFound, "/")
	}
	id := c.Param("id")
	item, err := model.Get(id)
	if err != nil {
		fmt.Println("Error", err)
		return c.Redirect(http.StatusFound, "/"+modelName)
	}
	fmt.Println("Item", item)

	itemMap := make(map[string]string)
	for fieldName := range GetFields(model) {
		itemMap[fieldName] = fmt.Sprintf("%v", reflect.ValueOf(item).Elem().FieldByName(fieldName))
	}
	fmt.Println("Item Map", itemMap)

	return c.Render(http.StatusOK, "get-item.html", map[string]interface{}{
		"model":     GetModelName(model),
		"fields":    GetFields(model),
		"item_name": item.GetItemName(),
		"pk_value":  GetPKValue(item, model.GetPKField()),
		"item":      &item,
		"item_map":  itemMap,
	})
}

func UpdateItemView(c echo.Context) error {
	fmt.Println("update item")
	svc := c.Get("svc").(*AdminSvc)
	modelName := c.Param("model")
	model := svc.GetModel(modelName)
	if model == nil {
		fmt.Println("Model not found")
		return c.Redirect(http.StatusFound, "/")
	}
	id := c.Param("id")
	item, err := model.Get(id)
	if err != nil {
		fmt.Println("Error", err)
		return c.Redirect(http.StatusFound, "/"+modelName)
	}
	fmt.Println("Item", item)

	itemMap := make(map[string]string)
	for fieldName := range GetFields(model) {
		itemMap[fieldName] = fmt.Sprintf("%v", reflect.ValueOf(item).Elem().FieldByName(fieldName))
	}
	fmt.Println("Item Map", itemMap)

	return c.Render(http.StatusOK, "update-item.html", map[string]interface{}{
		"model":     GetModelName(model),
		"fields":    GetFields(model),
		"item_name": item.GetItemName(),
		"pk_value":  GetPKValue(item, model.GetPKField()),
		"item":      &item,
		"item_map":  itemMap,
	})
}

func UpdateItem(c echo.Context) error {
	svc := c.Get("svc").(*AdminSvc)
	modelName := c.Param("model")
	model := svc.GetModel(modelName)
	if model == nil {
		return c.Redirect(http.StatusFound, "/")
	}
	id := c.Param("id")
	item, err := model.Get(id)
	if err != nil {
		c.Logger().Error(err)
		return c.Redirect(http.StatusFound, "/"+modelName)
	}
	err = model.Update(item)
	if err != nil {
		c.Logger().Error(err)
		return c.Redirect(http.StatusFound, "/"+modelName)
	}
	return c.Redirect(http.StatusFound, "/"+modelName+"/"+id)
}

func DeleteItem(c echo.Context) error {
	svc := c.Get("svc").(*AdminSvc)
	modelName := c.Param("model")
	model := svc.GetModel(modelName)
	if model == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Model not found",
		})
	}
	id := c.Param("id")
	item, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Item not found",
		})
	}
	model.Delete(item)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Item deleted",
	})

}

func seedEmptyClassrooms(c *fiber.Ctx) error {
	err := commands.GenerateEmptyRooms(&cli.Context{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "classrooms seed failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "classrooms has been initialized",
	})
}
