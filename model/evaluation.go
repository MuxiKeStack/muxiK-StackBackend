package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func (evaluation *CourseEvaluationModel) TableName() string {
	return "course_evaluation"
}

func (data *EvaluationLikeModel) TableName() string {
	return "course_evaluation_like"
}

func NewEvaluation(id uint32) *CourseEvaluationModel {
	return &CourseEvaluationModel{Id: id}
}

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

// Block a evaluation, cause of be reported > 5 times
func (evaluation *CourseEvaluationModel) Block() error {
	evaluation.IsValid = false
	d := DB.Self.Update(evaluation)
	if d.Error != nil {
		return d.Error
	}
	return nil
}

// Judge whether a course evaluation has already liked by the current user,
// return like-record id and bool type.
func (evaluation *CourseEvaluationModel) HasLiked(userId uint32) (uint32, bool) {
	var data EvaluationLikeModel
	d := DB.Self.Where("user_id = ? AND evaluation_id = ? ", userId, evaluation.Id).First(&data)
	return data.Id, !d.RecordNotFound()
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

// Cancel liking a course evaluation by the like-record id.
func (evaluation *CourseEvaluationModel) CancelLiking(id uint32) error {
	var data = EvaluationLikeModel{Id: id}
	d := DB.Self.Delete(&data)
	return d.Error
}

// Get evaluation by its id.
func (evaluation *CourseEvaluationModel) GetById() error {
	return DB.Self.Unscoped().Where("id = ?", evaluation.Id).First(evaluation).Error
}

// Update evaluation's total comment amount.
func (evaluation *CourseEvaluationModel) UpdateCommentNum(n int) error {
	num := int(evaluation.CommentNum)
	if num == 0 && n == -1 {
		return nil
	}
	num += n
	d := DB.Self.Model(evaluation).Update("comment_num", num)
	return d.Error
}

// Update evaluation's total like amount.
func (evaluation *CourseEvaluationModel) UpdateLikeNum(n int) error {
	num := int(evaluation.LikeNum)
	if num == 0 && n == -1 {
		return nil
	}
	num += n
	d := DB.Self.Model(evaluation).Update("like_num", num)
	return d.Error
}

// Get evaluation's total like amount by id.
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
	d := DB.Self.Where("course_id = ? AND like_num > 0", courseId).
		Order("like_num desc, id desc").
		Limit(limit).
		Find(&evaluations)

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
	d := DB.Self.Where("user_id = ? AND course_id = ?", userId, courseId).First(&evaluation)
	return !d.RecordNotFound()
}

// Update evaluation's commentNum by it's id
func UpdateCommentNumById(id uint32) error {
	evaluation := &CourseEvaluationModel{Id: id}
	if err := evaluation.GetById(); err != nil {
		return err
	}

	if err := evaluation.UpdateCommentNum(-1); err != nil {
		return err
	}
	return nil
}

// 获取一个课程中评价最多的点名方式
func GetAttendanceType(courseID string) uint8 {
	var res []struct{ AttendanceCheckType uint8 }
	d := DB.Self.Table("course_evaluation").
		Select("attendance_check_type, count( * ) AS count").
		Where("course_id = ?", courseID).
		Group("attendance_check_type").
		Order("count DESC").
		Limit(1).Scan(&res)
	if d.RecordNotFound() || len(res) == 0 {
		return 0
	}
	return res[0].AttendanceCheckType
}

// Get attendance check type amount of a course by identifier.
func GetAttendanceTypeNumChosenByCode(courseId string, code int) (count uint32) {
	DB.Self.Table("course_evaluation").Where("course_id = ? AND attendance_check_type = ?", courseId, code).Count(&count)
	return
}

// 获取一个课程中评价最多的点名方式
func GetExamCheckType(courseID string) uint8 {
	var res []struct{ ExamCheckType uint8 }
	d := DB.Self.Table("course_evaluation").
		Select("exam_check_type, count( * ) AS count").
		Where("course_id = ?", courseID).
		Group("exam_check_type").
		Order("count DESC").
		Limit(1).Scan(&res)
	if d.RecordNotFound() || len(res) == 0 {
		return 0
	}
	return res[0].ExamCheckType
}

// Get exam check type amount of a course by identifier.
func GetExamCheckTypeNumChosenByCode(courseId string, code int) (count uint32) {
	DB.Self.Table("course_evaluation").Where("course_id = ? AND exam_check_type = ?", courseId, code).Count(&count)
	return
}

// 新增评课时更新课程的评课信息
func UpdateCourseRateByEvaluation(id string, rate float32) error {
	var c HistoryCourseModel
	if d := DB.Self.First(&c, "hash = ?", id); d.Error != nil {
		return d.Error
	}

	c.Rate = (c.Rate*float32(c.StarsNum) + rate) / float32(c.StarsNum+1)
	c.StarsNum++

	return DB.Self.Save(&c).Error
}

// Update course's info after deleting an evaluation.
func UpdateCourseInfoAfterDeletingEvaluation(id string, rate float32) error {
	var c HistoryCourseModel
	if d := DB.Self.First(&c, "hash = ?", id); d.Error != nil {
		return d.Error
	}

	rate = c.Rate*float32(c.StarsNum) - rate
	if c.StarsNum <= 0 || rate < 0 {
		return errors.New("Unexpected data error ")
	}
	c.StarsNum--
	if c.StarsNum == 0 {
		c.Rate = 0
	} else {
		c.Rate = rate / float32(c.StarsNum)
	}

	return DB.Self.Save(&c).Error
}

// 通过评课id 获得评课人的userID block时候只知道是评课id，不知道是谁评的 用于消息提醒
func GetUIDByEvaluationID(eid uint32) (uint32, error) {
	var e CourseEvaluationModel
	d := DB.Self.Where("id = ?", eid).First(&e)
	return e.UserId, d.Error
}

// 删除评课，同时更新课程信息，事务
func DeleteEvaluation(evaluation *CourseEvaluationModel) error {
	tx := DB.Self.Begin()

	var course = &HistoryCourseModel{}
	if err := tx.Where("hash = ?", evaluation.CourseId).First(course).Error; err != nil {
		tx.Rollback()
		return err
	}

	rate := course.Rate*float32(course.StarsNum) - evaluation.Rate
	if course.StarsNum <= 0 || rate < 0 {
		tx.Rollback()
		return errors.New("Unexpected data error ")
	}
	course.StarsNum--
	if course.StarsNum == 0 {
		course.Rate = 0
	} else {
		course.Rate = rate / float32(course.StarsNum)
	}

	// 更新课程
	if err := tx.Save(course).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除评课
	if err := tx.Delete(evaluation).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
