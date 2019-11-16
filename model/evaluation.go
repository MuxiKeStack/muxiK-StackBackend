package model

import "github.com/jinzhu/gorm"

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
	DB.Self.Where("user_id = ? AND evaluation_id = ? ", userId, evaluation.Id).Find(&data).Count(&count)
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

// Update liked number of a course evaluation after liking or canceling it.
//func (evaluation *CourseEvaluationModel) UpdateLikeNum(num int) error {
//	likeNum := int(evaluation.LikeSum)
//	if likeNum == 0 && num == -1 {
//		return nil
//	}
//	likeNum += num
//	d := DB.Self.Model(evaluation).Update("like_sum", likeNum)
//	return d.Error
//}

// Get evaluation by its id.
func (evaluation *CourseEvaluationModel) GetById() error {
	d := DB.Self.First(evaluation, "id = ?", evaluation.Id)
	return d.Error
}

func (evaluation *CourseEvaluationModel) UpdateCommentNum(n int) error {
	num := int(evaluation.CommentNum)
	if num == 0 && n == -1 {
		return nil
	}
	num += n
	d := DB.Self.Model(evaluation).Update("comment_num", num)
	return d.Error
}

// Get evaluation's total like account by id.
func GetEvaluationLikeSum(id uint32) (count uint32) {
	var data EvaluationLikeModel
	DB.Self.Where("evaluation_id = ?", id).Find(&data).Count(&count)
	return
}

// Get course evaluations.
func GetEvaluations(lastId, limit int32) (*[]CourseEvaluationModel, error) {
	var evaluations []CourseEvaluationModel
	var d *gorm.DB
	if lastId != 0 {
		d = DB.Self.Where("id < ?", lastId).Order("id desc").Find(&evaluations).Limit(limit)
	} else {
		d = DB.Self.Order("id desc").Find(&evaluations).Limit(limit)
	}

	return &evaluations, d.Error
}

/*--------------- Course Operation -------------*/

func (course *HistoryCourseModel) TableName() string {
	return "history_course"
}

// 新增评课时更新课程的评课信息，先暂时放这里，避免冲突
func UpdateCourseRateByEvaluation(id string, rate float32) error {
	var c HistoryCourseModel
	if d := DB.Self.Find(&c, "hash = ?", id); d.Error != nil {
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
