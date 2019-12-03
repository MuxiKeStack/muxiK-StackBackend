package model

import (
	"github.com/jinzhu/gorm"
)

func (evaluation *CourseEvaluationModel) TableName() string {
	return "course_evaluation"
}

func (data *EvaluationLikeModel) TableName() string {
	return "course_evaluation_like"
}

/*-------------------------- Course Evaluation Operation --------------------------*/

// Create new course evaluation.
func (evaluation *CourseEvaluationModel) New() error {
	d := DB.Self.Create(evaluation)
	return d.Error
}

// Delete course evaluation.
func (evaluation *CourseEvaluationModel) Delete() error {
	d := DB.Self.Delete(evaluation)
	return d.Error
}

// Judge whether a course evaluation has already liked by the current user.
func (evaluation *CourseEvaluationModel) HasLiked(userId uint32) bool {
	var data EvaluationLikeModel
	var count int
	DB.Self.Where("user_id = ? AND evaluation_id = ? ", userId, evaluation.Id).First(&data).Count(&count)
	return count > 0
}

// Like a course evaluation by the current user.
func (evaluation *CourseEvaluationModel) Like(userId uint32) error {
	var data = EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}

	d := DB.Self.Create(&data)
	return d.Error
}

// Cancel liking a course evaluation by the current user.
func (evaluation *CourseEvaluationModel) CancelLiking(userId uint32) error {
	var data = EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}

	d := DB.Self.Delete(&data)
	return d.Error
}

// Get evaluation by its id.
func (evaluation *CourseEvaluationModel) GetById() error {
	d := DB.Self.Unscoped().First(evaluation, "id = ?", evaluation.Id)
	return d.Error
}

// Update evaluation's total comment account.
func (evaluation *CourseEvaluationModel) UpdateCommentNum(n int) error {
	num := int(evaluation.CommentNum)
	if num == 0 && n == -1 {
		return nil
	}
	num += n
	d := DB.Self.Model(evaluation).Update("comment_num", num)
	return d.Error
}

// Update evaluation's total like account.
func (evaluation *CourseEvaluationModel) UpdateLikeNum(n int) error {
	num := int(evaluation.LikeNum)
	if num == 0 && n == -1 {
		return nil
	}
	num += n
	d := DB.Self.Model(evaluation).Update("like_num", num)
	return d.Error
}

// Get evaluation's total like account by id.
func GetEvaluationLikeSum(id uint32) (count uint32) {
	var data EvaluationLikeModel
	DB.Self.Where("evaluation_id = ?", id).Find(&data).Count(&count)
	return
}

// Get all course evaluations.
func GetEvaluations(lastId, limit int32) (*[]CourseEvaluationModel, error) {
	var evaluations []CourseEvaluationModel
	var d *gorm.DB
	if lastId != 0 {
		d = DB.Self.Unscoped().Where("id < ?", lastId).Order("id desc").Limit(limit).Find(&evaluations)
	} else {
		d = DB.Self.Unscoped().Order("id desc").Limit(limit).Find(&evaluations)
	}

	if d.RecordNotFound() {
		return &evaluations, nil
	}
	return &evaluations, d.Error
}

// Get a course's all evaluations by id order by time.
func GetEvaluationsByCourseIdOrderByTime(id string, lastId, limit int32) (*[]CourseEvaluationModel, error) {
	var evaluations []CourseEvaluationModel
	var d *gorm.DB
	if lastId != 0 {
		d = DB.Self.Unscoped().Where("id < ? AND course_id = ?", lastId, id).Order("id desc").Limit(limit).Find(&evaluations)
	} else {
		d = DB.Self.Unscoped().Where("course_id = ?", id).Order("id desc").Limit(limit).Find(&evaluations)
	}

	if d.RecordNotFound() {
		return &evaluations, nil
	}
	return &evaluations, d.Error
}

// Get a course's hot evaluations by id.
func GetEvaluationsByCourseIdOrderByLikeNum(courseId string, limit int32) (*[]CourseEvaluationModel, error) {
	var evaluations []CourseEvaluationModel
	d := DB.Self.Where("course_id = ? AND like_num > 0", courseId).Order("like_num desc, id desc").Limit(limit).Find(&evaluations)

	if d.RecordNotFound() {
		return &evaluations, nil
	}
	return &evaluations, d.Error
}

// Get user's evaluations.
func GetEvaluationsByUserId(userId uint32, lastId, limit int32) (*[]CourseEvaluationModel, error) {
	var evaluations []CourseEvaluationModel
	var d *gorm.DB
	if lastId != 0 {
		d = DB.Self.Unscoped().Where("id < ? AND user_id = ?", lastId, userId).Order("id desc").Limit(limit).Find(&evaluations)
	} else {
		d = DB.Self.Unscoped().Where("user_id = ?", userId).Order("id desc").Limit(limit).Find(&evaluations)
	}

	if d.RecordNotFound() {
		return &evaluations, nil
	}
	return &evaluations, d.Error
}

// Whether user has evaluated the course.
func HasEvaluated(userId uint32, courseId string) bool {
	var evaluation CourseEvaluationModel
	d := DB.Self.Where("userId = ? AND course_id = ?", userId, courseId).First(&evaluation)
	if d.RecordNotFound() {
		return false
	}
	return true
}

/*--------------- Course Operation -------------*/

// 新增评课时更新课程的评课信息，先暂时放这里，避免冲突
func UpdateCourseRateByEvaluation(id string, rate float32) error {
	var c HistoryCourseModel
	if d := DB.Self.First(&c, "hash = ?", id); d.Error != nil {
		return d.Error
	}

	c.Rate = (c.Rate*float32(c.StarsNum) + rate) / float32(c.StarsNum+1)
	c.StarsNum++

	d := DB.Self.Save(&c)
	return d.Error
}

// 根据课程id获取教师名
func GetTeacherByCourseId(id string) (string, error) {
	var course HistoryCourseModel
	d := DB.Self.First(&course, "hash = ?", id)
	return course.Teacher, d.Error
}

// 判断课程是否存在
func IsCourseExisting(id string) bool {
	var course HistoryCourseModel
	d := DB.Self.Where("hash = ?", id).First(&course)
	if d.RecordNotFound() {
		return false
	}
	return true
}
