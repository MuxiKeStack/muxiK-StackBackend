package service

import (
	"errors"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

type ParentCommentInfoList struct {
	Lock  *sync.Mutex
	IdMap map[string]*model.ParentCommentInfo
}

type SubCommentInfoList struct {
	Lock  *sync.Mutex
	IdMap map[string]*model.CommentInfo
}

// Get comment list.
func CommentList(evaluationId uint32, limit, offset int32, userId uint32, visitor bool) (*[]model.ParentCommentInfo, error) {
	log.Info("CommentList function is called")

	// Get parent comments from database
	parentComments, err := model.GetParentComments(evaluationId, limit, offset)
	if err != nil {
		log.Error("GetParentComments", err)
		return nil, err
	}

	var parentIds []string
	for _, parentComment := range *parentComments {
		parentIds = append(parentIds, parentComment.Id)
	}

	parentCommentInfoList := ParentCommentInfoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[string]*model.ParentCommentInfo, len(*parentComments)),
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	for _, parentComment := range *parentComments {
		wg.Add(1)
		go func(parentComment model.ParentCommentModel) {
			defer wg.Done()

			// 获取父评论详情
			parentCommentInfo, err := GetParentCommentInfo(parentComment.Id, userId, visitor)
			if err != nil {
				log.Error("GetParentCommentInfo function error", err)
				errChan <- err
				return
			}

			parentCommentInfoList.Lock.Lock()
			defer parentCommentInfoList.Lock.Unlock()

			parentCommentInfoList.IdMap[parentCommentInfo.Id] = parentCommentInfo

		}(parentComment)
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

	var infos []model.ParentCommentInfo
	for _, id := range parentIds {
		infos = append(infos, *parentCommentInfoList.IdMap[id])
	}

	return &infos, nil
}

// Get the response data information of a parent comment.
func GetParentCommentInfo(id string, userId uint32, visitor bool) (*model.ParentCommentInfo, error) {
	// Get comment from Database
	comment := &model.ParentCommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		return nil, err
	}

	// Get the user of the parent comment if not anonymous
	var userInfo *model.UserInfoResponse
	var err error
	if !comment.IsAnonymous {
		userInfo, err = GetUserInfoById(comment.UserId)
		if err != nil {
			log.Error("GetUserInfoById function is called", err)
			return nil, err
		}
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = model.CommentHasLiked(userId, comment.Id)
	}

	// Get subComments' infos
	subCommentInfos, err := GetSubCommentInfosByParentId(comment.Id, userId, visitor)
	if err != nil {
		return nil, err
	}

	data := &model.ParentCommentInfo{
		Id:              comment.Id,
		Content:         comment.Content,
		LikeNum:         model.GetCommentLikeSum(comment.Id),
		IsLike:          isLike,
		IsValid:         true,
		Time:            comment.Time.Unix(),
		IsAnonymous:     comment.IsAnonymous,
		UserInfo:        userInfo,
		SubCommentsNum:  comment.SubCommentNum,
		SubCommentsList: subCommentInfos,
	}

	// The parent comment has been deleted or been reported
	if comment.DeletedAt != nil || !comment.IsValid {
		data.IsValid = false
		data.Content = ""
	}

	return data, nil
}

// Get subComments' infos by parent id.
func GetSubCommentInfosByParentId(id string, userId uint32, visitor bool) (*[]model.CommentInfo, error) {
	// Get subComments from Database
	comments, err := model.GetSubCommentsByParentId(id)
	if err != nil {
		log.Error("GetSubCommentsByParentId function error", err)
		return nil, err
	}

	var commentIds []string
	for _, comment := range *comments {
		commentIds = append(commentIds, comment.Id)
	}

	subCommentInfoList := SubCommentInfoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[string]*model.CommentInfo, len(*comments)),
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 并发获取子评论详情列表
	for _, comment := range *comments {
		wg.Add(1)

		go func(comment model.SubCommentModel) {
			defer wg.Done()

			// Get a subComment's info by its id
			info, err := GetSubCommentInfoById(comment.Id, userId, visitor)
			if err != nil {
				log.Error("GetSubCommentInfoById function error", err)
				errChan <- err
				return
			}

			subCommentInfoList.Lock.Lock()
			defer subCommentInfoList.Lock.Unlock()

			subCommentInfoList.IdMap[info.Id] = info

		}(comment) //传址会panic
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

	var commentInfos []model.CommentInfo
	for _, id := range commentIds {
		commentInfos = append(commentInfos, *subCommentInfoList.IdMap[id])
	}

	return &commentInfos, nil
}

// Get the response information of a subComment by id.
func GetSubCommentInfoById(id string, userId uint32, visitor bool) (*model.CommentInfo, error) {
	// Get comment from Database
	comment := &model.SubCommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		log.Error("comment.GetById function error", err)
		return nil, err
	}

	// Get the user of the subComment if not anonymous
	var commentUser *model.UserInfoResponse
	var err error
	if !comment.IsAnonymous {
		commentUser, err = GetUserInfoById(comment.UserId)
		if err != nil {
			log.Error("GetUserInfoById function error", err)
			return nil, err
		}
	}

	// Get the target user of the subComment if not anonymous (identified by 0)
	var targetUser *model.UserInfoResponse
	if comment.TargetUserId != 0 {
		targetUser, err = GetUserInfoById(comment.TargetUserId)
		if err != nil {
			log.Error("GetUserInfoById function error", err)
			return nil, err
		}
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = model.CommentHasLiked(userId, comment.Id)
	}

	data := &model.CommentInfo{
		Id:             comment.Id,
		Content:        comment.Content,
		LikeNum:        model.GetCommentLikeSum(comment.Id),
		IsLike:         isLike,
		IsValid:        true,
		Time:           comment.Time.Unix(),
		UserInfo:       commentUser,
		TargetUserInfo: targetUser,
	}

	// The subComment has been deleted or been reported
	if comment.DeletedAt != nil || !comment.IsValid {
		data.IsValid = false
		data.Content = ""
	}

	return data, nil
}

// Delete a parent comment.
func DeleteParentComment(id string, userId uint32) error {
	// Get evaluation by id
	comment := &model.ParentCommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		log.Error("comment.GetById error.", err)
		return err
	}

	// 验证当前用户是否有删除此评课的权限
	if comment.UserId != userId {
		return errors.New("With no permission to delete the comment ")
	}

	if err := comment.Delete(); err != nil {
		log.Error("comment.Delete() error.", err)
		return err
	}

	return nil
}

// Delete a subComment.
func DeleteSubComment(id string, userId uint32) error {
	// Get evaluation by id
	comment := &model.SubCommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		log.Error("comment.GetById error.", err)
		return err
	}

	// 验证当前用户是否有删除此评课的权限
	if comment.UserId != userId {
		return errors.New("With no permission to delete the comment ")
	}

	if err := comment.Delete(); err != nil {
		log.Error("comment.Delete() error.", err)
		return err
	}

	return nil
}

// Update liked number of a comment after liking or canceling it.
//func UpdateCommentLikeNum(commentId string, num int) (uint32, error) {
//	log.Info("UpdateCommentLikeNum function is called")
//
//	subComment, ok := model.IsSubComment(commentId)
//	if ok {
//		err := subComment.UpdateLikeNum(num)
//		return subComment.LikeNum, err
//	}
//
//	parentComment := &model.ParentCommentModel{Id: commentId}
//	if err := parentComment.GetById(); err != nil {
//		log.Error("parentComment.GetById function error", err)
//		return 0, err
//	}
//
//	err := parentComment.UpdateLikeNum(num)
//	return parentComment.LikeNum, err
//}
