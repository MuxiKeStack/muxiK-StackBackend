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
	tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NzM4MTc4NTksImlkIjoxLCJuYmYiOjE1NzM4MTc4NTl9.gfdq_WGp10Pxk3iGqRDascQ1wcSHaF37kMK3PCvYBlg"
	password    = "Yu@14796825550"
	sid         = "2018212576"
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
	u := g.Group("api/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("/info", PostInfo)
		u.GET("/info", GetInfo)
	}

	return g
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

// TestLogin function to test login router.
func TestLogin(t *testing.T) {
	g := getRouter(true)

	uri := "/login"
	u := CreateLoginRequest{
		model.LoginModel{
			Sid:      sid,
			Password: password,
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
	fmt.Println(tokenString)
	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", w.Code)
	}
}

// PostGetInfo function test post user information router.
func TestPostInfo(t *testing.T) {
	g := getRouter(true)
	uri := "api/v1/user/info"
	info := CreatePostInfoRequest{model.UserInfoRequest{
		Username: "Bowser",
		Avatar:   "https://www.gravatar.com/avatar/2af44a6505d5fa19f843ef83f1c61915?s=128&d=identicon",
	}}
	jsonByte, err := json.Marshal(info)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	w := util.PerformRequestWithBody(http.MethodPost, g, uri, jsonByte, tokenString)
	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", w.Code)
	}
}

// TestGetInfo function test get user information router.
func TestGetInfo(t *testing.T) {
	g := getRouter(true)
	uri := "api/v1/user/info"
	w := util.PerformRequest(http.MethodGet, g, uri, tokenString)
	// 读取响应body
	var data InfoResponse

	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get UserInfo Error:%s", err.Error())
	}
	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", w.Code)
	}

}
