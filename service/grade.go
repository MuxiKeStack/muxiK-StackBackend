package service

import (
	"strconv"
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
	//fmt.Println(data)

	for _, item := range data.Items {
		teacher := strings.ReplaceAll(item.Jsxm, ";", ",")
		hash := util.HashCourseId(item.Kch, teacher)

		if ok, err := model.GradeRecordExisting(userId, hash); err != nil {
			log.Error("GradeRecordExisting function error", err)
			return err
		} else if ok {
			//log.Info("The record has existed")
			continue
		}

		totalScore, err := strconv.ParseFloat(item.Cj, 32)
		if err != nil {
			log.Error("Parse totalScore error", err)
			return err
		}
		//fmt.Println(teacher, hash, float32(totalScore))

		g := &model.GradeModel{
			UserId:         userId,
			CourseHashId:   hash,
			CourseName:     item.Kcmc,
			TotalScore:     float32(totalScore),
			UsualScore:     0,
			FinalExamScore: 0,
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
	records, err := model.GetGradeRecordsByUserId(userId)
	if err != nil {
		log.Error("GetGradeRecordsByUserId function error", err)
		return err
	}

	// TO Do: 并发优化
	for _, record := range *records {
		if err := NewGradeDataAdditionForOneCourse(userId, record.CourseHashId); err != nil {
			log.Error("NewGradeDataAdditionForOneCourse function error", err)
			return err
		}
	}
	return nil
}

// 一门课程的成绩样本数据添加
func NewGradeDataAdditionForOneCourse(userId uint32, hashId string) error {
	course, err := model.GetHistoryCourseByHashId(hashId)
	if err != nil {
		log.Error("GetHistoryCourseByHashId function error", err)
		return err
	}

	data, _, err := model.GetGradeRecord(userId, hashId)
	if err != nil {
		log.Error("GetGradeRecord function error", err)
		return err
	}

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
