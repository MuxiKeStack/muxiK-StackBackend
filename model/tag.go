package model

import (
	"strconv"
	"strings"
)

// Get tag name
func GetTagName(tagId uint64) string {

	return ""
}

// Convert tags from string to array
func TagStrToArray(s string) []uint8 {
	var tagIds []uint8
	tagsStr := strings.Split(s, ",")
	for _, s := range tagsStr {
		id, _ := strconv.Atoi(s)
		tagIds = append(tagIds, uint8(id))
	}
	return tagIds
}

// Convert tags from array to string
func TagArrayToStr(tagIds []uint8) string {
	var s string
	for _, id := range tagIds {
		s = strconv.FormatUint(uint64(id), 10) + ","
	}
	return s
}
