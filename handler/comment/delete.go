package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/gin-gonic/gin"
)

// 删除评课
// @Summary 删除评课
// @Tags comment
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
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 验证当前用户是否有删除此评课的权限
	if evaluation.UserId != userId {
		handler.SendForbidden(c, errno.ErrDelete, nil, "Without permission to delete the evaluation. ")
		return
	}

	if err = evaluation.Delete(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
}
