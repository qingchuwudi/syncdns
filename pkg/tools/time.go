package tools

import (
	"math/rand/v2"
	"time"
)

// NextDay 下一天00:00
func NextDay() time.Duration {
	currentTime := time.Now()
	next := currentTime.AddDate(0, 0, 1)
	return next.Sub(currentTime)
}

// NextHour 计算并返回下一个整点小时距离当前时刻的时间
func NextHour() time.Duration {
	return NextSegment(time.Hour)
}

func NextHourWithRand() time.Duration {
	return NextSegmentWithRand(time.Hour, int(time.Minute))
}

// NextMinute 将24小时从00:00开始按照每1分钟做一次分割，计算并输出下一个分割点距离当前的时间
func NextMinute() time.Duration {
	return NextSegment(time.Minute)
}

// NextTenMinute 将24小时从00:00开始按照每10分钟做一次分割，计算并输出下一个分割点距离当前的时间
func NextTenMinute() time.Duration {
	return NextSegment(10 * time.Minute)
}

// NextSegment 根据时间间隔计算并返回下一个时间点
//
// segment : 时间间隔
func NextSegment(segment time.Duration) time.Duration {
	currentTime := time.Now()
	nextSegment := GetNextSegment(currentTime, segment)
	return nextSegment.Sub(currentTime)
}

// GetNextSegment 支持按照指定时间间隔进行分割
func GetNextSegment(currentTime time.Time, segmentDuration time.Duration) time.Time {
	// 计算当前时间距离当天的00:00的时间间隔
	timeSinceMidnight := currentTime.Sub(currentTime.Truncate(24 * time.Hour))

	// 计算下一个分割点的时间，segmentTrun 是为了对浮点型数据取整
	// nolint
	segmentTrun := timeSinceMidnight / segmentDuration * segmentDuration
	nextSegment := currentTime.Truncate(24 * time.Hour).Add(segmentTrun).Add(segmentDuration)

	return nextSegment
}

// NextSegmentWithRand 为时间添加一个随机的延迟，n是延迟的范围上限，时间单位是秒数
func NextSegmentWithRand(segment time.Duration, n int) time.Duration {
	second := rand.N(n)
	t := time.Duration(second) * time.Second
	segment = NextSegment(segment)
	return segment + t
}
