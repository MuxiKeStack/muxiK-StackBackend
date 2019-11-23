package evaluation

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
	g            *gin.Engine
	tokenStr     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NzM1Mjc1OTEsImlkIjoxLCJuYmYiOjE1NzM1Mjc1OTF9.1CZFz2OVeDfDnvEXwCpQjqNGpSCIRoZOgMkRpuPIgc8"
	courseId     = "112d34testsvggase"
	courseName   = "高等数学A"
	evaluationId uint32
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

// Test: create a new evaluation
func TestPublish(t *testing.T) {
	g := getRouter(true)
	uri := "api/v1/evaluation/"
	body := evaluationPublishRequest{
		CourseId:            courseId,
		CourseName:          courseName,
		Rate:                7.5,
		AttendanceCheckType: 1,
		ExamCheckType:       2,
		Content:             "老师讲课很棒",
		IsAnonymous:         false,
		Tags:                []uint8{5, 2, 1},
	}

	jsonByte, err := json.Marshal(body)
	if err != nil {
		t.Errorf("Test Error: Json Marshal Error: %s", err.Error())
	}

	w := util.PerformRequestWithBody(http.MethodPost, g, uri, jsonByte, tokenStr)

	var data struct {
		Code    int
		Message string
		Data    *evaluationPublishResponse
	}

	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Publish New Evaluation Error; Json Unmarshal Error: %s", err.Error())
	}

	evaluationId = data.Data.EvaluationId
	fmt.Printf("--- evaluationId = %d\n", evaluationId)

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: get an evaluation information by a specific user
func TestGetEvaluation(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("api/v1/evaluation/%d/", evaluationId)
	w := util.PerformRequest(http.MethodGet, g, uri, tokenStr)

	var data struct {
		Code    int
		Message string
		Data    model.EvaluationInfo
	}
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get Evaluation Info Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: get an evaluation information by a visitor
func TestGetEvaluation2(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("api/v1/evaluation/%d/", evaluationId)
	w := util.PerformRequest(http.MethodGet, g, uri, "")

	var data struct {
		Code    int
		Message string
		Data    model.EvaluationInfo
	}
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get Evaluation Info Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: get evaluations at evaluation playground by a specific user
func TestEvaluationPlayground(t *testing.T) {
	g := getRouter(true)
	uri := "api/v1/evaluation/?pageSize=20&lastEvaluationId=0"
	w := util.PerformRequest(http.MethodGet, g, uri, tokenStr)

	var data struct {
		Code    int
		Message string
		Data    playgroundResponse
	}
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get Evaluation List Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: get evaluations at evaluation playground by a visitor
func TestEvaluationPlayground2(t *testing.T) {
	g := getRouter(true)
	uri := "api/v1/evaluation/"
	w := util.PerformRequest(http.MethodGet, g, uri, "")

	var data struct {
		Code    int
		Message string
		Data    playgroundResponse
	}

	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get Evaluation List Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

func TestEvaluationsOfOneCourse(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("api/v1/course/%s/evaluations/?limit=10&lastId=0&sort=hot", courseId)
	w := util.PerformRequest(http.MethodGet, g, uri, tokenStr)

	var data struct {
		Code    int
		Message string
		Data    evaluationsOfCourseResponse
	}

	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get Evaluations of One Course Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

func TestEvaluationsOfOneCourse2(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("api/v1/course/%s/evaluations/?limit=10&lastId=0&sort=time", courseId)
	w := util.PerformRequest(http.MethodGet, g, uri, "")

	var data struct {
		Code    int
		Message string
		Data    evaluationsOfCourseResponse
	}

	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Get Evaluations of One Course Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}

// Test: change an evaluation's like state by a user
func TestUpdateEvaluationLike(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("api/v1/evaluation/%d/like/", evaluationId)
	body := likeDataRequest{LikeState: true}

	jsonByte, err := json.Marshal(body)
	if err != nil {
		t.Errorf("Test Error: Json Marshal Error: %s", err.Error())
	}

	w := util.PerformRequestWithBody(http.MethodPut, g, uri, jsonByte, tokenStr)

	var data struct {
		Code    int
		Message string
		Data    likeDataResponse
	}
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test Error: Update Evaluation Like State Error; Json Unmarshal Error: %s", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error: %d", w.Code)
	}
}


// Test: Delete a evaluation
func TestDelete(t *testing.T) {
	g := getRouter(true)
	uri := fmt.Sprintf("api/v1/evaluation/%d/", evaluationId)
	w := util.PerformRequest(http.MethodDelete, g, uri, tokenStr)

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

	// Router for course evaluations

	evaluation := g.Group("/api/v1/evaluation")
	evaluation.Use(middleware.VisitorAuthMiddleware())
	{
		evaluation.GET("/", EvaluationPlayground)
		evaluation.GET("/:id/", GetEvaluation)
	}

	g.GET("/api/v1/course/:id/evaluations/", EvaluationsOfOneCourse).Use(middleware.VisitorAuthMiddleware())

	evaluationWithAuth := g.Group("/api/v1/evaluation")
	evaluationWithAuth.Use(middleware.AuthMiddleware())
	{
		evaluationWithAuth.POST("/", Publish)
		evaluationWithAuth.DELETE("/:id/", Delete)
		evaluationWithAuth.PUT("/:id/like/", UpdateEvaluationLike)
	}

	return g
}
