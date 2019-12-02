package model

func (comment *ParentCommentModel) TableName() string {
	return "parent_comment"
}

func (comment *SubCommentModel) TableName() string {
	return "sub_comment"
}

func (data *CommentLikeModel) TableName() string {
	return "comment_like"
}

/*---------------------------- Parent Comment Operation --------------------------*/

// Create a new parent comment.
func (comment *ParentCommentModel) New() error {
	d := DB.Self.Create(comment)
	return d.Error
}

// Get a parent comment by its id.
func (comment *ParentCommentModel) GetById() error {
	d := DB.Self.First(comment, "id = ?", comment.Id)
	return d.Error
}

// Get parent comments by evaluation id.
func GetParentComments(EvaluationId uint32, limit, offset int32) (*[]ParentCommentModel, uint32, error) {
	var count uint32
	var comments []ParentCommentModel

	d := DB.Self.Where("evaluation_id = ?", EvaluationId).
		Order("time").Limit(limit).Offset(offset).Find(&comments).Count(&count)

	return &comments, count, d.Error
}

// Update parentComment's the total number of subComment
func (comment *ParentCommentModel) UpdateSubCommentNum(n int) error {
	num := int(comment.SubCommentNum)
	if num == 0 && n == -1 {
		return nil
	}
	num += n
	d := DB.Self.Model(comment).Update("sub_comment_num", num)
	return d.Error
}

/*---------------------------- SubComment Operation --------------------------*/

// Create a new subComment.
func (comment *SubCommentModel) New() error {
	d := DB.Self.Create(comment)
	return d.Error
}

// Get a subComment by its id.
func (comment *SubCommentModel) GetById() error {
	d := DB.Self.First(comment, "id = ?", comment.Id)
	return d.Error
}

// Get subComments by their parentId.
func GetSubCommentsByParentId(ParentId string) (*[]SubCommentModel, error) {
	var subComments []SubCommentModel
	d := DB.Self.Where("parent_id = ?", ParentId).Order("time").Find(&subComments)
	return &subComments, d.Error
}

// Judge whether it is a subComment by id, if so then also return the subComment.
func IsSubComment(id string) (*SubCommentModel, bool) {
	var comment SubCommentModel
	DB.Self.Where("id = ?", id).First(&comment)
	if comment.Id != "" {
		return &comment, true
	}
	return nil, false
}

/*---------------------------- Other Comment Operations --------------------------*/

// Like a comment by the current user.
func CommentLiking(userId uint32, commentId string) error {
	var data = &CommentLikeModel{
		UserId:    userId,
		CommentId: commentId,
	}
	d := DB.Self.Create(data)
	return d.Error
}

// Cancel liking a comment by the current user.
func CommentCancelLiking(userId uint32, commentId string) error {
	var data = &CommentLikeModel{
		UserId:    userId,
		CommentId: commentId,
	}
	d := DB.Self.Delete(data)
	return d.Error
}

// Judge whether a comment has already liked by the current user.
func CommentHasLiked(userId uint32, commentId string) bool {
	var data CommentLikeModel
	var count int
	DB.Self.Where("user_id = ? AND comment_id = ?", userId, commentId).Find(&data).Count(&count)
	return count > 0
}

// Get comment's total like account by commentId.
func GetCommentLikeSum(commentId string) (count uint32) {
	var data CommentLikeModel
	DB.Self.Where("comment_id = ?", commentId).Find(&data).Count(&count)
	return
}

// Get parentId by commentTargetId.
//func GetParentIdByCommentTargetId(id string) (string, bool) {
//	var comment SubCommentModel
//	DB.Self.Where("id = ?", id).First(&comment)
//	if comment.ParentId != "" {
//		return comment.ParentId, true
//	}
//
//	// The comment target is a parentComment instead of a subComment
//	var parentComment ParentCommentModel
//	DB.Self.Where("id = ?", id).First(&parentComment)
//	if parentComment.Id != "" {
//		return parentComment.Id, true
//	}
//	return "", false
//}

// Get comment target user's id by the commentTargetId.
//func GetTargetUserIdByCommentTargetId(id string) (uint32, bool) {
//	var comment SubCommentModel
//	DB.Self.Where("id = ?", id).First(&comment)
//	if comment.Id != "" {
//		return comment.UserId, true
//	}
//
//	// The comment target is a parentComment instead of a subComment
//	var parentComment ParentCommentModel
//	DB.Self.Where("id = ?", id).First(&parentComment)
//	if parentComment.Id != "" {
//		return parentComment.UserId, true
//	}
//	return 0, false
//}
