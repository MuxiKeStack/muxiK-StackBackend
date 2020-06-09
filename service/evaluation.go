package service

import (
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/constvar"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

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
		log.Error("GetEvaluations function error.", err)
		return nil, err
	}

	return GetEvaluationInfosByOriginModels(evaluations, userId, visitor)
}

// Get evaluations of one course.
func GetEvaluationsOfOneCourse(lastId, limit int32, userId uint32, visitor bool, courseId string) (*[]model.EvaluationInfo, error) {
	evaluations, err := model.GetEvaluationsByCourseIdOrderByTime(courseId, lastId, limit)
	if err != nil {
		log.Error("GetEvaluationsByCourseIdOrderByTime function error.", err)
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

			data, err := GetEvaluationInfo(&evaluation, userId, visitor)
			if err != nil {
				log.Error("GetEvaluationInfo function error.", err)
				errChan <- err
				return
			}

			evaluationInfoList.Lock.Lock()
			defer evaluationInfoList.Lock.Unlock()

			evaluationInfoList.IdMap[data.Id] = data

		}(evaluation)
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

// Get the response data information of a course evaluation by original evaluation model.
func GetEvaluationInfo(evaluation *model.CourseEvaluationModel, userId uint32, visitor bool) (*model.EvaluationInfo, error) {
	var err error

	// Get evaluation user info if not anonymous
	u := &model.UserInfoResponse{}
	if !evaluation.IsAnonymous {
		u, err = GetUserInfoById(evaluation.UserId)
		if err != nil {
			log.Error("GetUserInfoById function error.", err)
			return nil, err
		}
	}

	// Get teacher
	teacher, err := model.GetTeacherByCourseId(evaluation.CourseId)
	if err != nil {
		log.Error("GetTeacherByCourseId function error.", err)
		return nil, err
	}

	// Get like state
	var isLike = false
	if !visitor {
		_, isLike = evaluation.HasLiked(userId)
	}

	// Whether the evaluation can be deleted by the user
	canDelete := false
	if !visitor && evaluation.UserId == userId && evaluation.DeletedAt == nil {
		canDelete = true
	}

	// Get tag names
	tagNames, err := GetTagNamesByIdStr(evaluation.Tags)
	if err != nil {
		log.Error("GetTagNamesByIdStr function error.", err)
		return nil, err
	}

	date, time := util.ParseTime(evaluation.Time)

	var info = &model.EvaluationInfo{
		Id:                  evaluation.Id,
		CourseId:            evaluation.CourseId,
		CourseName:          evaluation.CourseName,
		Teacher:             teacher,
		Rate:                evaluation.Rate,
		AttendanceCheckType: GetAttendanceCheckTypeByCode(evaluation.AttendanceCheckType),
		ExamCheckType:       GetExamCheckTypeByCode(evaluation.ExamCheckType),
		Content:             evaluation.Content,
		Date:                date,
		Time:                time,
		IsAnonymous:         evaluation.IsAnonymous,
		IsLike:              isLike,
		LikeNum:             evaluation.LikeNum,
		CommentNum:          evaluation.CommentNum,
		Tags:                tagNames,
		UserInfo:            u,
		IsValid:             true,
		CanDelete:           canDelete,
	}

	// The evaluation has been deleted or been reported
	if evaluation.DeletedAt != nil || !evaluation.IsValid {
		info.IsValid = false
		info.Content = ""
		info.AttendanceCheckType = ""
		info.ExamCheckType = ""
		info.Rate = 0
		info.Tags = nil
	}

	return info, nil
}

// Get hot evaluations of one course.
func GetHotEvaluations(courseId string, limit int32, userId uint32, visitor bool) (*[]model.EvaluationInfo, error) {
	evaluations, err := model.GetEvaluationsByCourseIdOrderByLikeNum(courseId, limit)
	if err != nil {
		log.Error("GetEvaluationsByCourseIdOrderByLikeNum functions error.", err)
		return nil, err
	}

	return GetEvaluationInfosByOriginModels(evaluations, userId, visitor)
}

// Get history-evaluations for user by its id.
func GetHistoryEvaluationsByUserId(userId uint32, lastId, limit int32) (*[]model.EvaluationInfo, error) {
	evaluations, err := model.GetEvaluationsByUserId(userId, lastId, limit)
	if err != nil {
		log.Error("GetEvaluationsByUserId functions error.", err)
		return nil, err
	}

	return GetEvaluationInfosByOriginModels(evaluations, userId, false)
}

// 序号不可以为0, 0代表没有人评论过的情况

// Get attendance-check type name by identifier code.
func GetAttendanceCheckTypeByCode(code uint8) string {
	if attendance, ok := constvar.Attendance[code]; ok {
		return attendance
	}
	return ""
}

// Get exam-check type name by identifier code.
func GetExamCheckTypeByCode(code uint8) string {
	if exam, ok := constvar.Exam[code]; ok {
		return exam
	}
	return ""
}

// Get attendance-check type name by identifier code.
func GetAttendanceCheckTypeByCodeEnglish(code uint8) string {
	if attendance, ok := constvar.AttendanceEnglish[code]; ok {
		return attendance
	}
	return ""
}

// Get exam-check type name by identifier code.
func GetExamCheckTypeByCodeEnglish(code uint8) string {
	if exam, ok := constvar.ExamEnglish[code]; ok {
		return exam
	}
	return ""
}

// Get attendance check type amount by course hash id for course info.
func GetAttendanceCheckTypeNumForCourseInfo(courseId string) map[string]uint32 {
	result := make(map[string]uint32)
	for i := 1; i < 4; i++ {
		name := GetAttendanceCheckTypeByCode(uint8(i))
		result[name] = model.GetAttendanceTypeNumChosenByCode(courseId, i)
	}
	return result
}

// Get exam check type amount by course hash id for course info.
func GetExamCheckTypeNumForCourseInfo(courseId string) map[string]uint32 {
	result := make(map[string]uint32)
	for i := 1; i < 5; i++ {
		name := GetExamCheckTypeByCode(uint8(i))
		result[name] = model.GetExamCheckTypeNumChosenByCode(courseId, i)
	}
	return result
}

// Get attendance check type amount by course hash id for course info.
func GetAttendanceCheckTypeNumForCourseInfoEnglish(courseId string) map[string]uint32 {
	result := make(map[string]uint32)
	for i := 1; i < 4; i++ {
		name := GetAttendanceCheckTypeByCodeEnglish(uint8(i))
		result[name] = model.GetAttendanceTypeNumChosenByCode(courseId, i)
	}
	return result
}

// Get exam check type amount by course hash id for course info.
func GetExamCheckTypeNumForCourseInfoEnglish(courseId string) map[string]uint32 {
	result := make(map[string]uint32)
	for i := 1; i < 5; i++ {
		name := GetExamCheckTypeByCodeEnglish(uint8(i))
		result[name] = model.GetExamCheckTypeNumChosenByCode(courseId, i)
	}
	return result
}
