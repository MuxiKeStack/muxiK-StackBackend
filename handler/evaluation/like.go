package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/lexkong/log"

	"github.com/gin-gonic/gin"
)

type likeDataResponse struct {
	LikeState bool   `json:"like_state"`
	LikeNum   uint32 `json:"like_num"`
}

type likeDataRequest struct {
	LikeState bool `json:"like_state"`
}

// 评课点赞/取消点赞
// @Summary 评课点赞/取消点赞
// @Tags evaluation
// @Param token header string true "token"
// @Param id path string true "点赞的评课id"
// @Param data body evaluation.likeDataRequest true "当前点赞状态"
// @Success 200 {object} evaluation.likeDataResponse
// @Router /evaluation/{id}/like/ [put]
func UpdateEvaluationLike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	// 获取请求中当前的点赞状态
	var bodyData likeDataRequest
	if err := c.BindJSON(&bodyData); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	var evaluation = &model.CourseEvaluationModel{Id: uint32(id)}
	if err := evaluation.GetById(); err != nil {
		log.Info("evaluation.GetById function error")
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	hasLiked := evaluation.HasLiked(userId)

	// 判断点赞请求是否合理
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

	// 点赞或者取消点赞
	if bodyData.LikeState {
		err = evaluation.CancelLiking(userId)
	} else {
		err = evaluation.Like(userId)
	}

	if err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// 更新点赞数
	num := 1
	if bodyData.LikeState {
		num = -1
	}
	if err := evaluation.UpdateLikeNum(num); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, likeDataResponse{
		LikeState: !hasLiked,
		//LikeNum:   model.GetEvaluationLikeSum(uint32(id)),
		LikeNum: evaluation.LikeNum,
	})

	// New message reminder for liking an evaluation
	if !bodyData.LikeState {
		err = service.NewMessageForEvaluationLiking(userId, evaluation)
		if err != nil {
			log.Error("NewMessageForEvaluationLiking failed", err)
		}
	}
}
