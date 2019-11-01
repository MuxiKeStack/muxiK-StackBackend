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

// Get evaluation by its id.
func (evaluation *CourseEvaluationModel) GetById() error {
	d := DB.Self.First(evaluation)
	return d.Error
}

// Get course evaluations.
func GetEvaluations(lastId, limit int32) (*[]CourseEvaluationModel, error) {
	var evaluations *[]CourseEvaluationModel
	if lastId != -1 {
		DB.Self.Where("id < ?", lastId).Order("id desc").Find(evaluations).Limit(limit)
	} else {
		DB.Self.Order("id desc").Find(evaluations).Limit(limit)
	}

	return evaluations, nil
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

// Get a comment by its id.
func (comment *CommentModel) GetById() error {
	d := DB.Self.First(comment)
	return d.Error
}

// Get parentComments by evaluationId.
func GetParentComments(EvaluationId uint32, limit, offset int32) (*[]CommentModel, uint32, error) {
	var count uint32
	var comments []CommentModel

	DB.Self.Where("is_root = ? AND comment_target_id = ?", true, EvaluationId).
		Find(&comments).Limit(limit).Offset(offset).Count(&count)

	return &comments, count, nil
}

// Get subComments by their parentId.
func GetSubComments(ParentId uint32) (*[]CommentModel, error) {
	var subComments []CommentModel
	DB.Self.Find(&subComments, "parent_id = ?", ParentId)
	return &subComments, nil
}

// Get parentId by commentTargetId
func GetParentIdByCommentTargetId(id uint32) (uint32, error) {
	var data CommentModel
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

// 根据课程id获取教师名
func GetTeacherByCourseId(id string) (string, error) {
	var course HistoryCourseModel
	d := DB.Self.First(&course, "hash = ?", id)
	return course.Teacher, d.Error
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
