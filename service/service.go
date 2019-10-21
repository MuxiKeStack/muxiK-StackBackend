package service

import (
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type EvaluationInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.EvaluationInfo
}

type CommentInfoList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.ParentCommentInfo
}

func EvaluationList(lastId, size int32, userId uint32, visitor bool) (*[]model.EvaluationInfo, error)  {
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

			data, err := evaluation.GetInfo(userId, visitor)
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

func CommentList(evaluationId uint32, lastId, size int32, userId uint32, visitor bool) (*[]model.ParentCommentInfo, uint32, error) {
	parentComments, count, err := model.GetParentComments(evaluationId, lastId, size)
	if err != nil {
		return nil, count, err
	}

	var parentIds []uint32
	for _, comment := range *parentComments {
		parentIds = append(parentIds, comment.Id)
	}

	parentCommentInfoList := CommentInfoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint32]*model.ParentCommentInfo, len(*parentComments)),
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	for _, parentComment := range *parentComments {
		wg.Add(1)
		go func(parentComment *model.CommentModel) {
			defer wg.Done()

			var subCommentInfoList []model.CommentInfo
			subComments, err := model.GetSubComments(parentComment.Id)
			if err != nil {
				errChan <- err
				return
			}

			// 优化：并发
			for _, subComment := range *subComments {
				info, err := subComment.GetInfo(userId, visitor)
				if err != nil {
					errChan <- err
					return
				}
				subCommentInfoList = append(subCommentInfoList, *info)
			}

			parentCommentInfo, err := parentComment.GetParentCommentInfo(userId, visitor, &subCommentInfoList)
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
		wg.Wait()
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
