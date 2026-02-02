package convert

import "strconv"

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func Int32ToStr(num int32) string {
	return strconv.Itoa(int(num))
}

func Int16ToStr(num int16) string {
	return strconv.Itoa(int(num))
}
