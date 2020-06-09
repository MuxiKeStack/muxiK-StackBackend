package service

import (
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/lexkong/log"
)

// 爬取成绩，并发
func NewGradeRecord(userId uint32, sid, pwd string) error {
	// 获取现有成绩数
	curRecordNum, err := model.GetRecordsNum(userId)
	if err != nil {
		log.Error("GetRecordsNum function error", err)
		return err
	}

	//教务处获取成绩
	data, ok, err := util.GetGradeFromXK(sid, pwd, curRecordNum)
	if err != nil {
		log.Error("util.GetGradeFromXK function error", err)
		return err
	} else if !ok {
		// 若成绩未更新则不处理
		log.Info("Grades have not updated")
		return nil
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)
	gradeChan := make(chan *model.GradeModel, 1)

	// 新增记录，并发
	for _, item := range *data {
		wg.Add(1)

		go func(item util.ResultGradeItem) {
			defer wg.Done()

			teacher := strings.ReplaceAll(item.Teacher, ";", ",")
			hash := util.HashCourseId(item.CourseId, teacher)

			// 验证该记录是否已存在
			if ok, err := model.GradeRecordExisting(userId, hash); err != nil {
				log.Error("GradeRecordExisting function error", err)
				errChan <- err
				return
			} else if ok {
				//log.Info("The record has existed")
				return
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
			gradeChan <- g
		}(item)
	}

	go func() {
		wg.Wait()
		close(gradeChan)
	}()

	go func() {
		for g := range gradeChan {
			if err := g.New(); err != nil {
				log.Error("Add new grade record error", err)
				errChan <- err
				return
			}
		}
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return err
	}

	return nil
}

// 将成绩导入到各课程的成绩统计中
func NewGradeSampleFoCourses(userId uint32) error {
	// 获取未添加导入的成绩数据
	records, err := model.GetGradeRecordsByUserId(userId)
	if err != nil {
		log.Error("GetGradeRecordsByUserId function error", err)
		return err
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 逐个导入，并发
	for _, record := range *records {
		wg.Add(1)

		go func(record model.GradeModel) {
			defer wg.Done()

			if err := NewGradeDataAdditionForOneCourse(userId, &record); err != nil {
				log.Error("NewGradeDataAdditionForOneCourse function error", err)
				errChan <- err
			}
		}(record)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return err
	}
	return nil
}

// 一门课程的成绩样本数据添加
func NewGradeDataAdditionForOneCourse(userId uint32, data *model.GradeModel) error {
	// 获取课程
	course, _, err := model.GetHistoryCourseByHashId(data.CourseHashId)
	if err != nil {
		log.Error("GetHistoryCourseByHashId function error", err)
		return err
	}

	// 有数据库写入覆盖的隐患
	curSampleSize := course.GradeSampleSize
	course.TotalGrade = (course.TotalGrade*float32(curSampleSize) + data.TotalScore) / float32(curSampleSize+1)
	course.UsualGrade = (course.UsualGrade*float32(curSampleSize) + data.UsualScore) / float32(curSampleSize+1)
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

// 成绩服务，包括成绩爬取和导入统计样本
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
