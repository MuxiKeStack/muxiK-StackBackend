package model

func (table *ClassTableModel) TableName() string {
	return "class_table"
}

func (table *ClassTableModel) New() error {
	d := DB.Self.Create(table)
	return d.Error
}

func (table *ClassTableModel) Delete() error {
	d := DB.Self.Delete(table)
	return d.Error
}

func (table *ClassTableModel) GetById() error {
	d := DB.Self.First(table, "id = ?", table.Id)
	//d := DB.Self.First(table)
	return d.Error
}

func (table *ClassTableModel) Rename(newName string) error {
	DB.Self.First(table)
	table.Name = newName
	d := DB.Self.Save(table)
	return d.Error
}

func (table *ClassTableModel) Existing() bool {
	var count int
	DB.Self.First(table).Count(&count)
	return count == 1
}

func (table *ClassTableModel) UpdateClasses(classes string) error {
	table.Classes = classes
	d := DB.Self.Save(table)
	return d.Error
}

func GetTablesByUserId(userId uint32) (*[]ClassTableModel, error) {
	var tables []ClassTableModel
	d := DB.Self.Find(&tables, "user_id = ?", userId)
	return &tables, d.Error
}

/*--------- Class Operation -----------*/

func IsClassExisting(id string) bool {
	var class UsingCourseModel
	var count int8
	DB.Self.Where("hash = ?", id).First(&class).Count(&count)
	return count == 1
}

func GetClassByHashId(id string) (*UsingCourseModel, error) {
	var class UsingCourseModel
	d := DB.Self.Where("hash = ?", id).First(&class)
	return &class, d.Error
}

func GetCourseHashIdById(id string) (string, error) {
	var course HistoryCourseModel
	d := DB.Self.Where("id = ?", id).First(&course)
	return course.Hash, d.Error
}
