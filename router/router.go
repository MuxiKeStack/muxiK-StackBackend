package router

import (
	"net/http"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler/grade"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler/collection"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/report"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler/upload"

	_ "github.com/MuxiKeStack/muxiK-StackBackend/docs"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/comment"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/course"
	eva "github.com/MuxiKeStack/muxiK-StackBackend/handler/evaluation"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/message"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/sd"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/search"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/table"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/tag"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/user"
	"github.com/MuxiKeStack/muxiK-StackBackend/router/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
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
	g.POST("api/v1/login/", user.Login)

	// The user handlers, requiring authentication
	u := g.Group("/api/v1/user/")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("/info/", user.PostInfo)
		u.GET("/info/", user.GetInfo)
		u.POST("/licence/", user.JoinPro)

		u.POST("/courses/", course.GetSelfCourses)
		u.GET("/evaluations/", eva.GetHistoryEvaluations)
	}

	// Upload image to oss
	up := g.Group("api/v1/upload")
	up.Use(middleware.AuthMiddleware())
	{
		up.POST("/image/", upload.Image)
	}

	// The message handlers, required authentication
	m := g.Group("/api/v1/message")
	m.Use(middleware.AuthMiddleware())
	{
		m.GET("/", message.Get)
		m.GET("/count/", message.Count)
		m.POST("/readall/", message.ReadAll)
	}

	// 默认是转到docs文件下的json，不需要修改。
	// url := ginSwagger.URL("http://kstack.test.muxi-tech.xyz/swagger/doc.json") // The url pointing to API definition
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 选课手册课程
	// g.GET("/api/v1/course/using/home/", course.GetCourseInfo)
	// g.GET("/api/v1/course/using/:id/query/", course.QueryCourse)
	courses := g.Group("/api/v1/course/using")
	courses.Use(middleware.AuthMiddleware())
	{
		courses.PUT("/:id/", course.AddCourse)
		courses.POST("/:id/", course.ModifyCourse)
		courses.DELETE("/:id/", course.DeleteCourse)
		courses.PUT("/:id/favorite/", course.FavoriteCourse)
	}

	courses2 := g.Group("/api/v1/course/using/info")
	courses2.Use(middleware.VisitorAuthMiddleware())
	{
		courses2.GET("/:hash/", course.GetCourseInfo)
	}

	// Router for course evaluations

	evaluation := g.Group("/api/v1/evaluation")
	evaluation.Use(middleware.VisitorAuthMiddleware())
	{
		evaluation.GET("/", eva.EvaluationPlayground)
		evaluation.GET("/:id/", eva.GetEvaluation)

		// confirm reported, to block evaluation
		evaluation.GET("/:id/block", report.BlockEvaluation)

		// router for getting comment list
		evaluation.GET("/:id/comments/", comment.GetComments)
	}

	evaOfCourse := g.Group("/api/v1/course/history/:id/evaluations/")
	evaOfCourse.Use(middleware.VisitorAuthMiddleware())
	{
		evaOfCourse.GET("", eva.EvaluationsOfOneCourse)
	}

	evaluationWithAuth := g.Group("/api/v1/evaluation")
	evaluationWithAuth.Use(middleware.AuthMiddleware())
	{
		evaluationWithAuth.POST("/", eva.Publish)
		evaluationWithAuth.DELETE("/:id/", eva.Delete)
		evaluationWithAuth.PUT("/:id/like/", eva.UpdateEvaluationLike)

		evaluationWithAuth.POST("/:id/comment/", comment.CreateTopComment)

		// report a evaluation
		evaluation.POST("/:id/report/", report.ReportEvaluation)
	}

	// Router for comments requiring auth
	comments := g.Group("/api/v1/comment")
	comments.Use(middleware.AuthMiddleware())
	{
		comments.POST("/:id/", comment.Reply)
		comments.DELETE("/:id/", comment.Delete)
		comments.PUT("/:id/like/", comment.UpdateCommentLike)
	}

	// class table
	tables := g.Group("/api/v1/table")
	tables.Use(middleware.AuthMiddleware())
	{
		tables.GET("/", table.Get)
		tables.POST("/", table.AddTable)
		tables.POST("/:id/class/", table.AddClass)
		tables.PUT("/:id/rename/", table.Rename)
		tables.DELETE("/:id/", table.DeleteTable)
		tables.DELETE("/:id/class/", table.DeleteClass)
	}

	// tag
	g.GET("/api/v1/tags/", tag.Get)

	// grade
	grades := g.Group("/api/v1/grade")
	grades.Use(middleware.AuthMiddleware())
	{
		grades.GET("/", grade.Get)
	}

	// search
	searchGroup := g.Group("/api/v1/search")
	{
		searchGroup.GET("/course/", search.SearchCourse)
		searchGroup.GET("/historyCourse/", search.SearchHistoryCourse)
	}

	collections := g.Group("/api/v1/collection")
	collections.Use(middleware.AuthMiddleware())
	{
		collections.GET("/table/:id/", collection.CollectionsForTable)
		collections.GET("/", collection.GetCollections)
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
