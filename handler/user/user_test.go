package user

import (
	"encoding/json"
	"fmt"
	"github.com/MuxiKeStack/muxiK-StackBackend/config"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/router/middleware"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"testing"
)

var (
	g           *gin.Engine
	tokenString string
	username    string
	password    string
	sid         string
)

func TestMain(m *testing.M) {

	// init config
	if err := config.Init(""); err != nil {
		panic(err)
	}
	// init db
	model.DB.Init()
	defer model.DB.Close()

	os.Exit(m.Run())
}
func TestLogin(t *testing.T) {
	g := getRouter(true)

	uri := "/login"
	u := CreateRequest{
		model.LoginModel{
			Sid:      "2018212576",
			Password: "Yu@14796825550",
		},
	}
	jsonByte, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	w := util.PerformRequestWithBody(http.MethodPost, g, uri, jsonByte, "")

	// 读取响应body,获取tokenString
	var data LoginResponse

	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test error: Get LoginResponse Error:%s", err.Error())
	}
	tokenString = data.Data.Token
	fmt.Println(tokenString, data.Data.IsNew)
	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", w.Code)
	}
}

// Helper function to create a router during testing
func getRouter(withRouter bool) *gin.Engine {
	g = gin.New()
	if withRouter {
		loadRouters(
			// Cores.
			g,

			// Middlwares.
			middleware.Logging(),
			middleware.RequestId(),
		)
	}
	return g
}

// Load loads the middlewares, routes, handlers about Test
func loadRouters(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// api for authentication functionalities
	g.POST("/login", Login)

	// The user handlers, requiring authentication
	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		//u.POST("/info", PostInfo())
		//u.GET("/info", GetInfo())
	}

	return g
}
