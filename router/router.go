package router

import (
	"net/http"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler/comment"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/sd"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/table"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler/user"
	"github.com/MuxiKeStack/muxiK-StackBackend/router/middleware"

	"github.com/gin-gonic/gin"
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
	g.POST("/login", user.Login)

	// The user handlers, requiring authentication
	u := g.Group("/api/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("/info", user.PostInfo)
		u.GET("/info", user.GetInfo)
	}

	// 选课手册课程
	g.GET("/api/v1/course/using/home/", course.GetCourseInfo)
	g.GET("/api/v1/course/using/:id/query/", course.QueryCourse)
	course := g.Group("/api/v1/course/using")
	course.Use(middleware.AuthMiddleware())
	{
		course.PUT("/:id/add", course.AddCourse)
		course.POST("/:id/modify", course.ModifyCourse)
		course.DELETE("/:id/delete", course.DeleteCourse)
		course.POST("/:id/favorite/", course.FavoriteCourse)
	}
	
	// 云课堂课程
	g.GET("/api/v1/course/history/home/", course.GetCourseInfo)
	g.GET("/api/v1/course/history/:id/query/", course.QueryCourse)
	course := g.Group("/api/v1/course/history")
	course.Use(middleware.AuthMiddleware())
	{
		course.PUT("/:id/add", course.AddHistoryCourse)
		course.POST("/:id/modify", course.ModifyHistoryCourse)
		course.DELETE("/:id/delete", course.DeleteHistoryCourse)
	}
	
	// 评课
	g.GET("/api/v1/evaluation/list/", comment.EvaluationPlayground)
	g.GET("/api/v1/evaluation/:id/", comment.GetEvaluation)

	evaluation := g.Group("/api/v1/evaluation")
	evaluation.Use(middleware.AuthMiddleware())
	{
		evaluation.POST("/", comment.Publish)
		evaluation.POST("/:id/comment/", comment.CreateTopComment)
		evaluation.DELETE("/:id/", comment.Delete)
		evaluation.PUT("/:id/like/", comment.UpdateEvaluationLike)
	}

	// 评论
	g.GET("/api/v1/evaluation/:id/comments/", comment.GetComments)

	comments := g.Group("/api/v1/comment")
	comments.Use(middleware.AuthMiddleware())
	{
		comments.POST("/:id/", comment.Reply)
		comments.PUT("/:id/like/", comment.UpdateCommentLike)
	}

	// 排课课表
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
