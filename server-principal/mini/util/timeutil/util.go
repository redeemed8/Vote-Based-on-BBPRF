package timeutil

import (
	"strconv"
	"time"
)

// GetYMDStr 	获取 年月日
func GetYMDStr() (year, month, day string) {
	currentTime := time.Now()
	yearNumber := currentTime.Year()
	monthNumber := int(currentTime.Month())
	dayNumber := currentTime.Day()
	year = strconv.Itoa(yearNumber)
	month = strconv.Itoa(monthNumber)
	day = strconv.Itoa(dayNumber)
	if len(month) == 1 {
		month = "0" + month
	}
	if len(day) == 1 {
		day = "0" + day
	}
	return
}
