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

func (grade *GradeModel) New() error {
	return DB.Self.Create(grade).Error
}

func (grade *GradeModel) Get() error {
	return DB.Self.Create(grade).Error
}

func (grade *GradeModel) Update() error {
	return DB.Self.Save(grade).Error
}

func GetGradeRecord(userId uint32, hashId string) (*GradeModel, bool, error) {
	var g GradeModel
	d := DB.Self.Where("user_id = ? AND course_hash_id = ?", userId, hashId).First(&g)
	if d.RecordNotFound() {
		return nil, false, nil
	}
	return &g, true, d.Error
}

func GradeRecordExisting(userId uint32, hashId string) (bool, error) {
	var record GradeModel
	d := DB.Self.Where("user_id = ? AND course_hash_id = ?", userId, hashId).First(&record)
	if d.RecordNotFound() {
		return false, nil
	}
	return true, d.Error
}

func GetRecordsNum(userId uint32) (int, error) {
	var count int
	d := DB.Self.Table("grade").Where("user_id = ?", userId).Count(&count)
	if d.RecordNotFound() {
		return 0, nil
	}
	return count, d.Error
}

func GradeRecordUserHasGotten(userId uint32) (bool, error) {
	d := DB.Self.Where("user_id = ?", userId).First(GradeModel{})
	if d.RecordNotFound() {
		return false, nil
	}
	return true, d.Error
}

func GetGradeRecordsByUserId(userId uint32) (*[]GradeModel, error) {
	var records []GradeModel
	d := DB.Self.Where("user_id = ? AND has_added = 0", userId).Find(&records)
	return &records, d.Error
}
