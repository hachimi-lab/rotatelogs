package rotatelogs

import "time"

type (
	RotateType string
)

const (
	EveryMinute RotateType = "minute"
	EveryHour   RotateType = "hour"
	EveryDay    RotateType = "day"
)

func (slf RotateType) isValid() bool {
	return slf == EveryMinute || slf == EveryHour || slf == EveryDay
}

func (slf RotateType) TimeFormat() string {
	switch slf {
	case EveryMinute:
		return "20060102-1504"
	case EveryHour:
		return "20060102-1500"
	case EveryDay:
		return "20060102"
	default:
		return "20060102"
	}
}

func (slf RotateType) UntilNextTime(now time.Time) time.Duration {
	switch slf {
	case EveryMinute:
		year, month, day := now.Date()
		date := time.Date(year, month, day, now.Hour(), now.Minute(), 0, 0, now.Location())
		next := date.Add(time.Minute)
		return next.Sub(now)
	case EveryHour:
		year, month, day := now.Date()
		date := time.Date(year, month, day, now.Hour(), 0, 0, 0, now.Location())
		next := date.Add(time.Hour)
		return next.Sub(now)
	case EveryDay:
		year, month, day := now.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		next := date.Add(time.Hour * 24)
		return next.Sub(now)
	default:
		year, month, day := now.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		next := date.Add(time.Hour * 24)
		return next.Sub(now)
	}
}
