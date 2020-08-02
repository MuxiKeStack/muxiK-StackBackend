package collection

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type collectionListResponse struct {
	Sum  int                               `json:"sum"`
	List []*model.CourseInfoForCollections `json:"list"`
}

// @Summary 获取课程清单列表
// @Tags collection
// @Param token header string true "token"
// @Param limit query integer true "期望的数量"
// @Param last_id query string true "上一次请求的最后一个记录的id，若是初始请求则为0"
// @Success 200 {object} collection.collectionListResponse
// @Router /collection/ [get]
func GetCollections(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	lastIdStr := c.DefaultQuery("last_id", "0")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	data, err := service.GetCollectionList(userId, int32(lastId), int32(limit))
	if err != nil {
		log.Error("GetCollectionList function error", err)
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, collectionListResponse{
		Sum:  len(data),
		List: data,
	})
}
