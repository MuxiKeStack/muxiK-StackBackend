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

	tagStr := TagArrayToStr(data.Tags)

	newEvaluation := &CourseEvaluationModel{
		CourseId:            data.CourseId,
		CourseName:          data.CourseName,
		Rate:                data.Rate,
		AttendanceCheckType: data.AttendanceCheckType,
		ExamCheckType:       data.ExamCheckType,
		Content:             data.Content,
		IsAnonymous:         data.IsAnonymous,
		UserId:              userId,
		Tags:                tagStr,
		LikeNum:             0,
		CommentNum:          0,
		IsValid:             true,
		Time:                strconv.FormatInt(time.Now().Unix(), 10),
	}

	DB.Self.NewRecord(newEvaluation)

	id, _ := getLastInsertId()

	// 错误情况，未能获取id

	return id, nil
}

// 新增评论
func NewComment(data *NewCommentRequest, id uint64, isRoot bool, userId uint64) (uint64, error) {
	var newComment = &CommentModel{}
	var newCommentId uint64

	// 是否是评论评课
	if isRoot {
		newComment = &CommentModel{
			UserId:          userId,
			ParentId:        0,
			CommentTargetId: id,
			Content:         data.Content,
			LikeNum:         0,
			IsRoot:          true,
			Time:            strconv.FormatInt(time.Now().Unix(), 10),
			SubCommentNum:   0,
		}

		DB.Self.Create(newComment)
	} else if isTopComment(id) {
		// id是否是一级评论

		newComment = &CommentModel{
			UserId:          userId,
			ParentId:        id,
			CommentTargetId: id,
			Content:         data.Content,
			LikeNum:         0,
			IsRoot:          false,
			Time:            strconv.FormatInt(time.Now().Unix(), 10),
			SubCommentNum:   0,
		}

		DB.Self.Create(newComment)
	} else {
		// 根据id查找父评论

		parentId, err := getParentId(id)
		if err != nil {
			return 0, err
		}
		newComment = &CommentModel{
			UserId:          userId,
			ParentId:        parentId,
			CommentTargetId: id,
			Content:         data.Content,
			LikeNum:         0,
			IsRoot:          false,
			Time:            strconv.FormatInt(time.Now().Unix(), 10),
			SubCommentNum:   0,
		}
		DB.Self.Create(newComment)
	}

	newCommentId, _ = getLastInsertId()

	// 错误情况，未能获取id

	return newCommentId, nil
}

// 获取评论详情
func GetCommentInfo(id uint64, userId uint64) (*CommentInfo, error) {
	var data = &CommentInfo{}
	var c CommentModel

	d := DB.Self.First(&c, "id = ?", id)
	if d.Error != nil {
		return &CommentInfo{}, d.Error
	}

	commentUser, err := GetUserInfo(c.UserId)
	if err != nil {
		return &CommentInfo{}, nil
	}

	targetUser, err := GetUserInfo(c.CommentTargetId)
	if err != nil {
		return &CommentInfo{}, nil
	}

	data = &CommentInfo{
		Content:           c.Content,
		LikeNum:           c.LikeNum,
		IsLike:            GetCommentLike(id, userId),
		Time:              c.Time,
		UserInfo:          *commentUser,
		CommentTargetInfo: *targetUser,
	}

	return data, nil
}

// 获取评课详情
func GetEvaluationInfo(id uint64, userId uint64) (*EvaluationInfo, error) {
	var data = &EvaluationInfo{}
	var e CourseEvaluationModel

	d := DB.Self.Find(&e, "id = ?", id)
	if d.Error != nil {
		return &EvaluationInfo{}, d.Error
	}

	u, err := GetUserInfo(userId)
	if err != nil {
		return &EvaluationInfo{}, err
	}

	data = &EvaluationInfo{
		CourseId:            e.CourseId,
		CourseName:          e.CourseName,
		teacher:             "",
		Rate:                e.Rate,
		AttendanceCheckType: e.AttendanceCheckType,
		ExamCheckType:       e.ExamCheckType,
		Content:             e.Content,
		Time:                e.Time,
		IsAnonymous:         e.IsAnonymous,
		IsLike:              GetEvaluationLike(id, userId),
		LikeNum:             e.LikeNum,
		CommentNum:          e.CommentNum,
		Tags:                TagStrToArray(e.Tags),
		UserInfo:            *u,
	}

	return data, nil
}

// 修改评课点赞状态
func UpdateEvaluationLikeState(id uint64, userId uint64) error {
	return nil
}

// 修改评论点赞状态
func UpdateCommentLikeState(id uint64, userId uint64) error {
	return nil
}

// 获取最新插入数据的id
func getLastInsertId() (uint64, error) {
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

// 是否是一级评论
func isTopComment(id uint64) bool {
	var data CommentModel
	DB.Self.First(&data, "id = ?", id)

	// 还有数据不存在的情况

	return data.IsRoot
}

// Get parentId by commentTargetId
func getParentId(id uint64) (uint64, error) {
	var data = new(CommentModel)
	DB.Self.Where("id = ?", id).First(data)

	// 还有数据不存在的情况

	return data.ParentId, nil
}

// 获取评论点赞状态
func GetCommentLike(id uint64, userId uint64) bool {
	var data CommentLikeModel
	DB.Self.Where("user_id = ? AND comment_id = ?", userId, id).Find(&data)
	if data.Id != 0 {
		return true
	}
	return false
}

// 获取评课点赞状态
func GetEvaluationLike(id uint64, userId uint64) bool {
	var data EvaluationLikeModel
	DB.Self.Where("user_id = ? AND evaluation_id = ?", userId, id).Find(&data)
	if data.Id != 0 {
		return true
	}
	return false
}
