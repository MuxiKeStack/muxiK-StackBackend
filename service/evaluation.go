package service

import (
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/lexkong/log"
)

type EvaluationInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.EvaluationInfo
}

// Get course evaluation list.
func EvaluationList(lastId, size int32, userId uint32, visitor bool) (*[]model.EvaluationInfo, error) {
	log.Infof("EvaluationList is called")

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

	log.Infof("EvaluationList is returned.")
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
		u, err = GetUserInfoById(evaluation.UserId)
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
		AttendanceCheckType: GetAttendanceCheckTypeByCode(evaluation.AttendanceCheckType),
		ExamCheckType:       GetExamCheckTypeByCode(evaluation.ExamCheckType),
		Content:             evaluation.Content,
		Time:                evaluation.Time,
		IsAnonymous:         evaluation.IsAnonymous,
		IsLike:              isLike,
		LikeNum:             evaluation.LikeNum,
		CommentNum:          evaluation.CommentNum,
		Tags:                GetTagNamesByIdStr(evaluation.Tags),
		UserInfo:            u,
	}

	return info, nil
}

// Get attendance-check type name by identifier code.
func GetAttendanceCheckTypeByCode(code uint8) string {
	switch code {
	case 0:
		return "经常点名"
	case 1:
		return "偶尔点名"
	case 2:
		return "签到点名"
	}
	return ""
}

// Get exam-check type name by identifier code.
func GetExamCheckTypeByCode(code uint8) string {
	switch code {
	case 0:
		return "无考核"
	case 1:
		return "闭卷考试"
	case 2:
		return "开卷考试"
	case 3:
		return "论文考核"
	}
	return ""
}
