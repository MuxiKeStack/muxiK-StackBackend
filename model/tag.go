package model

import "strconv"

// Get tag name
func GetTagName(tagId uint64) string {

	return ""
}

// Convert tags from string to array
func TagStrToArray(s string) []uint8 {
	var data []uint8

	return data
}

// Convert tags from array to string
func TagArrayToStr(t []uint8) string {
	var s string
	for _, i := range t {
		s = strconv.FormatUint(uint64(i), 10) + ","
	}
	return s
}
