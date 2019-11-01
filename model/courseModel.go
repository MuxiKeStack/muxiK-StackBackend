package model

//云课堂课程物理表
type HistoryCourseModel struct {
	Id            uint32  `gorm:"column:id; primary_key"`
	Hash          string  `gorm:"column:hash; unique_key"`  //教师名和课程hash成的唯一标识
	Name          string  `gorm:"column:name"`              //课程名称
	Teacher       string  `gorm:"column:teacher"`           //教师性名
	Type          uint8   `gorm:"column:type"`              //课程类型，公共课为0，专业课为1
	Rate          float32 `gorm:"column:rate"`              //课程评价星级
	StarsNum      uint32  `gorm:"column:stars_num"`         //参与评分人数
	TotalScore    float32 `gorm:"column:total_score"`       //总评均分
	OrdinaryScore float32 `gorm:"column:ordinary_score"`    //平时均分
}

//选课手册课堂物理表
type UsingCourseModel struct {
	Id            uint32  `gorm:"column:id; primary_key"`
	Hash          string  `gorm:"column:hash; unique_key"`  //教师名和课程hash成的唯一标识
	Name          string  `gorm:"column:name"`              //课程名称
	Credit        float32 `gorm:"column:credit"`            //学分
	Teacher       string  `gorm:"column:teacher"`           //教师姓名
	CourseId      uint64  `gorm:"column:course_id"`         //UI上需要展示
	ClassId       uint8   `gorm:"column:class_id"`          //课堂编号，用于区分课堂
	Type          uint8   `gorm:"column:type"`              //通识必修，通识选修，通识核心，专业必修，专业选修分别为 0/1/2/3/4
	CreditType    uint8   `gorm:"column:credit_type"`       //学分类别，文科理科艺术之类的，加索引（筛选条件）
	Time1         string  `gorm:"column:time1"`             //上课时间1
	Place1        string  `gorm:"column:place1"`            //上课地点1
	Time2         string  `gorm:"column:time2"`             //上课时间2
	Place2        string  `gorm:"column:place2"`            //上课地点2
	Time3         string  `gorm:"column:time3"`             //上课时间3
	Place3        string  `gorm:"column:place3"`            //上课地点3
	Weeks         string  `gorm:"column:weeks"`             //上课周数
	Region        uint8   `gorm:"column:region"`            //上课地区，南湖，东区，西区。加索引（筛选条件）
}

type HistoryCourseInfo struct {
	CourseId      string
	Name          string
	Teacher       string
	Type          uint8
	Rate          float32
	StarsNum      uint32
	TotalScore    float32
	OrdinaryScore float32
}

type UsingCourseInfo struct {
	CourseId      string
	ClassId       string
	Name          string
	Credit        float32
	Teacher       string
	Type          uint8
	CreditType    uint8
	Time1         string
	Place1        string
	Time2         string
	Place2        string
	Time3         string
	Place3        string
	Weeks         string
	Region        uint8
}
