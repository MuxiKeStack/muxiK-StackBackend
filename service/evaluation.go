package service

import (
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type EvaluationInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.EvaluationInfo
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
