package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type selfCoursesResponse struct {
	Sum  int                           `json:"sum"`
	Data *[]service.ProducedCourseItem `json:"data"`
}

// @Summary 获取个人历史课程
// @Tags user
// @Param token header string true "token"
// @Param year query string true "学年，默认获取全部"
// @Param term query string true "学期，1/2/3，默认0表示获取全部"
// @Success 200 {object} model.ClassTableInfo
// @Router /user/courses/ [get]
func GetSelfCourses(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	var l model.LoginModel
	if err := c.ShouldBindJSON(&l); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	year := c.DefaultQuery("year", "0")
	term := c.DefaultQuery("term", "0")

	data, err := service.GetSelfCourseList(userId, l.Sid, l.Password, year, term)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, &selfCoursesResponse{
		Sum:  len(*data),
		Data: data,
	})
}
