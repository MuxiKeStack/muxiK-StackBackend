package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// 删除评课
// @Summary 删除评课
// @Tags evaluation
// @Param token header string true "token"
// @Param id path string true "评课id"
// @Success 200 "OK"
// @Router /evaluation/{id}/ [delete]
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	// Get evaluation by id
	evaluation := &model.CourseEvaluationModel{Id: uint32(id)}
	if err := evaluation.GetById(); err != nil {
		log.Infof("evaluation.GetById() error.")
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 验证当前用户是否有删除此评课的权限
	if evaluation.UserId != userId {
		handler.SendForbidden(c, errno.ErrDelete, nil, "With no permission to delete the evaluation. ")
		return
	}

	if err = evaluation.Delete(); err != nil {
		log.Info("evaluation.Delete() error.")
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 更新课程信息
	// ...

	handler.SendResponse(c, nil, nil)
}
