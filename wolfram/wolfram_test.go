package wolfram

import (
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	// FIXME - use local time zone
	tz, _ := time.LoadLocation("America/New_York")
	//tstt := time.Date(1969, time.July, 20, 20, 17, 0, 0, tz)
	var tests = map[string]time.Time{
		"8:17pm EDT on July 20th, 1969": time.Date(1969, time.July, 20, 20, 17, 0, 0, tz),
	}
	for tst, tstT := range tests {
		ans, err := QueryWolfram(tst)

		if err != nil {
			t.Errorf("err with query %s: %s", tst, err)
		}
		if ans.Sub(tstT).Seconds() != 0 {
			t.Errorf("error parasing %s. Expected %v got %v", tst, tstT, ans)
		}
	}

	//fmt.Println(ans, err)
}
