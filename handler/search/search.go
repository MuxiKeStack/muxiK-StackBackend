// create at 2019.11.24 by shiina orez
// search handlers
package search

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type searchResponse struct {
	Courses []service.SearchCourseInfo
	Length  int
	Page    uint64
}

// SearchCourse API means to search the courses by courseName, courseID and teacherName
// @Summary 搜索课程接口
// @Tags search
// @Param keyword query string true "关键字"
// @Param page query integer true "页码"
// @Param limit query integer true "每页最大数"
// @Success 200 {object} search.searchResponse
// @Router /search/course/ [get]
func SearchCourse(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseUint(limitStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
	}
	keyword := c.DefaultQuery("keyword", "")
	courseList, err := service.SearchCourses(keyword, page, limit)
	if err != nil {
		handler.SendError(c, errno.ErrSearchCourse, nil, err.Error())
	}
	response := searchResponse{
		Courses: courseList,
		Length:  len(courseList),
		Page:    page,
	}
	handler.SendResponse(c, nil, response)
}
