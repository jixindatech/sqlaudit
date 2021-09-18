package api

import (
	"github.com/jixindatech/sqlaudit/pkg/webserver/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type queryEventForm struct {
	Page  int    `json:"page" form:"page" query:"page" validate:"required,gte=1"`
	Size  int    `json:"size" form:"size" query:"size" validate:"required,min=1,max=50"`
	Name  string `json:"name" form:"name" query:"name" validate:"omitempty,max=254"`
	Type  int    `json:"type" form:"type" query:"type" validate:"omitempty,min=1,max=3"`
	Sql   string `json:"sql"  form:"sql" query:"sql"  validate:"omitempty,max=254"`
	Db    string `json:"db"   form:"db" query:"db"  validate:"omitempty,max=254"`
	IP    string `json:"ip"   form:"ip" query:"ip" validate:"omitempty,ip"`
	User  string `json:"user" form:"user" query:"user" validate:"omitempty,max=254"`
	Op    int    `json:"op"   form:"op" query:"op" validate:"omitempty,gte=1,lte=15"`
	Start int64  `json:"start" form:"start" query:"start" validate:"omitempty"`
	End   int64  `json:"end" form:"end" query:"end" validate:"omitempty"`
}

func GetEvents(c echo.Context) (err error) {
	form := new(queryEventForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	query := make(map[string]interface{})
	query["name"] = form.Name
	query["type"] = form.Type
	query["sql"] = form.Sql
	query["db"] = form.Db
	query["ip"] = form.IP
	query["user"] = form.User
	query["op"] = form.Op
	query["start"] = form.Start
	query["end"] = form.End

	data, err := models.GetEvents(query, form.Page, form.Size)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, data)
}

type queryEventInfoForm struct {
	Db    string `json:"db" query:"db" validate:"omitempty,max=254"`
	Start int64  `json:"start" query:"start" validate:"omitempty"`
	End   int64  `json:"end" query:"end" validate:"omitempty"`
}

func GetEventInfo(c echo.Context) (err error) {
	// query := make(map[string]interface{})
	form := new(queryEventInfoForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	query := make(map[string]interface{})
	query["db"] = form.Db
	query["start"] = form.Start / 1000
	query["end"] = form.End / 1000
	query["interval"] = 10

	data, err := models.GetEventInfo(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, data)
}
