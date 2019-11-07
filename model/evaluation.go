package model

/*-------------------------- Course Evaluation Operation --------------------------*/

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

	d := DB.Self.Create(data)
	return d.Error
}

// Cancel liking a course evaluation by the current user.
func (evaluation *CourseEvaluationModel) CancelLiking(userId uint32) error {
	var data = &EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}

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

/*--------------- Course Operation -------------*/

// 新增评课时更新课程的评课信息，先暂时放这里，避免冲突
func UpdateCourseRateByEvaluation(id uint32, rate float32) error {
	var c UsingCourseModel
	DB.Self.Find(&c, "id = ?", id)

	c.Rate = (c.Rate*float32(c.StarsNum) + rate) / float32(c.StarsNum+1)
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
