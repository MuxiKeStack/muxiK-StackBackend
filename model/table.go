package model

func (table *ClassTableModel) TableName() string {
	return "class_table"
}

// Create a new table.
func (table *ClassTableModel) New() error {
	return DB.Self.Create(table).Error
}

// Delete the table.
func (table *ClassTableModel) Delete() error {
	return DB.Self.Delete(table).Error
}

// Get table info by id.
func (table *ClassTableModel) Get() error {
	return DB.Self.First(table).Error
}

// Rename the table.
func (table *ClassTableModel) Rename(newName string) error {
	DB.Self.First(table)
	table.Name = newName
	d := DB.Self.Save(table)
	return d.Error
}

// Judge whether the table exists.
func (table *ClassTableModel) Existing() bool {
	d := DB.Self.Where("id = ? AND user_id = ?", table.Id, table.UserId).First(table)
	return !d.RecordNotFound()
}

// Update table's class info.
func (table *ClassTableModel) UpdateClasses(classes string) error {
	table.Classes = classes
	d := DB.Self.Save(table)
	return d.Error
}

// Get table by id.
func GetTableById(id uint32) (*ClassTableModel, error) {
	table := &ClassTableModel{Id: id}
	d := DB.Self.First(table)
	if d.RecordNotFound() {
		return nil, nil
	}
	return table, d.Error
}

// Judge whether the tableId is valid.
func TableIsExisting(tableId uint32, userId uint32) bool {
	var table ClassTableModel
	d := DB.Self.Where("id = ? AND user_id = ?", tableId, userId).First(&table)
	return !d.RecordNotFound()
}

// Get tables by userId.
func GetTablesByUserId(userId uint32) (*[]ClassTableModel, error) {
	var tables []ClassTableModel
	d := DB.Self.Find(&tables, "user_id = ?", userId)
	if d.RecordNotFound() {
		return &tables, nil
	}
	return &tables, d.Error
}

// Get user's table amount by userId.
func GetTableAmount(userId uint32) uint32 {
	var count uint32
	var table ClassTableModel
	DB.Self.Where("user_id = ?", userId).Find(&table).Count(&count)
	return count
}

/*--------- Class Operation -----------*/

func IsClassExisting(hashId string, classId string) bool {
	var class UsingCourseModel
	d := DB.Self.Where("hash = ? AND class_id = ?", hashId, classId).First(&class)
	return !d.RecordNotFound()
}

// 根据hash id 和教学班编号获取课堂
func GetClassByHashId(hashId string, classId string) (*UsingCourseModel, error) {
	var class UsingCourseModel
	d := DB.Self.Where("hash = ? AND class_id = ?", hashId, classId).First(&class)
	return &class, d.Error
}

func GetCourseHashIdById(id string) (string, error) {
	var course HistoryCourseModel
	d := DB.Self.Where("id = ?", id).First(&course)
	return course.Hash, d.Error
}
