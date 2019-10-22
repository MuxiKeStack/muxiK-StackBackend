package comment

import (
	"errors"
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

// 删除评课
func Delete(c *gin.Context) {
	var err error
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	userId := c.MustGet("id").(uint32)

	evaluation, err := model.GetEvaluationById(uint32(id))
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	if evaluation.UserId != userId {
		err := errors.New("Permission denied ")
		handler.SendError(c, err, nil, err.Error())
	}
	err = evaluation.Delete()

	handler.SendResponse(c, nil, nil)
}
