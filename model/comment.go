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
	return DB.Self.Create(comment).Error
}

// Delete the parent comment.
func (comment *ParentCommentModel) Delete() error {
	return DB.Self.Delete(comment).Error
}

// Get a parent comment by its id.
func (comment *ParentCommentModel) GetById() error {
	d := DB.Self.Unscoped().First(comment, "id = ?", comment.Id)
	return d.Error
}

// Get parent comments by evaluation id.
func GetParentComments(EvaluationId uint32, limit, offset int32) (*[]ParentCommentModel, error) {
	var comments []ParentCommentModel

	d := DB.Self.Unscoped().Where("evaluation_id = ?", EvaluationId).
		Order("time").Limit(limit).Offset(offset).Find(&comments)

	if d.RecordNotFound() {
		return &comments, nil
	}
	return &comments, d.Error
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
	return DB.Self.Create(comment).Error
}

// Delete the subComment.
func (comment *SubCommentModel) Delete() error {
	return DB.Self.Delete(comment).Error
}

// Get a subComment by its id.
func (comment *SubCommentModel) GetById() error {
	d := DB.Self.Unscoped().First(comment, "id = ?", comment.Id)
	return d.Error
}

// Get subComments by their parentId.
func GetSubCommentsByParentId(ParentId string) (*[]SubCommentModel, error) {
	var subComments []SubCommentModel
	d := DB.Self.Unscoped().Where("parent_id = ?", ParentId).Order("time").Find(&subComments)
	return &subComments, d.Error
}

// Judge whether it is a subComment by id, if so then also return the subComment.
func IsSubComment(id string) (*SubCommentModel, bool) {
	var comment SubCommentModel
	d := DB.Self.Where("id = ?", id).First(&comment)
	if !d.RecordNotFound() {
		return &comment, true
	}
	return nil, false
}

// Judge whether is a parent comment.
func IsParentComment(id string) bool {
	var comment ParentCommentModel
	d := DB.Self.Where("id = ?", id).First(&comment)
	return !d.RecordNotFound()
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

// Cancel liking a comment by the like-record id.
func CommentCancelLiking(id uint32) error {
	var data = CommentLikeModel{Id: id}
	d := DB.Self.Delete(&data)
	return d.Error
}

// Judge whether a comment has already liked by the current user,
// return like-record id and bool type.
func CommentHasLiked(userId uint32, commentId string) (uint32, bool) {
	var data CommentLikeModel
	d := DB.Self.Where("user_id = ? AND comment_id = ?", userId, commentId).Find(&data)
	return data.Id, !d.RecordNotFound()
}

// Get comment's total like amount by commentId.
func GetCommentLikeSum(commentId string) (count uint32) {
	var data CommentLikeModel
	DB.Self.Where("comment_id = ?", commentId).Find(&data).Count(&count)
	return
}
