package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"github.com/gin-gonic/gin"
)

type likeDataRequest struct {
	LikeState bool `json:"like_state"`
}

// 收藏课程
// @Summary 收藏课程/取消收藏
// @Tags course
// @Param token header string true "token"
// @Param id path string true "收藏的课程id"
// @Param data body course.likeDataRequest true "当前收藏状态"
// @Success 200 {object} string	"ok"
// @Router /course/using/{id}/favorite/ [put]
func FavoriteCourse(c *gin.Context) {
	log.Info("FavoriteCourse function is called.")
	hash := c.Param("id")
	if hash == "" {
		log.Info("get hash error")
		return
	}

	userId := c.MustGet("id").(uint32)

	recordId, hasLiked := model.HasFavorited(userId, hash)

	// 获取请求中当前的收藏状态
	var bodyData likeDataRequest
	if err := c.BindJSON(&bodyData); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 未收藏
	if bodyData.LikeState && !hasLiked {
		handler.SendResponse(c, errno.ErrNotLiked, nil)
		return
	}
	// 已收藏
	if !bodyData.LikeState && hasLiked {
		handler.SendResponse(c, errno.ErrHasLiked, nil)
		return
	}

	var err error

	// 收藏或者取消收藏
	if bodyData.LikeState {
		err = model.Unfavorite(recordId)
	} else {
		err = model.Favorite(userId, hash)
	}

	if err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	log.Info("success")
	handler.SendResponse(c, nil, "ok")
}
