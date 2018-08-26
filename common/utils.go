package common

import (
	"time"
)

func ConvertTime(timing string) time.Time {
	var m = int(time.Now().Month())
	loc, err := time.LoadLocation("Local")
	CheckError(err)
	parsedTime, err := time.ParseInLocation("15:04:05", timing, loc)
	if err != nil {
		parsedTime = time.Time{}
	}
	if !parsedTime.IsZero() {
		parsedTime = parsedTime.AddDate(time.Now().Year(), m-1, time.Now().Day()-1)
	}
	return parsedTime

}

func CheckError(e error) {
	if e != nil {
		Info.Println(e)
	}
}
