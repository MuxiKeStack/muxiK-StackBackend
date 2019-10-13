package model

import (
	"errors"
	"strconv"
	"time"
)

// 新增评课
func NewEvaluation(data *EvaluationPublish, userId uint64) (uint64, error) {
	if exit := DB.Self.HasTable(&CourseEvaluationModel{}); !exit {
		//DB.Self.CreateTable(&CourseEvaluationModel{})
		return 0, errors.New("Table does not exit.")
	}

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

	DB.Self.Create(newEvaluation)
	rows, err := DB.Self.Raw("select last_insert_id()").Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id uint64
	if err := DB.Self.ScanRows(rows, id); err != nil {
		return 0, err
	}

	return id, nil
}
