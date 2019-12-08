package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type evaluationsOfCourseResponse struct {
	HotSum     int                     `json:"hot_sum"`
	HotList    *[]model.EvaluationInfo `json:"hot_list"`
	NormalSum  int                     `json:"normal_sum"`
	NormalList *[]model.EvaluationInfo `json:"normal_list"`
}

// @Summary 课程所有评课和热评
// @Tags evaluation
// @Param token header string false "游客登录则不需要此字段或为空"
// @Param id path string true "课程id"
// @Param hot_limit query integer true "热评数"
// @Param limit query integer true "评课数"
// @Param last_id query integer true "上一次请求的最后一个评课的id，若是初始请求则为0"
// @Success 200 {object} evaluationsOfCourseResponse
// @Router /course/history/{id}/evaluations/ [get]
func EvaluationsOfOneCourse(c *gin.Context) {
	log.Info("EvaluationsOfOneCourse function is called.")

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "limit parse error.")
		return
	}

	hotLimitStr := c.DefaultQuery("hot_limit", "10")
	hotLimit, err := strconv.ParseInt(hotLimitStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "hot_limit parse error.")
		return
	}

	lastIdStr := c.DefaultQuery("last_id", "0")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
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

	// Get hot evaluations
	hotList, err := service.GetHotEvaluations(courseId, int32(hotLimit), userId, visitor)
	if err != nil {
		handler.SendError(c, errno.ErrGetHotEvaluations, nil, err.Error())
		return
	}

	// 获取评课列表
	list, err := service.GetEvaluationsOfOneCourse(int32(lastId), int32(limit), userId, visitor, courseId)
	if err != nil {
		handler.SendError(c, errno.ErrEvaluationList, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, evaluationsOfCourseResponse{
		HotSum:     len(*hotList),
		HotList:    hotList,
		NormalSum:  len(*list),
		NormalList: list,
	})
}
