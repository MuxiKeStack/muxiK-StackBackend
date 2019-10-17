package comment

import (
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
	userId := c.MustGet("id").(uint64)

	err = model.DeleteEvaluation(id, userId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	handler.SendResponse(c, nil, nil)
}