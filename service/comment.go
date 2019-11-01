package service

import (
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type EvaluationInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.EvaluationInfo
}

type ParentCommentInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.ParentCommentInfo
}

type SubCommentInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.CommentInfo
}

// Get course evaluation list.
func EvaluationList(lastId, size int32, userId uint32, visitor bool) (*[]model.EvaluationInfo, error) {
	evaluations, err := model.GetEvaluations(lastId, size)
	if err != nil {
		return nil, err
	}

	var ids []uint32
	for _, evaluation := range *evaluations {
		ids = append(ids, evaluation.Id)
	}

	evaluationInfoList := EvaluationInfoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint32]*model.EvaluationInfo, len(*evaluations)),
	}

	wg := sync.WaitGroup{}
	finished := make(chan bool, 1)
	errChan := make(chan error, 1)

	for _, evaluation := range *evaluations {
		wg.Add(1)
		go func(evaluation *model.CourseEvaluationModel) {
			defer wg.Done()

			data, err := GetEvaluationInfo(evaluation.Id, userId, visitor)
			if err != nil {
				errChan <- err
				return
			}

			evaluationInfoList.Lock.Lock()
			defer evaluationInfoList.Lock.Unlock()

			evaluationInfoList.IdMap[data.Id] = data

		}(&evaluation)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	var infos []model.EvaluationInfo
	for _, id := range ids {
		infos = append(infos, *evaluationInfoList.IdMap[id])
	}

	return &infos, nil
}

// Get comment list.
func CommentList(evaluationId uint32, limit, offset int32, userId uint32, visitor bool) (*[]model.ParentCommentInfo, uint32, error) {
	// Get parent comments from database
	parentComments, count, err := model.GetParentComments(evaluationId, limit, offset)
	if err != nil {
		return nil, count, err
	}

	var parentIds []uint32
	for _, comment := range *parentComments {
		parentIds = append(parentIds, comment.Id)
	}

	parentCommentInfoList := ParentCommentInfoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint32]*model.ParentCommentInfo, len(*parentComments)),
	}

	var wg1, wg2 *sync.WaitGroup
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	for _, parentComment := range *parentComments {
		wg1.Add(1)
		go func(parentComment *model.CommentModel) {
			defer wg1.Done()

			subComments, err := model.GetSubComments(parentComment.Id)
			if err != nil {
				errChan <- err
				return
			}

			var subCommentIds []uint32
			for _, comment := range *subComments {
				subCommentIds = append(subCommentIds, comment.Id)
			}

			subCommentInfoList := SubCommentInfoList{
				Lock:  new(sync.Mutex),
				IdMap: make(map[uint32]*model.CommentInfo, len(*subComments)),
			}

			errChan2 := make(chan error, 1)
			finished := make(chan bool, 1)

			// 并发获取子评论详情列表
			for _, subComment := range *subComments {
				wg2.Add(1)

				go func(subComment *model.CommentModel) {
					defer wg2.Done()

					info, err := GetCommentInfo(subComment.Id, userId, visitor)
					if err != nil {
						errChan2 <- err
						return
					}

					subCommentInfoList.Lock.Lock()
					defer subCommentInfoList.Lock.Unlock()

					subCommentInfoList.IdMap[info.Id] = info

				}(&subComment)
			}

			go func() {
				wg2.Wait()
				close(finished)
			}()

			select {
			case <-finished:
			case err := <-errChan2:
				errChan <- err
				return
			}

			var subCommentInfos []model.CommentInfo
			for _, id := range subCommentIds {
				subCommentInfos = append(subCommentInfos, *subCommentInfoList.IdMap[id])
			}

			// 获取父评论详情
			parentCommentInfo, err := GetParentCommentInfoById(parentComment.Id ,userId, visitor, &subCommentInfos)
			if err != nil {
				errChan <- err
				return
			}

			parentCommentInfoList.Lock.Lock()
			defer parentCommentInfoList.Lock.Unlock()

			parentCommentInfoList.IdMap[parentCommentInfo.Id] = parentCommentInfo

		}(&parentComment)
	}

	go func() {
		wg1.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	var infos []model.ParentCommentInfo
	for _, id := range parentIds {
		infos = append(infos, *parentCommentInfoList.IdMap[id])
	}

	return &infos, count, nil
}

// Get the response data information of a course evaluation.
func GetEvaluationInfo(id, userId uint32, visitor bool) (*model.EvaluationInfo, error) {
	var err error

	// Get evaluation from Database
	evaluation := &model.CourseEvaluationModel{Id: id}
	if err := evaluation.GetById(); err != nil {
		return nil, err
	}

	// Get evaluation user info if not anonymous
	u := &model.UserInfoResponse{}
	if !evaluation.IsAnonymous {
		u, err = model.GetUserInfoById(evaluation.UserId)
		if err != nil {
			return nil, err
		}
	}

	// Get teacher
	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		return nil, err
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = evaluation.HasLiked(userId)
	}

	var info = &model.EvaluationInfo{
		Id:                  evaluation.Id,
		CourseId:            evaluation.CourseId,
		CourseName:          evaluation.CourseName,
		Teacher:             teacher,
		Rate:                evaluation.Rate,
		AttendanceCheckType: evaluation.AttendanceCheckType,
		ExamCheckType:       evaluation.ExamCheckType,
		Content:             evaluation.Content,
		Time:                evaluation.Time,
		IsAnonymous:         evaluation.IsAnonymous,
		IsLike:              isLike,
		LikeNum:             evaluation.LikeNum,
		CommentNum:          evaluation.CommentNum,
		Tags:                model.TagStrToArray(evaluation.Tags),
		UserInfo:            u,
	}

	return info, nil
}

// Get the response data information of a comment.
func GetCommentInfo(id, userId uint32, visitor bool) (*model.CommentInfo, error) {
	// Get comment from Database
	comment := &model.CommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		return nil, err
	}

	// Get the user of the comment
	commentUser, err := model.GetUserInfoById(comment.UserId)
	if err != nil {
		return nil, nil
	}

	// Get the target user of the comment
	targetUser, err := model.GetUserInfoById(comment.CommentTargetId)
	if err != nil {
		return nil, nil
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = comment.HasLiked(userId)
	}

	data := &model.CommentInfo{
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
func GetParentCommentInfoById(id, userId uint32, visitor bool, subComments *[]model.CommentInfo) (*model.ParentCommentInfo, error) {
	// Get parent comment from Database
	comment := &model.CommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		return nil, err
	}

	// Get user info of the comment
	userInfo, err := model.GetUserInfoById(comment.UserId)
	if err != nil {
		return nil, err
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = comment.HasLiked(userId)
	}

	info := &model.ParentCommentInfo{
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
