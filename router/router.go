package router

import (
	"net/http"

	_ "github.com/MuxiKeStack/muxiK-StackBackend/docs"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/comment"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/message"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/sd"
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
	g.POST("api/v1/login", user.Login)

	// The user handlers, requiring authentication
	u := g.Group("/api/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("/info", user.PostInfo)
		u.GET("/info", user.GetInfo)
	}

	// The message handlers, required authentication
	m := g.Group("/api/v1/message")
	m.Use(middleware.AuthMiddleware())
	{
		m.GET("/", message.Get)
		m.GET("/count", message.Count)
		m.POST("/readall", message.ReadAll)
	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 选课手册课程
	// g.GET("/api/v1/course/using/home/", course.GetCourseInfo)
	// g.GET("/api/v1/course/using/:id/query/", course.QueryCourse)
	// course := g.Group("/api/v1/course/using")
	// course.Use(middleware.AuthMiddleware())
	// {
	// 	course.PUT("/:id/add", course.AddCourse)
	// 	course.POST("/:id/modify", course.ModifyCourse)
	// 	course.DELETE("/:id/delete", course.DeleteCourse)
	// 	course.POST("/:id/favorite/", course.FavoriteCourse)
	// }

	// // 云课堂课程
	// g.GET("/api/v1/course/history/home/", course.GetCourseInfo)
	// g.GET("/api/v1/course/history/:id/query/", course.QueryCourse)
	// course := g.Group("/api/v1/course/history")
	// course.Use(middleware.AuthMiddleware())
	// {
	// 	course.PUT("/:id/add", course.AddHistoryCourse)
	// 	course.POST("/:id/modify", course.ModifyHistoryCourse)
	// 	course.DELETE("/:id/delete", course.DeleteHistoryCourse)
	// }

	// Router for course evaluations

	evaluation := g.Group("/api/v1/evaluation")
	evaluation.Use(middleware.VisitorAuthMiddleware())
	{
		evaluation.GET("/", comment.EvaluationPlayground)
		evaluation.GET("/:id/", comment.GetEvaluation)

		// router for getting comment list
		evaluation.GET("/:id/comments/", comment.GetComments)
	}

	evaluationWithAuth := g.Group("/api/v1/evaluation")
	evaluationWithAuth.Use(middleware.AuthMiddleware())
	{
		evaluationWithAuth.POST("/", comment.Publish)
		evaluationWithAuth.DELETE("/:id/", comment.Delete)
		evaluationWithAuth.PUT("/:id/like/", comment.UpdateEvaluationLike)
		evaluationWithAuth.POST("/:id/comment/", comment.CreateTopComment)
	}

	// Router for comments
	comments := g.Group("/api/v1/comment")
	comments.Use(middleware.AuthMiddleware())
	{
		comments.POST("/:id/", comment.Reply)
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
