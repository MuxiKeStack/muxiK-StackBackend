package model

type ClassTableModel struct {
	Id      uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	UserId  uint32 `gorm:"column:user_id"`
	Courses string `gorm:"column:courses"`
}
