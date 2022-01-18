package webserver

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/jixindatech/sqlaudit/pkg/config"
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"github.com/jixindatech/sqlaudit/pkg/webserver/api"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/tylerb/graceful"
	"net/http"
	"time"
)

type ApiServer struct {
	cfg         *config.Config
	webAddr     string
	webUser     string
	webPassword string
	*echo.Echo
}

func NewApiServer(cfg *config.Config) (*ApiServer, error) {
	var err error

	s := new(ApiServer)
	s.cfg = cfg
	s.webAddr = cfg.WebAddr
	s.webUser = cfg.WebUser
	s.webPassword = cfg.WebPassword

	s.Echo = echo.New()

	return s, err
}

func (s *ApiServer) Run() {
	s.RegisterWebRoute()
	s.Echo.Server.Addr = s.webAddr
	s.Echo.Validator = &CustomValidator{validator: validator.New()}

	graceful.ListenAndServe(s.Echo.Server, 5*time.Second)
}

func (s *ApiServer) RegisterWebRoute() {
	//s.Use(mw.Logger())
	s.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: `{"remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
		Output: golog.GetLogger(),
	}))
	s.Use(mw.Recover())
	s.Use(mw.CORS())

	s.GET("/ping", api.Ping)
	s.POST("/login", s.Login)
	s.File("/", "dashboard/dist/index.html")
	s.File("/favicon.ico", "dashboard/dist/favicon.ico")
	s.Static("/static", "dashboard/dist/static")

	g := s.Group("/api/v1")
	g.Use(mw.JWT([]byte("secret")))

	g.POST("/rule", api.AddRule)
	g.DELETE("/rule/:id", api.DeleteRule)
	g.PUT("/rule/:id", api.UpdateRule)
	g.GET("/rule/:id", api.GetRule)
	g.GET("/rule", api.GetRules)

	g.GET("/event", api.GetEvents)
	g.GET("/event/info", api.GetEventInfo)
	/*
		g.POST("/fingerprint", api.AddFingerPrint)
		g.DELETE("/fingerprint/:id", api.DeleteFingerPrint)
		g.PUT("/fingerprint/:id", api.UpdateFingerPrint)
		g.GET("/fingerprint/:id", api.GetFingerPrint)
		g.GET("/fingerprint", api.GetFingerPrints)
	*/
}

type userLoginForm struct {
	Username string `json:"username" form:"username" validate:"required,max=254"`
	Password string `json:"password" form:"password" validate:"required,max=254"`
}

func (s *ApiServer) Login(c echo.Context) (err error) {
	form := new(userLoginForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err = c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if form.Username == s.webUser && form.Password == s.webPassword {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)
		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "s.webUser"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}
