package wolfram

import (
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	ans, err := QueryWolfram("8:17pm EDT on July 20th, 1969")
	//ans, err := QueryWolfram("3/14/19")
	if err != nil {
		t.Error("err with query", err)
	}
	tz, _ := time.LoadLocation("America/New_York")
	tstt := time.Date(1969, time.July, 20, 20, 17, 0, 0, tz)
	if ans.Sub(tstt).Seconds() != 0 {
		t.Error("err", ans, tstt)
	}
	//fmt.Println(ans, err)
}
