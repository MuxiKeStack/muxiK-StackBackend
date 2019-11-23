package service

import (
	"fmt"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

type EvaluationInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.EvaluationInfo
}

// Get course evaluation list for evaluation playground.
func GetEvaluationsForPlayground(lastId, limit int32, userId uint32, visitor bool) (*[]model.EvaluationInfo, error) {
	evaluations, err := model.GetEvaluations(lastId, limit)
	if err != nil {
		log.Info("GetEvaluations function error.")
		return nil, err
	}

	return GetEvaluationInfosByOriginModels(evaluations, userId, visitor)
}

// Get evaluations of one course.
func GetEvaluationsOfOneCourse(lastId, limit int32, userId uint32, visitor bool, courseId, sortKey string) (*[]model.EvaluationInfo, error) {
	var evaluations *[]model.CourseEvaluationModel
	var err error

	if sortKey == "hot" {
		evaluations, err = model.GetEvaluationsByCourseIdOrderByLikeNum(courseId, lastId, limit)
	} else {
		evaluations, err = model.GetEvaluationsByCourseIdOrderByTime(courseId, lastId, limit)
	}

	if err != nil {
		log.Info("GetEvaluationsByCourseId* function error.")
		return nil, err
	}

	return GetEvaluationInfosByOriginModels(evaluations, userId, visitor)
}

func GetEvaluationInfosByOriginModels(evaluations *[]model.CourseEvaluationModel, userId uint32, visitor bool) (*[]model.EvaluationInfo, error) {
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
		go func(evaluation model.CourseEvaluationModel) {
			defer wg.Done()

			data, err := GetEvaluationInfo(evaluation.Id, userId, visitor)
			if err != nil {
				log.Info("GetEvaluationInfo function error.")
				errChan <- err
				return
			}

			evaluationInfoList.Lock.Lock()
			defer evaluationInfoList.Lock.Unlock()

			evaluationInfoList.IdMap[data.Id] = data

		}(evaluation) // 不能传地址
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
		log.Info("evaluation.GetById() error.")
		return nil, err
	}

	// Get evaluation user info if not anonymous
	u := &model.UserInfoResponse{}
	if !evaluation.IsAnonymous {
		u, err = GetUserInfoById(evaluation.UserId)
		if err != nil {
			log.Info("GetUserInfoById function error.")
			return nil, err
		}
	}

	// Get teacher
	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		log.Info("GetTeacherByCourseId function error.")
		return nil, err
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = evaluation.HasLiked(userId)
	}

	// Get tag names
	tagNames, err := GetTagNamesByIdStr(evaluation.Tags)
	if err != nil {
		fmt.Println(tagNames, evaluation.Tags)
		log.Info("GetTagNamesByIdStr function error.")
		return nil, err
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
		Time:                evaluation.Time.Unix(),
		IsAnonymous:         evaluation.IsAnonymous,
		IsLike:              isLike,
		LikeNum:             evaluation.LikeNum,
		CommentNum:          evaluation.CommentNum,
		Tags:                tagNames,
		UserInfo:            u,
		IsValid:             true,
	}

	return info, nil
}

// Get attendance-check type name by identifier code.
func GetAttendanceCheckTypeByCode(code uint8) string {
	switch code {
	case 1:
		return "经常点名"
	case 2:
		return "偶尔点名"
	case 3:
		return "签到点名"
	}
	return ""
}

// Get exam-check type name by identifier code.
func GetExamCheckTypeByCode(code uint8) string {
	switch code {
	case 1:
		return "无考核"
	case 2:
		return "闭卷考试"
	case 3:
		return "开卷考试"
	case 4:
		return "论文考核"
	}
	return ""
}
