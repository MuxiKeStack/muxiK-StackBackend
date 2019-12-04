package service

import "github.com/spf13/viper"

var KeywordMap map[string]string

var (
	accessKey    string
	secretKey    string
	bucketName   string
	domainName   string
	upToken      string
	setTimeEpoch int64
	typeMap      map[string]bool
)

func init() {
	KeywordMap = map[string]string{
		// 信息管理学院
		"信技":   "信息技术基础",
		"c语言":  "高级语言程序设计",
		"C语言":  "高级语言程序设计",
		"c":    "高级语言程序设计",
		"C":    "高级语言程序设计",
		"数据库":  "数据库系统原理",
		"计网":   "计算机网络",
		"java": "Java",
		"运筹学":  "管理运筹学",
		"电商概论": "电子商务概论",
		// "ERP" : "ERP原理与应用",
		"erp": "ERP",
		"Erp": "ERP",
		"XML": "XML",
		"xml": "XML",
		"Xml": "XML",
		// web 由于对应Web站点设计与管理, WEB程序设计, 高级WEB程序设计, 不做搜索优化
		"Jsp":  "JSP",
		"jsp":  "JSP",
		"移动电商": "移动电子商务",
		"跨境电商": "跨境电子商务实务",
		"uml":  "UML",
		"Uml":  "UML",
		// 公共课
		"思修":  "思想道德修养与法律基础",
		"马基":  "马克思主义基本原理",
		"马原":  "马克思主义基本原理",
		"近代史": "中国近现代史纲要",
		"毛概":  "毛泽东思想和中国特色社会主义理论体系概论",
		"高数":  "高等数学",
		"线代":  "线性代数",
		"概统":  "概率论与数理统计",
		// 新闻传播学院
		"广告学": "广告学概论",
		// 计算机学院
		"c++":    "C++",
		"python": "Python",
		"linux":  "Linux",
		"unix":   "Unix",
		"微机":     "微机原理与接口技术",
		"计组":     "计算机组成原理",
		// 生命科学学院
		"人体解剖": "人体组织解剖学",
		// 物理科学与技术学院
		"大物": "大学物理",
		// 数学与统计学院
		"数分": "数学分析",
		"高代": "高等代数与解析几何",
		// 经济与工商管理学院
		"微经": "微观经济学",
		"宏经": "宏观经济学",
		// 政治与国际关系学院
		"马哲": "马克思主义哲学原理",
	}
	// 先初始化一些信息
	accessKey = viper.GetString("oss.access_key")
	secretKey = viper.GetString("oss.secret_key")
	bucketName = viper.GetString("oss.bucket_name")
	domainName = viper.GetString("oss.domain_name")
	typeMap = map[string]bool{".jpg": true, ".png": true, ".bmp": true, "jpeg": true, "gif": true}

}