package script

/*package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	//"github.com/lexkong/log"
)

var DB  gorm.DB

func delNull(c string) string{
	if c == ""{
		return "0"
	}else{
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
	default:
		return 9
	}
	return 9
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
	return 0
}

func judge3 (c string) string{
	switch c{
	case "单":
		return "1"
	case "双":
		return "2"
	}
	return "error"
}

func chToNum(a string) string{
	switch a{
	case "一":
		return "1"
	case "二":
		return "2"
	case "三":
		return "3"
	case "四":
		return "4"
	case "五":
		return "5"
	case "六":
		return "6"
	case "日":
		return "7"
	}
	return "0"
}

func analyzeTime(time string) string{
	if time == ""{
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

func preAnalyzeWeek(time string) string{
	if time == ""{
		return ""
	}
	split1 := strings.Index(time, "{")
	split2 := strings.Index(time, "}")
	finstr := analyzeWeek(time[split1+1:split2])
	//fmt.Println(finstr)
	return finstr
}

func analyzeWeek(time string) string{
	split3 := strings.Index(time, "(")
	split4 := strings.Index(time, ",")
	split5 := strings.Index(time, "周")
	var finstr string
	if split4 != -1{
		week := strings.SplitN(time, ",", -1)
		var i int
		for i = 0; i < len(week)-1; i++{
			//fmt.Println(week[0],week[1],len(week))
			finstr = finstr + analyzeManyWeek(week[i]) + ","
		}
		finstr = finstr + analyzeManyWeek(week[i]) + "#0"
		//fmt.Println(finstr)
		return finstr
	}else{
		if split3 != -1{
			finstr = time[0:split3-3] + "#" + judge3(time[split3+1:split3+4])
			//fmt.Println(time[split3+1:split3+4])
		}else{
			finstr = time[:split5] + "#0"
		}
		//fmt.Println(finstr)
	}
	return finstr
}

func analyzeManyWeek(section string) string{
	split1 := strings.Index(section, "周")
	var finstr string
	finstr = section[:split1]
	//fmt.Println(finstr)
	return finstr
}

func analyzeClass(classid string) string{
	split1 := strings.Index(classid, "堂")
	finstr := "(" + classid[split1+4:] + ")"
	return finstr
}

func main(){
	db, err := gorm.Open("mysql", "muxi:123@(127.0.0.1:3306)/muxikstack?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Println(err)
	}else {
		fmt.Println("connection succeed")
	}

	db.SingularTable(true)
	//db.CreateTable(User123{})
	//user := &User123{Sid:"2018212693"}
	//db.Create(user)

	var float float32
	f, err := excelize.OpenFile("./2.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows := f.GetRows("公共课")
	var scourseid string
	for n, row := range rows {
		if n == 0{
			continue
		}
		name := row[1]
		if strings.Contains(name, "大学体育"){
			name = name + analyzeClass(row[3])
		}
		scourseid = row[2]
		teachers := util.GetTeachersSqStrBySplitting(row[8])
		key := util.HashCourseId(scourseid, teachers)
		cred, _ := strconv.ParseFloat(row[4], 32)
		float = float32(cred)
		onecourse:= &model.HistoryCourseModel{
			Hash:      key,
			Name:     name,
			Credit:    float,
			Teacher:  teachers,
			Type:     judge1(row[2][3:4]),
		}
		db.Create(onecourse)
	}

	for i := 6; i <= 9; i++{
		stri := strconv.Itoa(i)
		rows = f.GetRows("201" +  stri + "级")
		for n, row := range rows {
			if n == 0{
				continue
			}
			scourseid = row[2]
			teachers := util.GetTeachersSqStrBySplitting(row[8])
			key := util.HashCourseId(scourseid, teachers)
			cred, _ := strconv.ParseFloat(row[4], 32)
			float = float32(cred)
			onecourse:= &model.HistoryCourseModel{
				Hash:      key,
				Name:     row[1],
				Credit:    float,
				Teacher:  teachers,
				Type:     judge1(row[2][3:4]),
			}
			db.Create(onecourse)
		}
	}
	defer db.Close()
}*/
