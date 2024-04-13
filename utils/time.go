package utils

import "time"

func CurrentMillisecond() int64 {
	return time.Now().UnixNano() / 1000000
}

func CurrentSecond() int64 {
	return time.Now().Unix()
}
