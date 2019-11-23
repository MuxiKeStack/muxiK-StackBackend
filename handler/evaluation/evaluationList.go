package evaluation

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"strconv"
)

type evaluationsOfCourseResponse struct {
	Sum int `json:"sum"`
	List *[]model.EvaluationInfo `json:"list"`
}

// @Summary 课程所有评课
// @Tags comment
// @Param token header string false "游客登录则不需要此字段或为空"
// @Param id path string true "课程id"
// @Param limit query integer true "评课数"
// @Param lastId query integer true "上一次请求的最后一个评课的id，若是初始请求则为0"
// @Param sort query string true "排序关键词：hot/time"
// @Success 200 {object} evaluationsOfCourseResponse
// @Router /course/{id}/evaluations/ [get]
func EvaluationsOfOneCourse(c *gin.Context) {
	size := c.DefaultQuery("limit", "20")
	limit, err := strconv.ParseInt(size, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	lastIdStr := c.DefaultQuery("lastId", "0")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	sortKey := c.DefaultQuery("sort", "hot")
	if sortKey != "hot" && sortKey != "time" {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "Error sort method key.")
		return
	}

	courseId := c.Param("id")

	// userId获取与游客模式判断
	var userId uint32
	visitor := false

	userIdInterface, ok := c.Get("id")
	if !ok {
		visitor = true
	} else {
		userId = userIdInterface.(uint32)
		log.Info("User auth successful.")
	}

	// 获取评课列表
	list, err := service.GetEvaluationsOfOneCourse(int32(lastId), int32(limit), userId, visitor, courseId, sortKey)
	if err != nil {
		handler.SendError(c, errno.ErrEvaluationList, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, evaluationsOfCourseResponse{
		Sum:  len(*list),
		List: list,
	})
}
