package table

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/MuxiKeStack/muxiK-StackBackend/config"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/router/middleware"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/gin-gonic/gin"
)

var (
	g *gin.Engine
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NzM1Mjc1OTEsImlkIjoxLCJuYmYiOjE1NzM1Mjc1OTF9.1CZFz2OVeDfDnvEXwCpQjqNGpSCIRoZOgMkRpuPIgc8"
	tableId uint32
	classId = "sadf23432234dfa"
)

func TestMain(m *testing.M) {
	// init config
	if err := config.Init(""); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// Run the other tests
	os.Exit(m.Run())
}

// Test: create new blank table
func TestAddTable(t *testing.T) {
	g := getRouter(true)
	uri := "/api/v1/table/"
	w := util.PerformRequest(http.MethodPost, g, uri, token)

	var data model.ClassTableInfo
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Creat A New Table Error; Json Unmarshal Error: %s", err.Error())
	}

	tableId = data.TableId

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: add a class into a table
func TestAddClass(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("/api/v1/table/%s/class/?classId=%s", tableId, classId)
	w := util.PerformRequest(http.MethodPost, g, uri, token)

	var data addClassResponseData
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Add New Class Into Table Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: create new table copied by a existing table
func TestAddTable2(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("/api/v1/table/?id=%s", tableId)
	w := util.PerformRequest(http.MethodPost, g, uri, token)

	var data model.ClassTableInfo
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Creat A New Table Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: get table information
func TestGet(t *testing.T) {
	g := getRouter(true)
	uri := "/api/v1/table/"
	w := util.PerformRequest(http.MethodGet, g, uri, token)

	var data getTablesResponse
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get Tables Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: change table's name
func TestRename(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("/api/v1/table/%s/rename/", tableId)
	body := renameBodyData{NewName: "new table"}

	jsonByte, err := json.Marshal(body)
	if err != nil {
		t.Errorf("Test Error: Json Marshal Error: %s", err.Error())
	}

	w := util.PerformRequestWithBody(http.MethodPut, g, uri, jsonByte, token)

	var data addClassResponseData
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Rename Table Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: remove a class from a table
func TestDeleteClass(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("/api/v1/table/%s/class/?classId=%s", tableId, classId)
	w := util.PerformRequest(http.MethodDelete, g, uri, token)

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: delete a table
func TestDeleteTable(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("/api/v1/table/%s/", tableId)
	w := util.PerformRequest(http.MethodDelete, g, uri, token)

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Helper function to create a router during testing
func getRouter(withRouter bool) *gin.Engine {
	g = gin.New()
	if withRouter {
		loadRouters(
			g,

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

	tables := g.Group("/api/v1/table")
	tables.Use(middleware.AuthMiddleware())
	{
		tables.GET("/", Get)
		tables.POST("/", AddTable)
		tables.POST("/:id/class/", AddClass)
		tables.PUT("/:id/rename/", Rename)
		tables.DELETE("/:id/", DeleteTable)
		tables.DELETE("/:id/class/", DeleteClass)
	}

	return g
}
