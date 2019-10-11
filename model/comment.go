package model

import (
	"strconv"
	"time"
)

// 新增评课
func NewEvaluation(data *EvaluationPublish, userId uint64) (uint64, error) {
	var tagStr string
	for _, tagBox := range data.Tags {
		tagStr = strconv.FormatUint(tagBox.TagId, 10) + ","
	}

	newEvaluation := &CourseEvaluationModel{
		CourseId: data.CourseId,
		CourseName: data.CourseName,
		Rate: data.Rate,
		AttendanceCheckType: data.AttendanceCheckType,
		ExamCheckType: data.ExamCheckType,
		Content: data.Content,
		IsAnonymous: data.IsAnonymous,
		UserId: userId,
		Tags: tagStr,
		LikeNum: 0,
		CommentNum: 0,
		IsValid: true,
		Time: strconv.FormatInt(time.Now().Unix(), 10),
	}

	DB.Self.NewRecord(newEvaluation)
	return 0, nil
}