package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

/*
对消息提醒的整理，分为三种，评论，点赞，举报，用kind来表识。
所有消息提醒都有一个tag，即为MessageInfo，信息有课程名，课程老师名，课程ID
所有的点击都是返回评课id，所以必须要有的id是前两个。
课程ID courseID
评课ID envaluationID

下面两个可以不要
一级评论ID parentCommentID
二级评论ID subCommentID

评论分为两种：
	对评课的评论为一级评论 		 返回 	一级评论内容/reply		评课ID 	评课内容
	对一级评论的评论即二级评论   返回 	二级评论内容/reply		一级评论ID 一级评论内容
点赞分为
	对评课的点赞 			  返回 							评课ID 	评课内容
	对评论的点赞			  返回        评论内容
举报分为
	对评课的举报
*/
func MessageList(page, limit, uid uint32) (*[]model.MessageSub, error) {
	messages, err := model.GetMessages(page, limit, uid)
	if err != nil {
		return nil, nil
	}
	var messageSubs []model.MessageSub
	for _, message := range *messages {
		messageSub := model.MessageSub{
			Kind:   message.Kind,
			IsRead: message.IsRead,
			Reply:  message.Reply,
			Time:   message.Time,
		}
		userInfo, err := GetUserInfoRById(uid)
		if err != nil {
			return nil, err
		}

		var MessageInfo model.MessageInfo
		err = json.Unmarshal([]byte(message.MessageInfo), &MessageInfo)
		if err != nil {
			return nil, err
		}
		messageSub.MessageInfo = MessageInfo
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
		Kind:      1,
		IsRead:    false,
		Reply:     comment.Content,
		Time:      strconv.FormatInt(comment.Time.Unix(), 10),
		MessageInfo: model.MessageInfo{
			CourseId:        evaluation.CourseId,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			EvaluationId:    evaluation.Id,
			Content:         evaluation.Content,
			Sid:             GetSidById(evaluation.UserId),
			ParentCommentId: comment.Id,
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
		Kind:      0,
		IsRead:    false,
		Reply:     comment.Content,
		Time:      strconv.FormatInt(comment.Time.Unix(), 10),
		MessageInfo: model.MessageInfo{
			CourseId:        evaluation.CourseId,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			EvaluationId:    parentComment.EvaluationId,
			Content:         parentComment.Content,
			Sid:             sid,
			ParentCommentId: comment.Id,
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
		Kind:      0,
		IsRead:    false,
		Reply:     "",
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
		MessageInfo: model.MessageInfo{
			CourseId:        evaluation.CourseId,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			EvaluationId:    evaluation.Id,
			Content:         evaluation.Content,
			Sid:             "",
			ParentCommentId: "",
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
		Kind:      0,
		IsRead:    false,
		Reply:     "",
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
		MessageInfo: model.MessageInfo{
			CourseId:        evaluation.CourseId,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			EvaluationId:    evaluation.Id,
			Content:         comment.Content,
			Sid:             "",
			ParentCommentId: "",
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
		Kind:      0,
		IsRead:    false,
		Reply:     "",
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
		MessageInfo: model.MessageInfo{
			CourseId:        evaluation.CourseId,
			CourseName:      evaluation.CourseName,
			Teacher:         teacher,
			EvaluationId:    evaluation.Id,
			Content:         comment.Content,
			Sid:             "",
			ParentCommentId: "",
		},
	}

	err = model.CreateMessage(message)
	if err != nil {
		log.Info("CreateMessage function error")
		return err
	}
	return nil
}

func NewMessageForReport(userId uint32) {}
