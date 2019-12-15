package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type likeDataRequest struct {
	LikeState bool `json:"like_state"`
}

//收藏课程
func FavoriteCourse(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		log.Info("get hash error")
		return
	}

	userId := c.MustGet("id").(uint32)

	var course = &model.UsingCourseModel{Hash: hash}
	if err := course.GetByHash(); err != nil {
		log.Info("Course.GetByHash function error")
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	hasLiked := course.HasFavorited(userId)

	// 获取请求中当前的点赞状态
	var bodyData likeDataRequest
	if err := c.BindJSON(&bodyData); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 未点赞
	if bodyData.LikeState && !hasLiked {
		handler.SendResponse(c, errno.ErrNotLiked, nil)
		return
	}
	//	已点赞
	if !bodyData.LikeState && hasLiked {
		handler.SendResponse(c, errno.ErrHasLiked, nil)
		return
	}

	var err error

	// 点赞或者取消点赞
	if bodyData.LikeState {
		err = course.Unfavorite(userId)
	} else {
		err = course.Favorite(userId)
	}

	if err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
}
