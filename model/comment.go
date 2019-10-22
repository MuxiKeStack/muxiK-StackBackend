package model

import (
	"errors"
)

/*---------------------- Course Evaluation Operation ---------------------*/

// Create new course evaluation.
func (evaluation *CourseEvaluationModel) New() error {
	d := DB.Self.Create(evaluation)
	return d.Error
}

// Delete course evaluation.
func (evaluation *CourseEvaluationModel) Delete() error {
	d := DB.Self.Delete(&evaluation)
	return d.Error
}

// Judge whether a course evaluation has already liked by the current user.
func (evaluation *CourseEvaluationModel) HasLiked(userId uint32) bool {
	var data = &EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}
	var count int
	DB.Self.Find(data).Count(&count)
	return count > 0
}

// Like a course evaluation by the current user.
func (evaluation *CourseEvaluationModel) Like(userId uint32) error {
	var data = &EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}
	//if evaluation.HasLiked(userId) {
	//	return errors.New("Has already liked. ")
	//}
	d := DB.Self.Create(data)
	return d.Error
}

// Cancel liking a course evaluation by the current user.
func (evaluation *CourseEvaluationModel) CancelLiking(userId uint32) error {
	var data = &EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}
	//if !evaluation.HasLiked(userId) {
	//	return errors.New("Has not liked yet. ")
	//}
	d := DB.Self.Delete(data)
	return d.Error
}

// Update liked number of a course evaluation after liking or canceling it.
func (evaluation *CourseEvaluationModel) UpdateLikeNum(num int) error {
	likeNum := int(evaluation.LikeNum)
	if likeNum == 0 {
		return nil
	}
	likeNum += num
	evaluation.LikeNum = uint32(likeNum)
	d := DB.Self.Save(evaluation)
	return d.Error
}

// Get the response data information of a course evaluation.
func (evaluation *CourseEvaluationModel) GetInfo(userId uint32, visitor bool) (*EvaluationInfo, error) {
	var err error
	var u = &UserInfo{}
	if !evaluation.IsAnonymous {
		u, err = GetUserInfoById(evaluation.UserId)
		if err != nil {
			return &EvaluationInfo{}, err
		}
	}

	// 获取教师名
	course := &UsingCourseModel{}
	DB.Self.First(course, "hash = ?", evaluation.CourseId)

	var isLike = false
	if !visitor {
		isLike = evaluation.HasLiked(userId)
	}

	var info = &EvaluationInfo{
		Id:                  evaluation.Id,
		CourseId:            evaluation.CourseId,
		CourseName:          evaluation.CourseName,
		Teacher:             course.Teacher,
		Rate:                evaluation.Rate,
		AttendanceCheckType: evaluation.AttendanceCheckType,
		ExamCheckType:       evaluation.ExamCheckType,
		Content:             evaluation.Content,
		Time:                evaluation.Time,
		IsAnonymous:         evaluation.IsAnonymous,
		IsLike:              isLike,
		LikeNum:             evaluation.LikeNum,
		CommentNum:          evaluation.CommentNum,
		Tags:                TagStrToArray(evaluation.Tags),
		UserInfo:            u,
	}
	return info, nil
}

// Get course evaluations.
func GetEvaluations(lastId, size int32) (*[]CourseEvaluationModel, error) {
	var evaluations *[]CourseEvaluationModel
	if lastId != -1 {
		DB.Self.Where("id < ?", lastId).Order("id desc").Find(evaluations).Limit(size)
	} else {
		DB.Self.Order("id desc").Find(evaluations).Limit(size)
	}

	return evaluations, nil
}

// Get course evaluation by evaluationId.
func GetEvaluationById(id uint32) (*CourseEvaluationModel, error) {
	var evaluation CourseEvaluationModel
	d := DB.Self.First(&evaluation, "id = ?", id)
	return &evaluation, d.Error
}

/*---------------------------- Comment Operation --------------------------*/

// Create new comment.
func (comment *CommentModel) New() error {
	if comment.IsRoot {
		comment.ParentId = 0
	} else {
		parentId, err := GetParentIdByCommentTargetId(comment.CommentTargetId)
		if err != nil {
			return err
		}
		comment.ParentId = parentId
	}

	d := DB.Self.Create(comment)
	return d.Error
}

// Judge whether a comment has already liked by the current user.
func (comment *CommentModel) HasLiked(userId uint32) bool {
	var data = &CommentLikeModel{
		CommentId: comment.Id,
		UserId:    userId,
	}
	var count int
	DB.Self.Find(data).Count(&count)
	return count > 0
}

// Like a comment by the current user.
func (comment *CommentModel) Like(userId uint32) error {
	var data = &CommentLikeModel{
		CommentId: comment.Id,
		UserId:    userId,
	}
	if comment.HasLiked(userId) {
		return errors.New("Have already liked ")
	}
	d := DB.Self.Create(data)
	return d.Error
}

// Cancel liking a course evaluation by the current user.
func (comment *CommentModel) CancelLiking(userId uint32) error {
	var data = &CommentLikeModel{
		CommentId: comment.Id,
		UserId:    userId,
	}
	if !comment.HasLiked(userId) {
		return errors.New("Have not liked ")
	}
	d := DB.Self.Delete(data)
	return d.Error
}

// Update liked number of a comment after liking or canceling it.
func (comment *CommentModel) UpdateLikeNum(num int) error {
	likeNum := int(comment.LikeNum)
	if likeNum == 0 {
		return nil
	}
	likeNum += num
	comment.LikeNum = uint32(likeNum)
	d := DB.Self.Save(comment)
	return d.Error
}

// Get the response data information of a comment.
func (comment *CommentModel) GetInfo(userId uint32, visitor bool) (*CommentInfo, error) {
	commentUser, err := GetUserInfoById(comment.UserId)
	if err != nil {
		return nil, nil
	}

	targetUser, err := GetUserInfoById(comment.CommentTargetId)
	if err != nil {
		return nil, nil
	}

	var isLike = false
	if !visitor {
		isLike = comment.HasLiked(userId)
	}

	data := &CommentInfo{
		Id:             comment.Id,
		Content:        comment.Content,
		LikeNum:        comment.LikeNum,
		IsLike:         isLike,
		Time:           comment.Time,
		UserInfo:       commentUser,
		TargetUserInfo: targetUser,
	}

	return data, nil
}

// Get the response data information of a parentComment.
func (comment *CommentModel) GetParentCommentInfo(userId uint32, visitor bool, subComments *[]CommentInfo) (*ParentCommentInfo, error)  {
	userInfo, err := GetUserInfoById(comment.UserId)
	if err != nil {
		return nil, err
	}

	var isLike = false
	if !visitor {
		isLike = comment.HasLiked(userId)
	}

	info := &ParentCommentInfo{
		CommentId:       comment.Id,
		Content:         comment.Content,
		LikeNum:         comment.LikeNum,
		IsLike:          isLike,
		Time:            comment.Time,
		UserInfo:        userInfo,
		SubCommentsNum:  comment.SubCommentNum,
		SubCommentsList: subComments,
	}
	return info, nil
}

// Get parentComments by evaluationId.
func GetParentComments(EvaluationId uint32, lastId, size int32) (*[]CommentModel, uint32, error) {
	var count uint32
	var comments []CommentModel

	if lastId != -1 {
		DB.Self.Where("is_root = ? AND comment_target_id = ?", true, EvaluationId).
			Find(&comments).Count(&count).Limit(size)
	} else {
		DB.Self.Where("id < ? AND is_root = ? AND comment_target_id = ?", lastId, true, EvaluationId).
			Find(&comments).Count(&count).Limit(size)
	}

	return &comments, count, nil
}

// Get subComments by their parentId.
func GetSubComments(ParentId uint32) (*[]CommentModel, error) {
	var subComments []CommentModel
	DB.Self.Find(&subComments, "parent_id = ?", ParentId)
	return &subComments, nil
}

// Get a comment by its id.
func GetCommentById(id uint32) (*CommentModel, error) {
	var comment CommentModel
	d := DB.Self.First(&comment, "id = ?", id)
	return &comment, d.Error
}

// Get parentId by commentTargetId
func GetParentIdByCommentTargetId(id uint32) (uint32, error) {
	var data  CommentModel
	d := DB.Self.Where("id = ?", id).First(&data)

	return data.ParentId, d.Error
}

/*--------------- Course Operation -------------*/

// 新增评课时更新课程的评课信息，先暂时放这里，避免冲突
func UpdateCourseRateByEvaluation(id uint32, rate uint8) error {
	var c UsingCourseModel
	DB.Self.Find(&c, "id = ?", id)

	c.Rate = (c.Rate*float32(c.StarsNum) + float32(rate)) / float32(c.StarsNum+1)
	c.StarsNum++
	DB.Self.Save(&c)

	return nil
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
