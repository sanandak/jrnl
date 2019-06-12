package entry

import (
	"fmt"
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
		t.Error("expected nil return got ", err)
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
	fmt.Printf("test parse entry: %+v when: %+v\n", entry, entry.when.String())
}
