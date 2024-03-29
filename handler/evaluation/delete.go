package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"github.com/gin-gonic/gin"
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
		log.Info("evaluation.GetById function error.")
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// Has been deleted.
	if evaluation.DeletedAt != nil {
		handler.SendBadRequest(c, errno.ErrDelete, nil, "Evaluation has already been deleted.")
		return
	}

	// 验证当前用户是否有删除此评课的权限
	if evaluation.UserId != userId {
		handler.SendError(c, errno.ErrDelete, nil, "With no permission to delete the evaluation. ")
		return
	}

	// 删除评课
	if err := model.DeleteEvaluation(evaluation); err != nil {
		log.Error("DeleteEvaluation function error", err)
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
}
