package tag

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

type tagGetResponse struct {
	Sum  int               `json:"sum"`
	List *[]model.TagModel `json:"list"`
}

// @Summary 获取课程评价标签列表
// @Tags tag
// @Success 200 {object} tag.tagGetResponse
// @Router /tags/ [get]
func Get(c *gin.Context) {
	tags := model.GetTags()

	result := tagGetResponse{
		Sum:  len(*tags),
		List: tags,
	}

	handler.SendResponse(c, nil, result)
}
