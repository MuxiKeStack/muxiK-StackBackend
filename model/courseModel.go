package model

//云课堂课程物理表
type HistoryCourseModel struct {
	Id              uint32  `gorm:"column:id; primary_key"`
	Hash            string  `gorm:"column:hash; unique_key"`  //教师名和课程hash成的唯一标识
	Name            string  `gorm:"column:name"`              //课程名称
	Teacher         string  `gorm:"column:teacher"`           //教师性名
	CourseId        string  `gorm:"column:course_id"`         //UI上需要展示
	Type            uint8   `gorm:"column:type"`              //课程类型，公共课为0，专业课为1
	Rate            float32 `gorm:"column:rate"`              //课程评价星级
	StarsNum        uint32  `gorm:"column:stars_num"`         //参与评分人数
	GradeSampleSize uint32  `gorm:"column:grade_sample_size"` // 成绩样本数
	TotalGrade      float32 `gorm:"column:total_grade"`       // 总成绩均分
	UsualGrade      float32 `gorm:"column:usual_grade"`       // 平时成绩均分
	GradeSection1   uint32  `gorm:"column:grade_section_1"`   // 成绩区间，85以上人数
	GradeSection2   uint32  `gorm:"column:grade_section_2"`   // 成绩区间，70-85人数
	GradeSection3   uint32  `gorm:"column:grade_section_3"`   // 成绩区间，70一下人数
}

//选课手册课堂物理表
type UsingCourseModel struct {
	Id       uint32  `gorm:"column:id; primary_key"`
	Hash     string  `gorm:"column:hash; unique_key"` //教师名和课程hash成的唯一标识
	Name     string  `gorm:"column:name"`             //课程名称
	Academy  string  `gorm:"column:academy"`          //开课学院
	Credit   float32 `gorm:"column:credit"`           //学分
	Teacher  string  `gorm:"column:teacher"`          //教师姓名
	CourseId string  `gorm:"column:course_id"`        //UI上需要展示
	ClassId  string  `gorm:"column:class_id"`         //课堂编号，用于区分课堂
	Type     uint8   `gorm:"column:type"`             //通识必修，通识选修，通识核心，专业必修，专业选修分别为 0/1/2/3/4
	Time1    string  `gorm:"column:time1"`            //上课时间1
	Place1   string  `gorm:"column:place1"`           //上课地点1
	Time2    string  `gorm:"column:time2"`            //上课时间2
	Place2   string  `gorm:"column:place2"`           //上课地点2
	Time3    string  `gorm:"column:time3"`            //上课时间3
	Place3   string  `gorm:"column:place3"`           //上课地点3
	Weeks1   string  `gorm:"column:weeks1"`           //上课周数
	Weeks2   string  `gorm:"column:weeks2"`           //上课周数
	Weeks3   string  `gorm:"column:weeks3"`           //上课周数
	Region   uint8   `gorm:"column:region"`           //上课地区，暂定：东区，西区，其他。加索引（筛选条件）
}

// 添加了搜索时需要返回的字段
type UsingCourseSearchModel struct {
	UsingCourseModel
	Rate     float32 `gorm:"column:rate"`      //课程评价星级
	StarsNum uint32  `gorm:"column:stars_num"` //参与评分人数
}

// time格式：1-2#1 ==> 周一的第一到第二节，#后面的数字代表周几(1-7)
// week格式：2-17#0 ==> 2-17周，全周；0为全周，1为单周，2为双周

// 个人教学课程物理表
type SelfCourseModel struct {
	Id      uint32 `gorm:"column:id"`
	UserId  uint32 `gorm:"column:user_id"`
	Num     uint32 `gorm:"column:num"`
	Courses string `gorm:"column:courses"`
}

type HistoryCourseInfo struct {
	Hash     string
	CourseId string
	Name     string
	Teacher  string
	Type     uint8
	Rate     float32
	StarsNum uint32
}

type UsingCourseInfo struct {
	CourseId string
	ClassId  string
	Name     string
	Academy  string
	Credit   float32
	Teacher  string
	Type     uint8
	Time1    string
	Place1   string
	Time2    string
	Place2   string
	Time3    string
	Place3   string
	Weeks1   string
	Weeks2   string
	Weeks3   string
	Region   uint8
}
