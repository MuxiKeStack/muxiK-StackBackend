package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
)

//获取课程信息
func GetCourseInfo(c *gin.Context) {
	courseId, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	courseInfo, err := service.GetCourseInfo(courseid)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, courseInfo)
}
