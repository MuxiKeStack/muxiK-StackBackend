package model

import (
	"errors"
	"strconv"
	"time"
)

/*--------------- Course Evaluation Operation -------------*/

// 新增评课
func NewEvaluation(data *EvaluationPublish, userId uint32) (uint32, error) {
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

// 获取评课详情
func GetEvaluationInfo(id, userId uint32, visitor bool) (*EvaluationInfo, error) {
	var evaluation CourseEvaluationModel
	DB.Self.Find(&evaluation, "id = ?", id)

	data, err := evaluationSQLDataToResponseInfo(&evaluation, userId, visitor)
	if err != nil {
		return &EvaluationInfo{}, err
	}

	return data, nil
}

// 评课表数据转换为返回的评课信息数据
func evaluationSQLDataToResponseInfo(e *CourseEvaluationModel, userId uint32, visitor bool) (*EvaluationInfo, error) {
	var err error
	var u = &UserInfo{}
	if !e.IsAnonymous {
		u, err = GetUserInfoById(e.UserId)
		if err != nil {
			return &EvaluationInfo{}, err
		}
	}

	// 获取教师名
	course := &UsingCourseModel{}
	DB.Self.Find(course, "hash = ?", e.CourseId)

	var data = &EvaluationInfo{
		CourseId:            e.CourseId,
		CourseName:          e.CourseName,
		Teacher:             course.Teacher,
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
func GetLatestEvaluationList(lastId, size int32, userId uint32, visitor bool) (*[]EvaluationInfo, error) {
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

// 删除评课
func DeleteEvaluation(id, userId uint32) error {
	var e CourseEvaluationModel
	DB.Self.Find(&e, "id = ?", id)

	if e.UserId != userId {
		return errors.New("Permission denied ")
	}
	DB.Self.Delete(&e)
	return nil
}

// 修改评课点赞状态
func UpdateEvaluationLikeState(id, userId uint32, like bool) error {
	d := &EvaluationLikeModel{
		EvaluationId: id,
		UserId:       userId,
	}
	var count uint32
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

// 获取评课点赞状态
func GetEvaluationLikeState(id, userId uint32) bool {
	var data EvaluationLikeModel
	DB.Self.Where("user_id = ? AND evaluation_id = ?", userId, id).Find(&data)
	if data.Id != 0 {
		return true
	}
	return false
}

// 获取评课点赞数
func GetEvaluationLikeNum(id uint32) uint32 {
	var count uint32
	DB.Self.Where("id = ?", id).Find(&CourseEvaluationModel{}).Count(&count)

	return count
}

// 修改评课点赞数
func updateEvaluationLikeNum(id uint32, num int8) uint32 {
	var d = &CourseEvaluationModel{}
	DB.Self.Where("id = ?", id).Find(d)

	if num < 0 {
		if d.LikeNum < uint32(-num) {
			d.LikeNum = 0
		} else {
			d.LikeNum -= uint32(-num)
		}
	} else {
		d.LikeNum += uint32(num)
	}
	DB.Self.Save(d)

	return d.LikeNum
}

/*--------------- Comment Operation -------------*/

// 新增评论
func NewComment(data *NewCommentRequest, id uint32, isRoot bool, userId uint32) (uint32, error) {
	var newComment = &CommentModel{}
	var newCommentId uint32

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
func GetCommentInfo(id, userId uint32) (*CommentInfo, error) {
	var data = &CommentInfo{}
	var c CommentModel

	d := DB.Self.First(&c, "id = ?", id)
	if d.Error != nil {
		return nil, d.Error
	}

	commentUser, err := GetUserInfoById(c.UserId)
	if err != nil {
		return nil, nil
	}

	targetUser, err := GetUserInfoById(c.CommentTargetId)
	if err != nil {
		return nil, nil
	}

	data = &CommentInfo{
		Content:        c.Content,
		LikeNum:        c.LikeNum,
		IsLike:         GetCommentLikeState(id, userId),
		Time:           c.Time,
		UserInfo:       commentUser,
		TargetUserInfo: targetUser,
	}

	return data, nil
}

// 获取评论列表
func GetCommentList(id uint32, lastId, size int32, userId uint32, visitor bool) (*[]ParentCommentInfo, uint32, error) {
	var count uint32
	var list []ParentCommentInfo
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
			commentUser, err := GetUserInfoById(k.UserId)
			if err != nil {
				return nil, 0, err
			}

			targetUser, err := GetUserInfoById(k.CommentTargetId)
			if err != nil {
				return nil, 0, err
			}

			isLiked := false
			if !visitor {
				isLiked = GetCommentLikeState(i.Id, userId)
			}

			subComments = append(subComments, CommentInfo{
				Content:        k.Content,
				LikeNum:        k.LikeNum,
				IsLike:         isLiked,
				Time:           k.Time,
				UserInfo:       commentUser,
				TargetUserInfo: targetUser,
			})
		}

		userInfo, err := GetUserInfoById(i.UserId)
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

// 修改评论点赞状态
func UpdateCommentLikeState(id, userId uint32, like bool) error {
	d := &CommentLikeModel{
		CommentId: id,
		UserId:    userId,
	}
	var count uint32
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
func GetCommentLikeState(id, userId uint32) bool {
	var data CommentLikeModel
	var count int
	DB.Self.Where("user_id = ? AND comment_id = ?", userId, id).Find(&data).Count(&count)

	return count > 0
}

// 获取评论点赞数
func GetCommentLikeNum(id uint32) uint32 {
	var count uint32
	DB.Self.Where("id = ?", id).Find(&CommentModel{}).Count(&count)

	return count
}

// 修改评论数
func updateCommentLikeNum(id uint32, num int8) uint32 {
	var d = &CommentModel{}
	DB.Self.Where("id = ?", id).Find(d)

	if num < 0 {
		if d.LikeNum < uint32(-num) {
			d.LikeNum = 0
		} else {
			d.LikeNum -= uint32(-num)
		}
	} else {
		d.LikeNum += uint32(num)
	}
	DB.Self.Save(d)

	return d.LikeNum
}

/*--------------- Other Tools -------------*/

// 获取最新插入数据的id
func getLastInsertId() (uint32, error) {
	rows, err := DB.Self.Raw("select last_insert_id()").Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id uint32
	if err := DB.Self.ScanRows(rows, id); err != nil {
		return 0, err
	}
	return id, nil
}

// 是否是一级评论
func isTopComment(id uint32) bool {
	var data CommentModel
	DB.Self.First(&data, "id = ?", id)

	// 还有数据不存在的情况

	return data.IsRoot
}

// Get parentId by commentTargetId
func getParentId(id uint32) (uint32, error) {
	var data = new(CommentModel)
	DB.Self.Where("id = ?", id).First(data)

	// 还有数据不存在的情况

	return data.ParentId, nil
}

// 新增评课时更新课程的评课信息，先暂时放这里，避免冲突
func UpdateCourseByEva(id uint32, rate uint8) error {
	var c UsingCourseModel
	DB.Self.Find(&c, "id = ?", id)

	c.Rate = (c.Rate*float32(c.StarsNum) + float32(rate)) / float32(c.StarsNum+1)
	c.StarsNum++
	DB.Self.Save(&c)

	return nil
}
