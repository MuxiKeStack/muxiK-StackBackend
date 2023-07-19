package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/config"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
)

// 爬取成绩，并发
func NewGradeRecords(userId uint32, sid, pwd string) error {
	// 获取现有成绩数
	curRecordNum, err := model.GetRecordsNum(userId)
	if err != nil {
		log.Error("GetRecordsNum function error", err)
		return err
	}

	// 教务处获取成绩
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
	finished := make(chan interface{}, 1)
	gradeChan := make(chan *model.GradeModel, 1)

	// 新增记录，并发
	for _, item := range data {
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
				// log.Info("The record has existed")
				return
			}

			gradeChan <- &model.GradeModel{
				UserId:       userId,
				CourseHashId: hash,
				CourseName:   item.CourseName,
				TotalScore:   item.TotalScore,
				UsualScore:   item.UsualScore,
				FinalScore:   item.FinalScore,
				HasAdded:     false,
			}
		}(*item)
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
	finished := make(chan interface{}, 1)

	// 逐个导入，并发
	for _, record := range records {
		wg.Add(1)

		go func(record model.GradeModel) {
			defer wg.Done()

			if err := NewGradeDataAdditionForOneCourse(userId, &record); err != nil {
				log.Error(fmt.Sprintf("NewGradeDataAdditionForOneCourse function error for %t", record), err)
				errChan <- err
			}
		}(*record)
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
// 更新 historyCourse 的成绩字段和 grade 的 hasAdded 字段
func NewGradeDataAdditionForOneCourse(userId uint32, grade *model.GradeModel) error {
	// 获取课程
	course, err := model.GetHistoryCourseByHashId(grade.CourseHashId)
	if err != nil {
		log.Error(fmt.Sprintf("GetHistoryCourseByHashId function error; [hash: %s]", grade.CourseHashId), err)
		return err
	}

	// 数据库数据覆盖问题
	// TO DO: 事务
	curSampleSize := course.GradeSampleSize
	course.TotalGrade = (course.TotalGrade*float32(curSampleSize) + grade.TotalScore) / float32(curSampleSize+1)
	course.UsualGrade = (course.UsualGrade*float32(curSampleSize) + grade.UsualScore) / float32(curSampleSize+1)
	course.GradeSampleSize++

	if grade.TotalScore > 85 {
		course.GradeSection1++
	} else if grade.TotalScore >= 70 {
		course.GradeSection2++
	} else {
		course.GradeSection3++
	}

	// 更新课程数据
	if err := course.Update(); err != nil {
		log.Error("Update history error", err)
		return err
	}

	grade.HasAdded = true
	if err := grade.Update(); err != nil {
		log.Error("Update grade error", err)
		return err
	}

	return nil
}

// 成绩服务，包括成绩爬取和导入统计样本
func GradeImportService(userId uint32, sid, pwd string) {
	log.Info("Crawling grades begins")

	// 获取成绩
	if err := NewGradeRecords(userId, sid, pwd); err != nil {
		log.Error(fmt.Sprintf("Grade import failed for (userId=%d, sid=%s, psw=%s)", userId, sid, pwd), err)
		return
	}
	// 导入成绩样本数据
	if err := NewGradeSampleFoCourses(userId); err != nil {
		log.Error("NewGradeSampleFoCourses function error", err)
		return
	}
	log.Info("Grade sample imported successfully")
}

// 成绩爬取，需要检查是否开启异步爬取
func GradeCrawlHandler(userId uint32, sid, pwd string) {
	// 环境变量设置，是否爬取成绩
	if config.GradeSwitch != "on" {
		return
	}

	// 检查是否加入成绩共享计划
	if ok, err := model.UserHasLicence(userId); err != nil {
		log.Error("UserHasLicence function error", err)
		return
	} else if !ok {
		log.Info(fmt.Sprintf("user(%d) has no licence", userId))
		return
	}

	GradeImportService(userId, sid, pwd)
}
