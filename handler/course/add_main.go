package course
/*
package main

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
	return 0
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

//增加一个课程
// func AddCourse() {
// 	var float float32
// 	f, err := excelize.OpenFile("./1.xlsx")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	rows := f.GetRows("公共课")
// 	var scourseid string
// 	var clas uint64
// 	for n, row := range rows {
// 		if n == 0{
// 			continue
// 		}
// 		scourseid = row[2]
// 		teachers := util.GetTeachersSqStrBySplitting(row[8])
// 		key := util.HashCourseId(scourseid, teachers)
// 		cred, _ := strconv.ParseFloat(row[4], 32)
// 		float = float32(cred)
// 		clas, _ = strconv.ParseUint(row[3], 10, 64)
// 		onecourse:= &model.UsingCourseModel{
// 			Hash:      key,
// 			Name:     test(row[1]),
// 			CourseId: test(row[2]),
// 			ClassId:  clas,
// 			Credit:   float,
// 			Teacher:  teachers,
// 			Type:     judge1(row[2][4:5]),
// 			Time1:    test(row[10]),
// 			Place1:   test(row[11]),
// 			Time2:    test(row[12]),
// 			Place2:   test(row[13]),
// 			Time3:    test(row[14]),
// 			Place3:   test(row[15]),
// 			Weeks1:   test(row[10]),
// 			Weeks2:   test(row[12]),
// 			Weeks3:   test(row[14]),
// 			//Region:   judge2(row[11][0:1]),
// 		}
// 		DB.Create(onecourse)
// 		// fmt.Println(onecourse)
// 		// if err := onecourse.Add(); err != nil {
// 		// 	log.Info("add onecourse error")
// 		// 	return
// 		//  }
// 	}
	
// }

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

func analyzeWeek(time string) string{
	if time == ""{
		return ""
	}
	split1 := strings.Index(time, "{")
	split2 := strings.Index(time, "}")
	split3 := strings.Index(time, "(")
	var finstr string
	if split3 != -1{
		finstr = time[split1+1:split2-8] + "#" + judge3(time[split3+1:split3+4])
		fmt.Println(time[split3+1:split3+4])
	}else{
		finstr = time[split1+1:split2-3] + "#0"
	}
	fmt.Println(finstr)
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
	f, err := excelize.OpenFile("./1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows := f.GetRows("公共课")
	var scourseid string
	var clas uint64
	for n, row := range rows {
		if n == 0{
			continue
		}
		scourseid = row[2]
		teachers := util.GetTeachersSqStrBySplitting(row[8])
		key := util.HashCourseId(scourseid, teachers)
		cred, _ := strconv.ParseFloat(row[4], 32)
		float = float32(cred)
		clas, _ = strconv.ParseUint(row[3], 10, 64)
		onecourse:= &model.UsingCourseModel{
			Hash:      key,
			Name:     row[1],
			CourseId: row[2],
			ClassId:  clas,
			Credit:    float,
			Teacher:  teachers,
			Type:     judge1(row[2][4:5]),
			Time1:    analyzeTime(row[10]),
			Place1:   row[11],
			Time2:    analyzeTime(row[12]),
			Place2:  row[13],
			Time3:   analyzeTime(row[14]),
			Place3:   row[15],
			Weeks1:   analyzeWeek(row[10]),
			Weeks2:   analyzeWeek(row[12]),
			Weeks3:   analyzeWeek(row[14]),
			Region:   judge2(delNull(row[11])),
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
			scourseid = row[4]
			teachers := util.GetTeachersSqStrBySplitting(row[10])
			key := util.HashCourseId(scourseid, teachers)
			cred, _ := strconv.ParseFloat(row[6], 32)
			float = float32(cred)
			clas, _ = strconv.ParseUint(row[5], 10, 64)
			onecourse:= &model.UsingCourseModel{
				Hash:      key,
				Name:     row[2],
				CourseId: row[4],
				ClassId:  clas,
				Credit:    float,
				Teacher:  teachers,
				Type:     judge1(row[4][4:5]),
				Time1:    analyzeTime(row[12]),
				Place1:   row[13],
				Time2:    analyzeTime(row[14]),
				Place2:  row[15],
				Time3:   analyzeTime(row[16]),
				Place3:   row[17],
				Weeks1:   analyzeWeek(row[12]),
				Weeks2:   analyzeWeek(row[14]),
				Weeks3:   analyzeWeek(row[16]),
				Region:   judge2(delNull(row[13])),
			}
			db.Create(onecourse)
		}
	}
	defer db.Close()
}*/