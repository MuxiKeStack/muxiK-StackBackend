package comment

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// @Summary 删除评论
// @Tags evaluation
// @Param token header string true "token"
// @Param id path string true "评论id"
// @Success 200 "OK"
// @Router /comment/{id}/ [delete]
func Delete(c *gin.Context) {
	id := c.Param("id")

	userId := c.MustGet("id").(uint32)

	var err error
	if ok := model.IsParentComment(id); ok {
		err = service.DeleteParentComment(id, userId)
	} else {
		err = service.DeleteSubComment(id, userId)
	}

	if err != nil {
		handler.SendError(c, errno.ErrDeleteComment, nil, err.Error())
	}

	handler.SendResponse(c, nil, nil)
}
