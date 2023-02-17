package helper

import "time"

type PrintType string

const (
	TimeLayout PrintType = "2006-01-02 15:04:05"
	ExcelTime  PrintType = "2006-01-02 15-04-05"
	NumberTime PrintType = "20060102150405"
)

func PrintTime(inType PrintType) string {
	return time.Now().Format(string(inType))
}

//获取传入时间戳月份的最后一天
func GetMonthLastDay(d time.Time) time.Time {
	return GetMonthFirstDay(d).AddDate(0, 1, -1)
}

//获取传入时间戳月份的第一天
func GetMonthFirstDay(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func GetMonday(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return GetZeroTime(t).AddDate(0, 0, offset)
}

func GetSunday(monday time.Time) time.Time {
	sunday := monday.AddDate(0, 0, 6)
	return GetZeroTime(sunday)
}

func StringTimeFormat(timeStr string) time.Time {
	t, e := time.ParseInLocation(string(TimeLayout), timeStr, time.Local)
	if e != nil {
		return time.Time{}
	}
	return t
}
