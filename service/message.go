package service

import (
	"encoding/json"
	"fmt"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/lexkong/log"
	"strconv"
	"time"
)

func MessageList(page, limit, uid uint32) (*[]model.MessageSub, error) {
	messages, err := model.GetMessages(page, limit, uid)
	if err != nil {
		return nil, nil
	}
	var messageSubs []model.MessageSub
	for _, message := range *messages {
		messageSub := model.MessageSub{
			IsLike: message.IsLike,
			IsRead: message.IsRead,
			Reply:  message.Reply,
			Time:   message.Time,
		}
		userInfo, err := GetUserInfoRById(uid)
		if err != nil {
			return nil, err
		}

		var courseInfo model.CourseInfo
		err = json.Unmarshal([]byte(message.CourseInfo), &courseInfo)
		if err != nil {
			return nil, err
		}
		messageSub.CourseInfo = courseInfo
		messageSub.UserInfo = *userInfo
		messageSubs = append(messageSubs, messageSub)
	}
	return &messageSubs, nil
}

func NewMessageForParentComment(userId uint32, comment *model.ParentCommentModel, evaluation *model.CourseEvaluationModel) error {
	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		log.Info("GetTeacherByCourseId function error")
		return err
	}

	message := &model.MessagePub{
		PubUserId: userId,
		SubUserId: evaluation.UserId,
		IsLike:    false,
		IsRead:    false,
		Reply:     comment.Content,
		Time:      strconv.FormatInt(comment.Time.Unix(), 10),
		CourseInfo: model.CourseInfo{
			EvaluationId:    evaluation.Id,
			Sid:             "",
			ParentCommentId: comment.Id,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			Content:         evaluation.Content,
		},
	}

	err = model.CreateMessage(message)
	if err != nil {
		log.Info("CreateMessage function error")
		return err
	}
	return nil
}

func NewMessageForSubComment(userId uint32, sid string, comment *model.SubCommentModel, parentComment *model.ParentCommentModel) error {
	evaluation := &model.CourseEvaluationModel{Id: parentComment.EvaluationId}
	if err := evaluation.GetById(); err != nil {
		fmt.Println(parentComment.EvaluationId)
		log.Info("evaluation.GetById function error")
		return err
	}

	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		log.Info("GetTeacherByCourseId function error")
		return err
	}

	message := &model.MessagePub{
		PubUserId: userId,
		SubUserId: parentComment.UserId,
		IsLike:    false,
		IsRead:    false,
		Reply:     comment.Content,
		Time:      strconv.FormatInt(comment.Time.Unix(), 10),
		CourseInfo: model.CourseInfo{
			EvaluationId:    parentComment.EvaluationId,
			Sid:             sid,
			ParentCommentId: comment.Id,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			Content:         parentComment.Content,
		},
	}

	err = model.CreateMessage(message)
	if err != nil {
		log.Info("CreateMessage function error")
		return err
	}
	return nil
}

func NewMessageForEvaluationLiking(userId uint32, evaluation *model.CourseEvaluationModel) error {
	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		fmt.Println(evaluation.CourseId, teacher)
		log.Info("GetTeacherByCourseId function error")
		return err
	}

	message := &model.MessagePub{
		PubUserId: userId,
		SubUserId: evaluation.UserId,
		IsLike:    true,
		IsRead:    false,
		Reply:     "",
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
		CourseInfo: model.CourseInfo{
			EvaluationId:    evaluation.Id,
			Sid:             "",
			ParentCommentId: "",
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			Content:         evaluation.Content,
		},
	}

	err = model.CreateMessage(message)
	if err != nil {
		log.Info("CreateMessage function error")
		return err
	}
	return nil
}

func NewMessageForCommentLiking(userId uint32, commentId string) error {
	comment, ok := model.IsSubComment(commentId)
	if ok {
		return NewMessageForSubCommentLiking(userId, comment)
	}
	return NewMessageForParentCommentLiking(userId, commentId)
}

func NewMessageForParentCommentLiking(userId uint32, commentId string) error {
	comment := &model.ParentCommentModel{Id: commentId}
	if err := comment.GetById(); err != nil {
		log.Info("comment.GetById function error")
		return err
	}

	evaluation := &model.CourseEvaluationModel{Id: comment.EvaluationId}
	if err := evaluation.GetById(); err != nil {
		log.Info("evaluation.GetById function error")
		return err
	}

	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		log.Info("GetTeacherByCourseId function error")
		return err
	}

	message := &model.MessagePub{
		PubUserId: userId,
		SubUserId: comment.UserId,
		IsLike:    true,
		IsRead:    false,
		Reply:     "",
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
		CourseInfo: model.CourseInfo{
			EvaluationId:    evaluation.Id,
			Sid:             "",
			ParentCommentId: comment.Id,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			Content:         comment.Content,
		},
	}

	err = model.CreateMessage(message)
	if err != nil {
		log.Info("CreateMessage function error")
		return err
	}
	return nil
}

func NewMessageForSubCommentLiking(userId uint32, comment *model.SubCommentModel) error {
	parentComment := &model.ParentCommentModel{Id: comment.Id}
	if err := parentComment.GetById(); err != nil {
		log.Info("parentComment.GetById function error")
		return err
	}

	evaluation := &model.CourseEvaluationModel{Id: parentComment.EvaluationId}
	if err := evaluation.GetById(); err != nil {
		log.Info("evaluation.GetById function error")
		return err
	}

	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		log.Info("GetTeacherByCourseId function error")
		return err
	}

	message := &model.MessagePub{
		PubUserId: userId,
		SubUserId: comment.UserId,
		IsLike:    true,
		IsRead:    false,
		Reply:     "",
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
		CourseInfo: model.CourseInfo{
			EvaluationId:    evaluation.Id,
			Sid:             "",
			ParentCommentId: comment.Id,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			Content:         comment.Content,
		},
	}

	err = model.CreateMessage(message)
	if err != nil {
		log.Info("CreateMessage function error")
		return err
	}
	return nil
}
