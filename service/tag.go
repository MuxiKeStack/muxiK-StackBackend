package service

import (
	"strconv"
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

// Get tag real names by id string.
func GetTagNamesByIdStr(s string) []string {
	var tagNames []string
	tagsStr := strings.Split(s, ",")

	for _, idStr := range tagsStr {
		id, _ := strconv.Atoi(idStr)
		name := model.GetTagNameById(id)

		tagNames = append(tagNames, name)
	}

	return tagNames
}

// Convert tags from array to string
func TagArrayToStr(tagIds []uint8) string {
	var s string
	for _, id := range tagIds {
		s = strconv.FormatUint(uint64(id), 10) + ","
	}
	return s
}
