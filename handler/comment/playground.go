package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
)

type playgroundResponse struct {
	Sum  int                     `json:"sum"`
	List *[]model.EvaluationInfo `json:"list"`
}

// 评课广场获取评课列表
func EvaluationPlayground(c *gin.Context) {
	pageSize := c.DefaultQuery("pageSize", "20")
	size, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	} else if size <= 0 {
		handler.SendBadRequest(c, err, nil, "PageSize error")
	}

	lastIdStr := c.DefaultQuery("lastEvaluationId", "-1")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var userId uint64
	visitor := false
	// 游客登录
	if t := c.Request.Header.Get("token"); len(t) == 0 {
		visitor = true
	} else {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
		}
		userId = c.MustGet("id").(uint64)
	}

	list, err := model.GetLatestEvaluationList(lastId, size, userId, visitor)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	data := playgroundResponse{
		Sum:  len(*list),
		List: list,
	}

	handler.SendResponse(c, nil, data)
}