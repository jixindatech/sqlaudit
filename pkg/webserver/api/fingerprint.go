package api

import (
	"github.com/jixindatech/sqlaudit/pkg/webserver/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type fingerPrintForm struct {
	FingerPrint string `json:"fingerprint"   form:"fingerprint" validate:"required,max=254"`
	Remark      string `json:"remark" validate:"max=254"`
}

func AddFingerPrint(c echo.Context) (err error) {
	form := new(fingerPrintForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	data := make(map[string]interface{})
	data["fingerprint"] = form.FingerPrint
	data["remark"] = form.Remark

	err = models.AddFingerPrint(data)
	if err != nil {
		return err
	}

	err = updateFingerPrintConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, "ok")
}

func DeleteFingerPrint(c echo.Context) error {
	_id := c.Param("id")
	if _id == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	id, err := strconv.Atoi(_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := models.DeleteFingerPrint(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err = updateFingerPrintConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, "ok")
}

func UpdateFingerPrint(c echo.Context) (err error) {
	_id := c.Param("id")
	if _id == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	id, err := strconv.Atoi(_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	form := new(fingerPrintForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	data := make(map[string]interface{})
	data["fingerprint"] = form.FingerPrint
	data["remark"] = form.Remark

	err = models.UpdateFingerPrint(uint(id), data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err = updateFingerPrintConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}

func GetFingerPrint(c echo.Context) error {
	_id := c.Param("id")
	if _id == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	id, err := strconv.Atoi(_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	rule, err := models.GetFingerPrint(uint(id))
	if err != nil || rule.ID == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, rule)
}

type queryFingerPrintForm struct {
	Page int    `json:"page" query:"page" validate:"required,gte=1"`
	Size int    `json:"size" query:"size" validate:"required,min=1,max=50"`
	Sort int    `json:"sort" query:"sort" validate:"required,gte=1,lte=4"`
	Name string `json:"name" query:"name" validate:"omitempty,max=254"`
}

func GetFingerPrints(c echo.Context) (err error) {
	form := new(queryFingerPrintForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	query := make(map[string]interface{})
	fingerprints, err := models.GetFingerPrints(query, form.Page, form.Size)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	count, err := models.GetFingerPrintCount(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	data := make(map[string]interface{})
	data["data"] = fingerprints
	data["count"] = count

	return c.JSON(http.StatusOK, data)
}

func updateFingerPrintConfig() error {
	return nil
}
