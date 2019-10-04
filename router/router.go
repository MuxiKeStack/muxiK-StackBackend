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
	comments := g.Group("/api/v1/comment")
	comments.Use(middleware.AuthMiddleware())
	{
		comments.POST("/", comment.Publish)
		comments.POST("/:id/new/", comment.Create)
		comments.DELETE("/:id/", comment.Delete)
		comments.GET("/:id/", comment.GetCourseCommentInfo)
		comments.GET("/:id/commentList/", comment.GetCommentList)
		comments.PUT("/:id/like/", comment.UpdateLike)
	}

	// 排课课表
	tables := g.Group("/api/v1/table")
	tables.Use(middleware.AuthMiddleware())
	{
		tables.GET("/", table.Get)
		tables.POST("/", table.AddTable)
		tables.POST("/:id/class/:classId/", table.AddClass)
		tables.PUT("/:id/rename/", table.Rename)
		tables.DELETE("/:id/", table.DeleteTable)
		tables.DELETE("/:id/class/:classId/", table.DeleteClass)
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
