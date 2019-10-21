package comment

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Publish_2(c *gin.Context) {
	var data model.EvaluationPublish
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("id").(uint32)

	var evaluation = &model.CourseEvaluationModel{
		CourseId:            data.CourseId,
		CourseName:          data.CourseName,
		UserId:              userId,
		Rate:                data.Rate,
		AttendanceCheckType: data.AttendanceCheckType,
		ExamCheckType:       data.ExamCheckType,
		Content:             data.Content,
		LikeNum:             0,
		CommentNum:          0,
		Tags:                model.TagArrayToStr(data.Tags),
		IsAnonymous:         data.IsAnonymous,
		IsValid:             true,
		Time:                strconv.FormatInt(time.Now().Unix(), 10),
	}

	if err := evaluation.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	// 更新数据库中课程的评分信息
	if err := model.UpdateCourseByEva(evaluation.Id, data.Rate); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, &responseData{EvaluationId: evaluation.Id})
}

func CreateTopComment_2(c *gin.Context) {
	var data model.NewCommentRequest
	if err := c.BindJSON(&data); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("id").(uint32)
	evaluationId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var comment = &model.CommentModel{
		UserId:          userId,
		ParentId:        0,
		CommentTargetId: uint32(evaluationId),
		Content:         data.Content,
		LikeNum:         0,
		IsRoot:          true,
		Time:            strconv.FormatInt(time.Now().Unix(), 10),
		SubCommentNum:   0,
	}

	if err := comment.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	commentInfo, err := comment.GetInfo(userId, false)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, commentInfo)
}

// 获取评论列表
func GetComments_2(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	pageSize := c.DefaultQuery("pageSize", "20")
	size, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	} else if size <= 0 {
		handler.SendBadRequest(c, err, nil, "PageSize error")
	}

	lastIdStr := c.DefaultQuery("lastCommentId", "-1")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var userId uint32
	visitor := false
	// 游客登录
	if t := c.Request.Header.Get("token"); len(t) == 0 {
		visitor = true
	} else {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
		}
		userId = c.MustGet("id").(uint32)
	}

	list, count, err := model.GetCommentList(uint32(id), int32(lastId), int32(size), userId, visitor)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	data := commentListResponse{
		ParentCommentNum:  count,
		ParentCommentList: list,
	}

	handler.SendResponse(c, nil, data)
}
