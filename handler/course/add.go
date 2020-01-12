package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/util"
	//"regexp"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/gin-gonic/gin"
)

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
}

func judge2(c string) uint8 {
	switch c {
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
		return 0
	}
}

//增加一个课程
func AddCourse(c *gin.Context) {
	var float float32
	f, err := excelize.OpenFile("./1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows := f.GetRows("公共课")
	for _, row := range rows {
		scourseid := row[2]
		teachers := util.GetTeachersSqStrBySplitting(row[8])
		//key := scourseid + row[8]
		key := util.HashCourseId(scourseid, teachers)
		//result :=  hex.EncodeToString(md5.Sum(key))
		//result := hex.EncodeToString(key.Sum(nil))
		//result := md5lnst.Sum([]byte(""))
		cred, _ := strconv.ParseFloat(row[4], 32)
		float = float32(cred)
		clas, _ := strconv.ParseUint(row[3], 10, 64)
		//reg := regexp.MustCompile(`^.{8}`)
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("panic %s\n", err)
			}
		}()
		onecourse := &model.UsingCourseModel{
			Hash:     key,
			Name:     test(row[1]),
			CourseId: test(row[2]),
			ClassId:  clas,
			Credit:   float,
			Teacher:  teachers,
			Type:     judge1(row[2][4:5]),
			Time1:    test(row[10]), //regexp.FindString("^.{8}",row[10]),
			Place1:   test(row[11]),
			Time2:    test(row[12]),
			Place2:   test(row[13]),
			Time3:    test(row[14]),
			Place3:   test(row[15]),
			Weeks1:   test(row[10]),
			Weeks2:   test(row[12]),
			Weeks3:   test(row[14]),
			Region:   judge2(row[11][0:1]),
		}
		//			fmt.Println(onecourse)
		if err := onecourse.Add(); err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}
	}
}

func test(row string) string {
	if row == "" {
		return "nil"
	}
	return row
}
