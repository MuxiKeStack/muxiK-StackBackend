package model

import (
	"errors"
)

/*--------------- Course Evaluation Operation -------------*/

func (evaluation *CourseEvaluationModel) New() error {
	d := DB.Self.Create(evaluation)
	return d.Error
}

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
		isLike = evaluation.HaveLiked(userId)
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

func (evaluation *CourseEvaluationModel) Delete() error {
	d := DB.Self.Delete(&evaluation)
	return d.Error
}

func (evaluation *CourseEvaluationModel) HaveLiked(userId uint32) bool {
	var data = &EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}
	var count int
	DB.Self.Find(data).Count(&count)
	return count > 0
}

func (evaluation *CourseEvaluationModel) Like(userId uint32) error {
	var data = &EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}
	if evaluation.HaveLiked(userId) {
		return errors.New("Have already liked ")
	}
	d := DB.Self.Create(data)
	return d.Error
}

func (evaluation *CourseEvaluationModel) Dislike(userId uint32) error {
	var data = &EvaluationLikeModel{
		EvaluationId: evaluation.Id,
		UserId:       userId,
	}
	if !evaluation.HaveLiked(userId) {
		return errors.New("Have not liked ")
	}
	d := DB.Self.Delete(data)
	return d.Error
}

func (evaluation *CourseEvaluationModel) UpdateLikeNum(num int) {
	likeNum := int(evaluation.LikeNum)
	if likeNum == 0 {
		return
	}
	likeNum += num
	evaluation.LikeNum = uint32(likeNum)
	DB.Self.Save(evaluation)
}

func GetEvaluations(lastId, size int32) (*[]CourseEvaluationModel, error) {
	var evaluations *[]CourseEvaluationModel
	if lastId != -1 {
		DB.Self.Where("id < ?", lastId).Order("id desc").Find(evaluations).Limit(size)
	} else {
		DB.Self.Order("id desc").Find(evaluations).Limit(size)
	}

	return evaluations, nil
}

func GetEvaluationById(id uint32) (*CourseEvaluationModel, error) {
	var evaluation CourseEvaluationModel
	d := DB.Self.First(&evaluation, "id = ?", id)
	return &evaluation, d.Error
}

/*--------------- Comment Operation -------------*/

func (comment *CommentModel) New() error {
	if comment.IsRoot {
		comment.ParentId = 0
	} else {
		parentId, err := getParentId(comment.CommentTargetId)
		if err != nil {
			return err
		}
		comment.ParentId = parentId
	}

	DB.Self.Create(comment)

	newCommentId, _ := getLastInsertId()

	// 错误情况，未能获取id

	comment.Id = newCommentId

	return nil
}

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
		isLike = comment.HaveLiked(userId)
	}

	data := &CommentInfo{
		Content:        comment.Content,
		LikeNum:        comment.LikeNum,
		IsLike:         isLike,
		Time:           comment.Time,
		UserInfo:       commentUser,
		TargetUserInfo: targetUser,
	}

	return data, nil
}

func (comment *CommentModel) GetParentCommentInfo(userId uint32, visitor bool, subComments *[]CommentInfo) (*ParentCommentInfo, error)  {
	userInfo, err := GetUserInfoById(comment.UserId)
	if err != nil {
		return nil, err
	}

	var isLike = false
	if !visitor {
		isLike = comment.HaveLiked(userId)
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

func (comment *CommentModel) HaveLiked(userId uint32) bool {
	var data = &CommentLikeModel{
		CommentId: comment.Id,
		UserId:    userId,
	}
	var count int
	DB.Self.Find(data).Count(&count)
	return count > 0
}

func (comment *CommentModel) Like(userId uint32) error {
	var data = &CommentLikeModel{
		CommentId: comment.Id,
		UserId:    userId,
	}
	if comment.HaveLiked(userId) {
		return errors.New("Have already liked ")
	}
	d := DB.Self.Create(data)
	return d.Error
}

func (comment *CommentModel) Dislike(userId uint32) error {
	var data = &CommentLikeModel{
		CommentId: comment.Id,
		UserId:    userId,
	}
	if !comment.HaveLiked(userId) {
		return errors.New("Have not liked ")
	}
	d := DB.Self.Delete(data)
	return d.Error
}

func (comment *CommentModel) UpdateLikeNum(num int) {
	likeNum := int(comment.LikeNum)
	if likeNum == 0 {
		return
	}
	likeNum += num
	comment.LikeNum = uint32(likeNum)
	DB.Self.Save(comment)
}

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

func GetSubComments(ParentId uint32) (*[]CommentModel, error) {
	var subComments []CommentModel
	DB.Self.Find(&subComments, "parent_id = ?", ParentId)
	return &subComments, nil
}

func GetCommentById(id uint32) (*CommentModel, error) {
	var comment CommentModel
	d := DB.Self.First(&comment, "id = ?", id)
	return &comment, d.Error
}

/*--------------- Course Operation -------------*/

func (course *UsingCourseModel) UpdateRateByNewEvaluation() {

}
