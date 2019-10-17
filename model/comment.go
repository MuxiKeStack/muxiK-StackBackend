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
		return 0, errors.New("Table does not exit. ")
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

// 获取评课详情
func GetEvaluationInfo(id, userId uint64, visitor bool) (*EvaluationInfo, error) {
	var e CourseEvaluationModel
	d := DB.Self.Find(&e, "id = ?", id)
	if d.Error != nil {
		return &EvaluationInfo{}, d.Error
	}

	data, err := evaluationSQLDataToResponseInfo(&e, userId, visitor)
	if err != nil {
		return &EvaluationInfo{}, err
	}

	return data, nil
}

// 评课表数据转换为返回的评课信息数据
func evaluationSQLDataToResponseInfo(e *CourseEvaluationModel, userId uint64, visitor bool) (*EvaluationInfo, error) {
	var err error
	var u = &UserInfo{}
	if !e.IsAnonymous {
		u, err = GetUserInfo(e.UserId)
		if err != nil {
			return &EvaluationInfo{}, err
		}
	}

	// 获取教师名

	var data = &EvaluationInfo{
		CourseId:            e.CourseId,
		CourseName:          e.CourseName,
		Teacher:             "",
		Rate:                e.Rate,
		AttendanceCheckType: e.AttendanceCheckType,
		ExamCheckType:       e.ExamCheckType,
		Content:             e.Content,
		Time:                e.Time,
		IsAnonymous:         e.IsAnonymous,
		IsLike:              false,
		LikeNum:             e.LikeNum,
		CommentNum:          e.CommentNum,
		Tags:                TagStrToArray(e.Tags),
		UserInfo:            u,
	}
	if !visitor {
		data.IsLike = GetEvaluationLikeState(e.Id, userId)
	}

	return data, nil
}

// 获取最新的评课列表
func GetLatestEvaluationList(lastId, size int64, userId uint64, visitor bool) (*[]EvaluationInfo, error) {
	var result []EvaluationInfo
	var data []CourseEvaluationModel
	if lastId != -1 {
		DB.Self.Where("id < ?", lastId).Order("id desc").Find(&data).Limit(size)
	} else {
		DB.Self.Order("id desc").Find(&data).Limit(size)
	}

	// 待优化：并发
	for _, e := range data {
		d, err := evaluationSQLDataToResponseInfo(&e, userId, visitor)
		if err != nil {
			return nil, err
		}
		result = append(result, *d)
	}

	return &result, nil
}

// 获取评论详情
func GetCommentInfo(id, userId uint64) (*CommentInfo, error) {
	var data = &CommentInfo{}
	var c CommentModel

	d := DB.Self.First(&c, "id = ?", id)
	if d.Error != nil {
		return nil, d.Error
	}

	commentUser, err := GetUserInfo(c.UserId)
	if err != nil {
		return nil, nil
	}

	targetUser, err := GetUserInfo(c.CommentTargetId)
	if err != nil {
		return nil, nil
	}

	data = &CommentInfo{
		Content:           c.Content,
		LikeNum:           c.LikeNum,
		IsLike:            GetCommentLikeState(id, userId),
		Time:              c.Time,
		UserInfo:          commentUser,
		CommentTargetInfo: targetUser,
	}

	return data, nil
}

// 获取评论列表
func GetCommentList(id uint64, lastId, size int64, userId uint64, visitor bool) (*[]ParentCommentInfo, uint64, error) {
	var count uint64
	var list  []ParentCommentInfo
	var data []CommentModel

	if lastId != -1 {
		DB.Self.Where("is_root = ? AND comment_target_id = ?", true, id).
			Find(&data).Count(&count).Limit(size)
	} else {
		DB.Self.Where("id < ? AND is_root = ? AND comment_target_id = ?", lastId, true, id).
			Find(&data).Count(&count).Limit(size)
	}

	// 优化：并发
	for _, i := range data {
		var subComments []CommentInfo
		var comments []CommentModel
		DB.Self.Find(&comments, "parent_id = ?", i.Id)

		// 优化：并发
		for _, k := range comments {
			commentUser, err := GetUserInfo(k.UserId)
			if err != nil {
				return nil, 0, nil
			}

			targetUser, err := GetUserInfo(k.CommentTargetId)
			if err != nil {
				return nil, 0, nil
			}

			isLiked := false
			if !visitor {
				isLiked = GetCommentLikeState(i.Id, userId)
			}

			subComments = append(subComments, CommentInfo{
				Content:           k.Content,
				LikeNum:           k.LikeNum,
				IsLike:            isLiked,
				Time:              k.Time,
				UserInfo:          commentUser,
				CommentTargetInfo: targetUser,
			})
		}

		userInfo, err := GetUserInfo(i.UserId)
		if err != nil {
			return nil, 0, nil
		}

		isLiked := false
		if !visitor {
			isLiked = GetCommentLikeState(i.Id, userId)
		}

		list = append(list, ParentCommentInfo{
			CommentId:       i.Id,
			Content:         i.Content,
			LikeNum:         i.LikeNum,
			IsLike:          isLiked,
			Time:            i.Time,
			UserInfo:        userInfo,
			SubCommentsNum:  i.SubCommentNum,
			SubCommentsList: &subComments,
		})
	}

	return &list, count, nil
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

// 修改评课点赞状态
func UpdateEvaluationLikeState(id, userId uint64, like bool) error {
	d := &EvaluationLikeModel{
		EvaluationId: id,
		UserId:       userId,
	}
	var count uint64
	DB.Self.Find(d).Count(&count)

	// 点赞
	if !like {
		// 检查是否含有记录
		if count != 0 {
			return errors.New("Have liked already. ")
		}
		// 添加点赞记录
		DB.Self.Create(d)
		updateEvaluationLikeNum(id, 1)
	} else {
		if count == 0 {
			return errors.New("Have not liked. ")
		}
		// 删除记录
		DB.Self.Delete(d)
		updateEvaluationLikeNum(id, -1)
	}

	return nil
}

// 修改评论点赞状态
func UpdateCommentLikeState(id, userId uint64, like bool) error {
	d := &CommentLikeModel{
		CommentId: id,
		UserId:    userId,
	}
	var count uint64
	DB.Self.Find(d).Count(&count)

	// 点赞
	if !like {
		// 检查是否含有记录
		if count != 0 {
			return errors.New("Have liked already. ")
		}
		// 添加点赞记录
		DB.Self.Create(d)
		updateCommentLikeNum(id, 1)
	} else {
		if count == 0 {
			return errors.New("Have not liked. ")
		}
		// 删除记录
		DB.Self.Delete(d)
		updateCommentLikeNum(id, -1)
	}

	return nil
}

// 获取评论点赞状态
func GetCommentLikeState(id, userId uint64) bool {
	var data CommentLikeModel
	DB.Self.Where("user_id = ? AND comment_id = ?", userId, id).Find(&data)
	if data.Id != 0 {
		return true
	}
	return false
}

// 获取评课点赞状态
func GetEvaluationLikeState(id, userId uint64) bool {
	var data EvaluationLikeModel
	DB.Self.Where("user_id = ? AND evaluation_id = ?", userId, id).Find(&data)
	if data.Id != 0 {
		return true
	}
	return false
}

// 获取评课点赞数
func GetEvaluationLikeNum(id uint64) uint64 {
	var count uint64
	DB.Self.Where("id = ?", id).Find(&CourseEvaluationModel{}).Count(&count)

	return count
}

// 获取评论点赞数
func GetCommentLikeNum(id uint64) uint64 {
	var count uint64
	DB.Self.Where("id = ?", id).Find(&CommentModel{}).Count(&count)

	return count
}

// 修改评课点赞数
func updateEvaluationLikeNum(id uint64, num int8) uint64 {
	var d = &CourseEvaluationModel{}
	DB.Self.Where("id = ?", id).Find(d)

	if num < 0 {
		if d.LikeNum < uint64(-num) {
			d.LikeNum = 0
		} else {
			d.LikeNum -= uint64(-num)
		}
	} else {
		d.LikeNum += uint64(num)
	}
	DB.Self.Save(d)

	return d.LikeNum
}

// 修改评论数
func updateCommentLikeNum(id uint64, num int8) uint64 {
	var d = &CommentModel{}
	DB.Self.Where("id = ?", id).Find(d)

	if num < 0 {
		if d.LikeNum < uint64(-num) {
			d.LikeNum = 0
		} else {
			d.LikeNum -= uint64(-num)
		}
	} else {
		d.LikeNum += uint64(num)
	}
	DB.Self.Save(d)

	return d.LikeNum
}

// 删除评课
func DeleteEvaluation(id, userId uint64) error {
	var e CourseEvaluationModel
	DB.Self.Find(&e, "id = ?", id)

	if e.UserId != userId {
		return errors.New("Permission denied ")
	}
	DB.Self.Delete(&e)
	return nil
}