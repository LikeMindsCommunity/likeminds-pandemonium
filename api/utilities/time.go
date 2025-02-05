package utilities

func GetMillisecondsInMinutes(minutes int) int64 {
	return int64(minutes) * 60 * 1000
}
