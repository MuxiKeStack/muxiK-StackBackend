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
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	// 评课&评论
	comments := g.Group("/api/v1/course")
	comments.Use(middleware.AuthMiddleware())
	{
		comments.POST("/comment/", comment.Publish)
		comments.POST("/:courseId/comment/", comment.Create)
		comments.DELETE("/:courseId/comment/:courseCommentId/", comment.Delete)
		comments.GET("/:courseId/comment/:courseCommentId/", comment.GetCourseCommentInfo)
		comments.GET("/:courseId/comment/:courseCommentId/commentList/", comment.GetCommentList)
		comments.PUT("/:courseId/comment/:targetId/like/", comment.UpdateLike)
	}

	// 排课课表
	tables := g.Group("/api/v1/table")
	tables.Use(middleware.AuthMiddleware())
	{
		tables.GET("/", table.Get)
		tables.PUT("/:tableId/rename/", table.Rename)
		tables.POST("/", table.AddTable)
		tables.POST("/:tableId/class/:classId/", table.AddClass)
		tables.DELETE("/:tableId/", table.DeleteTable)
		tables.DELETE("/:tableId/class/:classId/", table.DeleteClass)
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
