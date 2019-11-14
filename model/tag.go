package model

// Get tag name by id.
func GetTagNameById(id int) string {
	var tag TagModel
	DB.Self.Where("id = ?", id).First(&tag)
	return tag.TagName
}

// Get tag total number.
func GetTagSum() int {
	var count int
	var tag TagModel
	DB.Self.Find(&tag).Count(&count)
	return count
}

// Get all tags including ids and names.
func GetTags() ([]TagModel, error) {
	var tags []TagModel
	d := DB.Self.Find(&tags)
	return tags, d.Error
}
