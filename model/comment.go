package model

import (
	"errors"
)

/*---------------------------- Parent Comment Operation --------------------------*/

// Create a new parent comment.
func (comment *ParentCommentModel) New() error {
	d := DB.Self.Create(comment)
	return d.Error
}

// Update liked number of a parent comment after liking or canceling it.
func (comment *ParentCommentModel) UpdateLikeNum(num int) error {
	likeNum := int(comment.LikeNum)
	if likeNum == 0 {
		return nil
	}
	likeNum += num
	comment.LikeNum = uint32(likeNum)
	d := DB.Self.Save(comment)
	return d.Error
}

// Get a parent comment by its id.
func (comment *ParentCommentModel) GetById() error {
	d := DB.Self.First(comment)
	return d.Error
}

// Get parent comments by evaluation id.
func GetParentComments(EvaluationId uint32, limit, offset int32) (*[]ParentCommentModel, uint32, error) {
	var count uint32
	var comments []ParentCommentModel

	d := DB.Self.Where("evaluation_id = ?", EvaluationId).
		Find(&comments).Limit(limit).Offset(offset).Count(&count)

	return &comments, count, d.Error
}

/*---------------------------- SubComment Operation --------------------------*/

// Create a new subComment.
func (comment *SubCommentModel) New() error {
	d := DB.Self.Create(comment)
	return d.Error
}

// Update liked number of a subComment after liking or canceling it.
func (comment *SubCommentModel) UpdateLikeNum(num int) error {
	likeNum := int(comment.LikeNum)
	if likeNum == 0 {
		return nil
	}
	likeNum += num
	comment.LikeNum = uint32(likeNum)
	d := DB.Self.Save(comment)
	return d.Error
}

// Get a subComment by its id.
func (comment *SubCommentModel) GetById() error {
	d := DB.Self.First(comment)
	return d.Error
}

// Get subComments by their parentId.
func GetSubCommentsByParentId(ParentId string) (*[]SubCommentModel, error) {
	var subComments []SubCommentModel
	DB.Self.Find(&subComments, "parent_id = ?", ParentId)
	return &subComments, nil
}

// Judge whether it is a subComment by id, if so then also return the subComment.
func IsSubComment(id string) (*SubCommentModel, bool) {
	var comment SubCommentModel
	DB.Self.Where("id = ?", id).First(&comment)
	if comment.Id == "" {
		return &comment, true
	}
	return nil, false
}


/*---------------------------- Other Comment Operations --------------------------*/

// Like a comment by the current user.
func CommentLike(userId uint32, commentId string) error {
	if CommentHasLiked(userId, commentId) {
		return errors.New("Have already liked ")
	}

	var data = &CommentLikeModel{
		UserId:    userId,
		CommentId: commentId,
	}
	d := DB.Self.Create(data)
	return d.Error
}

// Cancel liking a comment by the current user.
func CommentCancelLiking(userId uint32, commentId string) error {
	if !CommentHasLiked(userId, commentId) {
		return errors.New("Have not liked ")
	}

	var data = &CommentLikeModel{
		UserId:    userId,
		CommentId: commentId,
	}
	d := DB.Self.Delete(data)
	return d.Error
}

// Judge whether a comment has already liked by the current user.
func CommentHasLiked(userId uint32, commentId string) bool {
	var data = &CommentLikeModel{
		UserId:    userId,
		CommentId: commentId,
	}
	var count int
	DB.Self.First(data).Count(&count)
	return count == 1
}

// Get parentId by commentTargetId.
func GetParentIdByCommentTargetId(id string) (string, error) {
	var comment SubCommentModel
	d := DB.Self.Where("id = ?", id).First(&comment)
	if d.Error != nil || comment.ParentId != "" {
		return comment.ParentId, d.Error
	}

	// The comment target is a parentComment instead of a subComment
	var parentComment ParentCommentModel
	d = DB.Self.Where("id = ?", id).Find(&parentComment)
	return parentComment.Id, d.Error
}

// Get comment target user's id by the commentTargetId.
func GetTargetUserIdByCommentTargetId(id string) (uint32, error) {
	var comment SubCommentModel
	d := DB.Self.Where("id = ?", id).First(&comment)
	if d.Error != nil || comment.Id == "" {
		return comment.UserId, d.Error
	}

	// The comment target is a parentComment instead of a subComment
	var parentComment ParentCommentModel
	d = DB.Self.Where("id = ?", id).Find(&parentComment)
	return parentComment.UserId, d.Error
}
