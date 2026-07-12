package router

import "time"

// jwtExpire 把小时数转为 time.Duration。
func jwtExpire(hours int) time.Duration {
	if hours <= 0 {
		hours = 168 // 默认 7 天
	}
	return time.Duration(hours) * time.Hour
}
