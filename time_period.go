package rotatelogs

import "time"

type (
	TimePeriod string
)

const (
	Minutely TimePeriod = "minutely"
	Hourly   TimePeriod = "hourly"
	Daily    TimePeriod = "daily"
	Monthly  TimePeriod = "monthly"
)

func (slf TimePeriod) isValid() bool {
	return slf == Minutely || slf == Hourly || slf == Daily || slf == Monthly
}

func (slf TimePeriod) TimeFormat() string {
	switch slf {
	case Minutely:
		return "20060102-1504"
	case Hourly:
		return "20060102-1500"
	case Daily:
		return "20060102"
	case Monthly:
		return "20060101"
	default:
		return "20060102"
	}
}

func (slf TimePeriod) UntilNextTime(now time.Time) time.Duration {
	switch slf {
	case Minutely:
		year, month, day := now.Date()
		date := time.Date(year, month, day, now.Hour(), now.Minute(), 0, 0, now.Location())
		next := date.Add(time.Minute)
		return next.Sub(now)
	case Hourly:
		year, month, day := now.Date()
		date := time.Date(year, month, day, now.Hour(), 0, 0, 0, now.Location())
		next := date.Add(time.Hour)
		return next.Sub(now)
	case Daily:
		year, month, day := now.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		next := date.AddDate(0, 0, 1)
		return next.Sub(now)
	case Monthly:
		year, month, _ := now.Date()
		date := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		next := date.AddDate(0, 1, 0)
		return next.Sub(now)
	default:
		year, month, day := now.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		next := date.AddDate(0, 0, 1)
		return next.Sub(now)
	}
}
