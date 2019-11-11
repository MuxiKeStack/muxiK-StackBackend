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

// Get all tag names.
func GetTagNames() ([]string, error) {
	var tags []TagModel
	var names []string
	d := DB.Self.Find(&tags)

	for _, tag := range tags {
		names = append(names, tag.TagName)
	}

	return names, d.Error
}
