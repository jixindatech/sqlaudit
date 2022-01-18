package api

import (
	"fmt"
	"github.com/jixindatech/sqlaudit/pkg/apps/mysql"
	"github.com/jixindatech/sqlaudit/pkg/webserver/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ruleForm struct {
	Name     string `json:"name" form:"name" validate:"required,max=254"`
	Time     int    `json:"time" form:"time" validate:"omitempty,min=0"`
	Type     int    `json:"type" form:"type" validate:"required,gte=1,lte=2"`
	User     string `json:"user" form:"user" validate:"omitempty,max=254"`
	IP       string `json:"ip"   form:"ip" validate:"omitempty,ip"`
	Db       string `json:"db"   form:"db" validate:"omitempty,max=254"`
	RuleType int    `json:"ruletype" form:"ruletype" validate:"required,min=1,max=2"`
	Op       int    `json:"op"   form:"op" validate:"omitempty,gte=0,lte=15"`
	Match    int    `json:"match"   form:"match" validate:"omitempty,gte=1,lte=2"`
	Priority int    `json:"priority" form:"priority" validate:"omitempty,gte=0"`
	Sql      string `json:"sql"   form:"sql" validate:"omitempty,max=254"`
	Alert    int    `json:"alert"   form:"alert" validate:"omitempty,gte=0,lte=1"`
	Remark   string `json:"remark" validate:"max=254"`
}

func AddRule(c echo.Context) (err error) {
	form := new(ruleForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(form); err != nil {
		fmt.Println("err:", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if form.RuleType == 1 {
		if form.Op == 0 {
			return c.JSON(http.StatusBadRequest, "invalid ruletype parameters")
		}
		if form.Match > 0 && len(form.Sql) == 0 {
			return c.JSON(http.StatusBadRequest, "invalid ruletype parameters")
		}
	} else {
		if len(form.Sql) == 0 {
			return c.JSON(http.StatusBadRequest, "invalid ruletype parameters")
		}
		form.Op = 0
		form.Match = 0
		form.Priority = 0
	}

	data := make(map[string]interface{})
	data["name"] = form.Name
	data["time"] = form.Time
	data["type"] = form.Type
	data["user"] = form.User
	data["ip"] = form.IP
	data["db"] = form.Db
	data["ruletype"] = form.RuleType
	data["op"] = form.Op
	data["sql"] = form.Sql
	data["match"] = form.Match
	data["priority"] = form.Priority
	data["alert"] = form.Alert
	data["remark"] = form.Remark

	err = models.AddRule(data)
	if err != nil {
		return err
	}

	err = updateRuleConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, "ok")
}

func DeleteRule(c echo.Context) error {
	_id := c.Param("id")
	if _id == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	id, err := strconv.Atoi(_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := models.DeleteRule(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err = updateRuleConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, "ok")
}

func UpdateRule(c echo.Context) (err error) {
	_id := c.Param("id")
	if _id == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	id, err := strconv.Atoi(_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	form := new(ruleForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	data := make(map[string]interface{})
	data["name"] = form.Name
	data["time"] = form.Time
	data["type"] = form.Type
	data["user"] = form.User
	data["ip"] = form.IP
	data["db"] = form.Db
	data["op"] = form.Op
	data["sql"] = form.Sql
	data["match"] = form.Match
	data["priority"] = form.Priority
	data["alert"] = form.Alert
	data["remark"] = form.Remark

	err = models.UpdateRule(uint(id), data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err = updateRuleConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}

func GetRule(c echo.Context) error {
	_id := c.Param("id")
	if _id == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	id, err := strconv.Atoi(_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	rule, err := models.GetRule(uint(id))
	if err != nil || rule.ID == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, rule)
}

type queryRuleForm struct {
	Page int    `json:"page" query:"page" validate:"required,gte=1"`
	Size int    `json:"size" query:"size" validate:"required,min=1,max=50"`
	Sort int    `json:"sort" query:"sort" validate:"required,gte=1,lte=4"`
	Name string `json:"name" query:"name" validate:"omitempty,max=254"`
}

func GetRules(c echo.Context) (err error) {
	form := new(queryRuleForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	query := make(map[string]interface{})
	query["name"] = form.Name
	query["sort"] = form.Sort

	rules, err := models.GetRules(query, form.Page, form.Size)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	count, err := models.GetRuleCount(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	data := make(map[string]interface{})
	data["data"] = rules
	data["count"] = count

	return c.JSON(http.StatusOK, data)
}

func updateRuleConfig() error {
	err := mysql.ParserSqlRules()
	if err != nil {
		return err
	}

	return nil
}
