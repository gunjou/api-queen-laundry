package utils

import "time"

func GetNowWITA() time.Time {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return time.Now().In(loc)
}