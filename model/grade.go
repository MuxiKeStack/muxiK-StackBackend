package model

type GradeModel struct {
	Id           uint32  `gorm:"column:id"`
	UserId       uint32  `gorm:"column:user_id"`
	CourseHashId string  `gorm:"column:course_hash_id"`
	CourseName   string  `gorm:"column:course_name"`
	TotalScore   float32 `gorm:"column:total_score"`
	UsualScore   float32 `gorm:"column:usual_score"`
	FinalScore   float32 `gorm:"column:final_score"`
	HasAdded     bool    `gorm:"column:has_added"`
}

func (*GradeModel) TableName() string {
	return "grade"
}

// Create new grade record.
func (grade *GradeModel) New() error {
	return DB.Self.Create(grade).Error
}

// Get grade record by id.
func (grade *GradeModel) Get() error {
	return DB.Self.Create(grade).Error
}

// Update grade record.
func (grade *GradeModel) Update() error {
	return DB.Self.Save(grade).Error
}

// Get grade record by userId and courseHashId.
func GetGradeRecord(userId uint32, hashId string) (*GradeModel, bool, error) {
	var g GradeModel
	d := DB.Self.Where("user_id = ? AND course_hash_id = ?", userId, hashId).First(&g)
	if d.RecordNotFound() {
		return nil, false, nil
	}
	return &g, true, d.Error
}

// Check whether the grade record exists by userId and courseHashId.
func GradeRecordExisting(userId uint32, hashId string) (bool, error) {
	var count int
	d := DB.Self.Table("grade").Where("user_id = ? AND course_hash_id = ?", userId, hashId).Count(&count)
	//d := DB.Self.Exec("SELECT id FROM grade WHERE user_id = ? AND course_hash_id = ?;", userId, hashId)
	//d := DB.Self.Where("user_id = ? AND course_hash_id = ?", userId, hashId).First(&GradeModel{})
	return count != 0, d.Error
}

// Get the amount of grade records by userId.
func GetRecordsNum(userId uint32) (int, error) {
	var count int
	d := DB.Self.Table("grade").Where("user_id = ?", userId).Count(&count)
	return count, d.Error
}

// Get user's all grade records that have not added to statistical sample.
func GetGradeRecordsByUserId(userId uint32) (*[]GradeModel, error) {
	var records []GradeModel
	d := DB.Self.Where("user_id = ? AND has_added = 0", userId).Find(&records)
	return &records, d.Error
}
