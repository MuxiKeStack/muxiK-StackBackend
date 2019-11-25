package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

//获取课程信息
func GetCourseInfo(c *gin.Context) {
	courseId := c.DefaultQuery("courseId", "")
	if courseId == "" || len(courseId) != 8 {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, "CourseId must be required and length must be 8")
		return
	}

	// courseInfo, err := service.GetCourseInfo(courseid)
	// if err != nil {
	// 	handler.SendError(c, err, nil, err.Error())
	// 	return
	// }

	// handler.SendResponse(c, nil, courseInfo)
}
