package service

import (
	"fmt"
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/lexkong/log"
)

// 爬取成绩
func NewGradeRecord(userId uint32, sid, pwd string) error {
	data, err := util.GetGradeFromXK(sid, pwd)
	if err != nil {
		log.Error("util.GetGradeFromXK function error", err)
		return err
	}
	fmt.Println(data)

	// TO DO: 并发
	for _, item := range *data {
		teacher := strings.ReplaceAll(item.Teacher, ";", ",")
		hash := util.HashCourseId(item.CourseId, teacher)

		if ok, err := model.GradeRecordExisting(userId, hash); err != nil {
			log.Error("GradeRecordExisting function error", err)
			return err
		} else if ok {
			//log.Info("The record has existed")
			continue
		}

		g := &model.GradeModel{
			UserId:       userId,
			CourseHashId: hash,
			CourseName:   item.CourseName,
			TotalScore:   item.TotalScore,
			UsualScore:   item.UsualScore,
			FinalScore:   item.FinalScore,
			HasAdded:     false,
		}
		if err := g.New(); err != nil {
			log.Error("Add new grade record error", err)
			return err
		}
	}
	return nil
}

// 将成绩导入到各课程的成绩统计中
func NewGradeSampleFoCourses(userId uint32) error {
	// 获取未添加的成绩数据
	records, err := model.GetGradeRecordsByUserId(userId)
	if err != nil {
		log.Error("GetGradeRecordsByUserId function error", err)
		return err
	}

	// TO DO: 并发优化
	for _, record := range *records {
		if err := NewGradeDataAdditionForOneCourse(userId, &record); err != nil {
			log.Error("NewGradeDataAdditionForOneCourse function error", err)
			return err
		}
	}
	return nil
}

// 一门课程的成绩样本数据添加
func NewGradeDataAdditionForOneCourse(userId uint32, data *model.GradeModel) error {
	// 获取课程
	course, err := model.GetHistoryCourseByHashId(data.CourseHashId)
	if err != nil {
		log.Error("GetHistoryCourseByHashId function error", err)
		return err
	}

	// 获取成绩数据
	//data, _, err := model.GetGradeRecord(userId, hashId)
	//if err != nil {
	//	log.Error("GetGradeRecord function error", err)
	//	return err
	//}

	// 有数据库写入覆盖的隐患
	curSampleSize := course.GradeSampleSize
	course.TotalGrade = (course.TotalGrade*float32(curSampleSize) + data.TotalScore) / float32(curSampleSize+1)
	course.UsualGrade += (course.UsualGrade*float32(curSampleSize) + data.UsualScore) / float32(curSampleSize+1)
	course.GradeSampleSize++

	if data.TotalScore > 85 {
		course.GradeSection1++
	} else if data.TotalScore >= 70 {
		course.GradeSection2++
	} else {
		course.GradeSection3++
	}

	if err := course.UpdateGradeInfo(); err != nil {
		log.Error("UpdateGradeInfo function error", err)
		return err
	}

	data.HasAdded = true
	if err := data.Update(); err != nil {
		log.Error("gradeModel.Update function error", err)
		return err
	}

	return nil
}

func GradeImportService(userId uint32, sid, pwd string) error {
	log.Info("Grade import service is called")

	// 获取成绩
	if err := NewGradeRecord(userId, sid, pwd); err != nil {
		log.Error("NewGradeRecord function error", err)
		return err
	}
	// 导入成绩样本数据
	if err := NewGradeSampleFoCourses(userId); err != nil {
		log.Error("NewGradeSampleFoCourses function error", err)
		return err
	}
	return nil
}
