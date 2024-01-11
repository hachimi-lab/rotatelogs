package rotatelogs

import "time"

type (
	RotateTime int
)

const (
	EveryMinute RotateTime = 60
	EveryHour              = EveryMinute * 60
	EveryDay               = EveryHour * 24
)

func (slf RotateTime) TimeFormat() string {
	switch slf {
	case EveryMinute:
		return "20060102-1504"
	case EveryHour:
		return "20060102-1500"
	case EveryDay:
		return "20060102"
	default:
		return ""
	}
}

func (slf RotateTime) UntilNextTime(now time.Time) time.Duration {
	switch slf {
	case EveryMinute:
		year, month, day := now.Date()
		date := time.Date(year, month, day, now.Hour(), now.Minute(), 0, 0, now.Location())
		next := date.Add(time.Duration(slf) * time.Second)
		return next.Sub(now)
	case EveryHour:
		year, month, day := now.Date()
		date := time.Date(year, month, day, now.Hour(), 0, 0, 0, now.Location())
		next := date.Add(time.Duration(slf) * time.Second)
		return next.Sub(now)
	case EveryDay:
		year, month, day := now.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		next := date.Add(time.Duration(slf) * time.Second)
		return next.Sub(now)
	default:
		return 0
	}
}
