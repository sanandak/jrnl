package entry

import (
	"math"
	"testing"
	"time"
)

var (
	titleStr   = "This is the title"
	textStr    = "First sentence.      Second one!"
	textStrStd = "First sentence. Second one!"
)

func TestEmptyEntry(t *testing.T) {
	_, err := NewEntry("")
	if err == nil {
		t.Error("expected \"empty entry\" return got nil")
	}
}
func TestTime(t *testing.T) {
	now := time.Now()
	entry, err := NewEntry("test")
	if err != nil {
		t.Error("expected nil err, got ", err)
	}
	if entry.rawStr != "test" {
		t.Errorf("expected \"%+v\" got %+v", "test", entry.rawStr)
	}
	if math.Abs(entry.entryTime.Sub(now).Seconds()) > 1 {
		t.Errorf("expected entrytime of %v got %v", now, entry.entryTime)
	}
}

func TestParseTitle(t *testing.T) {
	entry, _ := NewEntry(titleStr + ".")
	if entry.title != titleStr {
		t.Errorf("expected \"%v\" got \"%v\"", titleStr, entry.title)
	}
	//fmt.Printf("test parse title: %+v got %+v\n", titleStr, entry.title)
}

func TestParseText(t *testing.T) {
	entry, _ := NewEntry(titleStr + "." + textStr)
	if entry.text != textStrStd {
		t.Errorf("expected \"%v\" got \"%v\"", textStrStd, entry.text)
	}
	//fmt.Printf("test parse text: %+v %+v\n", textStr, entry.text)
}

func TestParseWhen(t *testing.T) {
	entry, _ := NewEntry("friday: title. text")
	// but is it the right friday???
	if entry.when.Weekday() != time.Friday {
		t.Errorf("err parsing %s. Got %+v expected %+v", entry.whenStr, entry.when.Weekday(), time.Friday)
	}
}

func TestParseTime(t *testing.T) {
	entry, _ := NewEntry("7pm today: title. text")
	if entry.when.Hour() != 19 {
		t.Errorf("err parsing %s. Got %+v expected %+v", entry.whenStr, entry.when.Hour(), 19)
	}
}
func TestParseDate(t *testing.T) {
	entry, _ := NewEntry("3/14/2019: title. text")
	if entry.when.Month() != time.March && entry.when.Day() != 14 {
		t.Errorf("err parsing %s. Got %+v expected 2019/3/14", entry.whenStr, entry.when)
	}
}
func TestTags(t *testing.T) {
	entry, _ := NewEntry("today: title. has @tag.")
	if entry.tags[0] != "@tag" {
		t.Errorf("err parsing tags in %s. got %s expected @tag", entry.text, entry.tags)
	}
}
