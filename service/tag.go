package service

import (
	"fmt"
	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"strconv"
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

// Get tag real names by id string.
func GetTagNamesByIdStr(s string) ([]string, error) {
	var tagNames []string
	if s == "" {
		return tagNames, nil
	}

	tagsStr := strings.Split(s, ",")

	for _, idStr := range tagsStr {
		// id should not be zero, skip it
		if idStr == "0" {
			continue
		}

		id, _ := strconv.Atoi(idStr)
		name, err := model.GetTagNameById(id)
		if err != nil {
			return nil, err
		}

		tagNames = append(tagNames, name)
	}

	return tagNames, nil
}

// Convert tags from array to string
func TagArrayToStr(tagIds []uint8) string {
	var s string
	for i, id := range tagIds {
		// tagId should not be zero, skip it
		if id == 0 {
			continue
		}

		if i > 0 {
			s = fmt.Sprintf("%s,%d", s, id)
			continue
		}
		s = fmt.Sprintf("%d", id)
	}
	return s
}

// Update course's tag amount after publishing a new evaluation.
func NewTagsAfterPublishing(tags *[]uint8, courseId string) error {
	for _, tag := range *tags {
		err := model.NewTagsForCourse(uint32(tag), courseId)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get two most tags of a course by its hash id.
func GetTwoMostTagsOfOneCourse(courseId string) ([]string, error) {
	tagIds, err := model.GetTwoMostTagIdsOfCourseByHashId(courseId)
	if err != nil {
		log.Error("GetTwoMostTagsOfCourseByHashId function error", err)
		return nil, err
	}

	var result []string
	for _, id := range tagIds {
		tag, err := model.GetTagNameById(id)
		if err != nil {
			log.Error("GetTagNameById function error", err)
			return nil, err
		}
		result = append(result, tag)
	}

	return result, nil
}
