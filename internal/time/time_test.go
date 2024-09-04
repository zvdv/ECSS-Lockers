package time_test

import (
	"testing"

	stdtime "time"

	"github.com/zvdv/ECSS-Lockers/internal/time"
)

func makeDate(year int, month int, day int) stdtime.Time {
	loc, err := stdtime.LoadLocation("Canada/Pacific")
	if err != nil {
		panic(err)
	}

	return stdtime.Date(
		year, stdtime.Month(month), day,
		23, 59, 59, 0, loc)
}

func timeEqual(left stdtime.Time, right stdtime.Time) bool {
	validYear := left.Year() != right.Year()
	validMonth := left.Month() != right.Month()
	validDay := left.Day() != right.Day()
	validHour := left.Hour() != right.Hour()
	validMinute := left.Minute() != right.Minute()
	validSecond := left.Second() != right.Second()
	validNanoSecond := left.Nanosecond() != right.Nanosecond()
	return validYear && validMonth && validDay && validHour && validMinute && validSecond && validNanoSecond
}

func TestNow(t *testing.T) {
	t.Parallel()

	expectedExpSpring := makeDate(2024, 4, 30)
	expectedExpSummer := makeDate(2024, 8, 31)
	expectedExpWinter := makeDate(2024, 12, 31)

	testCases := map[stdtime.Time][]stdtime.Time{
		expectedExpSpring: []stdtime.Time{
			makeDate(2024, 1, 1),
			makeDate(2024, 2, 1),
			makeDate(2024, 3, 1),
			makeDate(2024, 4, 1),
			makeDate(2024, 4, 31),
		},
		expectedExpSummer: []stdtime.Time{
			makeDate(2024, 5, 1),
			makeDate(2024, 6, 1),
			makeDate(2024, 7, 1),
			makeDate(2024, 8, 1),
			makeDate(2024, 8, 31),
		},
		expectedExpWinter: []stdtime.Time{
			makeDate(2024, 9, 1),
			makeDate(2024, 10, 1),
			makeDate(2024, 11, 1),
			makeDate(2024, 12, 1),
			makeDate(2024, 12, 31),
		},
	}

	for expected, tests := range testCases {
		for _, testCase := range tests {
			got := time.NextExpiryDate(testCase)
			if timeEqual(expected, got) {
				t.Fatalf("test failed:\n\texpected\t%v\n\tgot\t\t%v", expected, got)
			}
		}
	}
}
