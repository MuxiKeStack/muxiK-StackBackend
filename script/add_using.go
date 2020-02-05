package script

/*
package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var (
	DB     gorm.DB
	DBAddr string
	DBUser string
	DBPwd  string
	fileFg = flag.String("file", "sample.xlsx", "set using-course manual excel file (*.xlsx)")
)

func init() {
	flag.Parse()

	// 配置环境变量
	// export MUXIKSTACK_DB_ADDR=127.0.0.1:3306
	// export MUXIKSTACK_DB_USERNAME=muxi
	// export MUXIKSTACK_DB_PASSWORD=muxi
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MUXIKSTACK")
	DBAddr = viper.GetString("DB_ADDR")
	DBUser = viper.GetString("DB_USERNAME")
	DBPwd = viper.GetString("DB_PASSWORD")
}

func delNull(c string) string {
	if c == "" {
		return "0"
	} else {
		return c[0:1]
	}
}

func judge1(c string) uint8 {
	switch c {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "5":
		return 5
	}
	return 4
}

func judge2(c string) uint8 {
	switch c {
	case "0":
		return 0
	case "6":
		return 1
	case "7":
		return 1
	case "8":
		return 2
	case "9":
		return 1
	case "y":
		return 2
	default:
		return 9
	}
}

func judge3(c string) string {
	switch c {
	case "单":
		return "1"
	case "双":
		return "2"
	}
	return "error"
}

func chToNum(a string) string {
	dayMap := map[string]string{"一": "1", "二": "2", "三": "3", "四": "4", "五": "5", "六": "6", "日": "7"}
	if day, ok := dayMap[a]; ok {
		return day
	}
	return "error"

	//switch a {
	//case "一":
	//	return "1"
	//case "二":
	//	return "2"
	//case "三":
	//	return "3"
	//case "四":
	//	return "4"
	//case "五":
	//	return "5"
	//case "六":
	//	return "6"
	//case "日":
	//	return "7"
	//}
	//return "error"
}

func analyzeTime(time string) string {
	if time == "" {
		return ""
	}
	split1 := strings.Index(time, "第")
	//fmt.Println(split1)
	split2 := strings.Index(time, "节")
	//fmt.Println(split2)
	finstr := time[split1+3:split2] + "#" + chToNum(time[split1-3:split1])
	//fmt.Println(finstr)
	return finstr
}

func preAnalyzeWeek(time string) string {
	if time == "" {
		return ""
	}
	split1 := strings.Index(time, "{")
	split2 := strings.Index(time, "}")
	finstr := analyzeWeek(time[split1+1 : split2])
	//fmt.Println(finstr)
	return finstr
}

func analyzeWeek(time string) string {
	split3 := strings.Index(time, "(")
	split4 := strings.Index(time, ",")
	split5 := strings.Index(time, "周")
	var finstr string
	if split4 != -1 {
		week := strings.SplitN(time, ",", -1)
		var i int
		for i = 0; i < len(week)-1; i++ {
			//fmt.Println(week[0],week[1],len(week))
			finstr = finstr + analyzeManyWeek(week[i]) + ","
		}
		finstr = finstr + analyzeManyWeek(week[i]) + "#0"
		//fmt.Println(finstr)
		return finstr
	} else {
		if split3 != -1 {
			finstr = time[0:split3-3] + "#" + judge3(time[split3+1:split3+4])
			//fmt.Println(time[split3+1:split3+4])
		} else {
			finstr = time[:split5] + "#0"
		}
		//fmt.Println(finstr)
	}
	return finstr
}

func analyzeManyWeek(section string) string {
	split1 := strings.Index(section, "周")
	var finStr string
	finStr = section[:split1]
	//fmt.Println(finStr)
	return finStr
}

func analyzeClass(classid string) string {
	split1 := strings.Index(classid, "堂")
	finStr := "(" + classid[split1+4:] + ")"
	return finStr
}

func main() {
	if DBAddr == "" || DBUser == "" {
		fmt.Println("Database config error, required env settings")
		return
	}
	dbOpenCmd := fmt.Sprintf("%s:%s@(%s)/muxikstack?charset=utf8&parseTime=True", DBUser, DBPwd, DBAddr)
	db, err := gorm.Open("mysql", dbOpenCmd)
	//db, err := gorm.Open("mysql", "*:*@(*.*.*.*:*)/muxikstack?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	fmt.Println("connection succeed")

	db.SingularTable(true)
	//db.CreateTable(User123{})
	//user := &User123{Sid:"2018212693"}
	//db.Create(user)

	//f, err := excelize.OpenFile("./2.xlsx")
	f, err := excelize.OpenFile(*fileFg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Open excel file successfully")
	fmt.Println("Start importing...")

	rows := f.GetRows("公共课")
	var sCourseId string
	var float float32
	for n, row := range rows {
		if n == 0 {
			continue
		}
		name := row[1]
		if strings.Contains(name, "大学体育") {
			name = name + analyzeClass(row[3])
		}
		sCourseId = row[2]
		teachers := util.GetTeachersSqStrBySplitting(row[8])
		key := util.HashCourseId(sCourseId, teachers)
		cred, _ := strconv.ParseFloat(row[4], 32)
		float = float32(cred)
		oneCourse := &model.UsingCourseModel{
			Hash:     key,
			Academy:  row[0],
			Name:     name,
			CourseId: row[2],
			ClassId:  row[3],
			Credit:   float,
			Teacher:  teachers,
			Type:     judge1(row[2][4:5]),
			Time1:    analyzeTime(row[10]),
			Place1:   row[11],
			Time2:    analyzeTime(row[12]),
			Place2:   row[13],
			Time3:    analyzeTime(row[14]),
			Place3:   row[15],
			Weeks1:   preAnalyzeWeek(row[10]),
			Weeks2:   preAnalyzeWeek(row[12]),
			Weeks3:   preAnalyzeWeek(row[14]),
			Region:   judge2(delNull(row[11])),
		}
		db.Create(oneCourse)

		fmt.Printf("正在导入第  %d  条记录...\r", oneCourse.Id)
	}

	for i := 6; i <= 9; i++ {
		rows = f.GetRows("201" + strconv.Itoa(i) + "级")
		for n, row := range rows {
			if n == 0 {
				continue
			}
			sCourseId = row[2]
			teachers := util.GetTeachersSqStrBySplitting(row[8])
			key := util.HashCourseId(sCourseId, teachers)
			cred, _ := strconv.ParseFloat(row[4], 32)
			float = float32(cred)
			oneCourse := &model.UsingCourseModel{
				Hash:     key,
				Academy:  row[0],
				Name:     row[1],
				CourseId: row[2],
				ClassId:  row[3],
				Credit:   float,
				Teacher:  teachers,
				Type:     judge1(row[2][4:5]),
				Time1:    analyzeTime(row[10]),
				Place1:   row[11],
				Time2:    analyzeTime(row[12]),
				Place2:   row[13],
				Time3:    analyzeTime(row[14]),
				Place3:   row[15],
				Weeks1:   preAnalyzeWeek(row[10]),
				Weeks2:   preAnalyzeWeek(row[12]),
				Weeks3:   preAnalyzeWeek(row[14]),
				Region:   judge2(delNull(row[11])),
			}
			db.Create(oneCourse)

			fmt.Printf("正在导入第  %d  条记录...\r", oneCourse.Id)
		}
	}
	fmt.Println("Import has completed")
}*/
