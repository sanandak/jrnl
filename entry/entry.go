// Package entry manages diary entries by parsing strings and populating
// the Entry type.
//
// Strings are of the form [when:] [title.] entry
//
//   `[when.]` is an optional entry in natural language (today, next wednesday, etc.)
//      followed by a period (`.`)
//   `[title:]` is an optional title string ending with a colon (`:`)
//   `entry` is the text of the diary entry
//
// This string is parsed and formatted as an org entry.
// For example:
// `jrnl today. title is here: entry text with @tag1 and @tag2`
/*
   ** Title Is Here                                    :tag1:tag2:
   <2019/06/11 Tue 15:32>
   entry text with @tag1 and @tag2
*/
// I use the WolframAlpha to parse the time.
package entry

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/olebedev/when"
	"github.com/sanandak/jrnl/wolfram"

	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

// Entry defines a journal entry
type Entry struct {
	entryTime time.Time
	rawStr    string
	title     string
	text      string
	tags      []string
	whenStr   string
	when      time.Time
	whenErr   bool
}

const (
	// TimeStampLayout is the format for the org-mode timestamp
	TimeStampLayout = "<2006/01/02 Mon 15:04>"
)

// NewEntry creates a new entry from `rawStr`
// returns a pointer *Entry and error
func NewEntry(rawStr string) (entry *Entry, err error) {
	err = nil
	if len(strings.TrimSpace(rawStr)) == 0 {
		return nil, errors.New("empty entry")
	}
	entry = &Entry{}
	entry.rawStr = rawStr
	rawStr = strings.TrimSpace(rawStr)

	// find tags of form @xxx
	tagre := regexp.MustCompile(`@([A-Za-z0-9_]+)`)
	tags := tagre.FindAllString(rawStr, -1)
	entry.tags = tags

	var whenIdx, titleIdx int
	// search for `when string.`
	whenIdx = strings.IndexByte(rawStr, '.')
	if whenIdx > 0 {
		entry.whenStr = rawStr[:whenIdx]
		rawStr = rawStr[whenIdx+1:]
	}
	if whenIdx == 0 { // bare `.`
		entry.whenStr = "today"
		rawStr = rawStr[1:]
	}
	if whenIdx < 0 { // no when string
		entry.whenStr = "today"
	}

	// search for `title string:`
	titleIdx = strings.IndexByte(rawStr, ':')
	if titleIdx > 0 {
		entry.title = standardizeSpaces(rawStr[:titleIdx])
		entry.text = standardizeSpaces(rawStr[titleIdx+1:])
	}
	if titleIdx == 0 { // bare : at start; no title
		entry.text = standardizeSpaces(rawStr[1:])
	}
	entry.entryTime = time.Now()
	entry.parseWhenWolfram()
	return entry, err

	// tags @([A-Za-z0-9_]+)
}

// Print makes an org-mode string from entry
func (entry *Entry) Print() []byte {
	var buf []byte
	var tagStr = ""
	var out = bytes.NewBuffer(buf)

	// put the tags at the end of the headline
	// org-mode tags are delimited by : `:tag:tag2:`
	titleLen := len(entry.title) + 2
	if len(entry.tags) > 0 {
		var noAt []string
		atRE := regexp.MustCompile(`^@`)
		for _, t := range entry.tags {
			noAt = append(noAt, atRE.ReplaceAllString(t, ""))
		}
		tagStr = strings.Join(noAt, ":")
	}
	//fmt.Println(tagStr, entry.tags)
	// headline should have tags right justified
	// ** Title            :tag:<--column 80
	if tagLen := len(tagStr); tagLen > 0 {
		tagStr = ":" + tagStr + ":"
		tagFmtLen := 80 - titleLen
		if tagFmtLen < tagLen {
			tagFmtLen = tagLen
		}
		//fmt.Println(tagFmtLen)
		fmt.Fprintf(out, "** %s %[2]*s\n", strings.Title(entry.title), tagFmtLen, tagStr)
	} else {
		fmt.Fprintf(out, "** %s\n", strings.Title(entry.title))
	}

	out.Write([]byte(entry.when.Format(TimeStampLayout)))
	if entry.whenErr {
		fmt.Fprintf(out, " [* %s]\n", entry.whenStr)
	} else {
		fmt.Fprint(out, "\n")
	}
	fmt.Fprintf(out, "%s\n\n", entry.text)
	//fmt.Println(out)
	return out.Bytes()
}

// parseWhen uses olebedev/when to parse the `when` string
func (entry *Entry) parseWhen() (err error) {
	w := when.New(nil)

	w.Add(en.All...)
	w.Add(common.All...)

	r, err := w.Parse(entry.whenStr, time.Now())
	if err == nil && r != nil {
		entry.when = r.Time
	} else {
		entry.when = time.Now()
		entry.whenErr = true
	}
	//fmt.Printf("r: %+v %+v\n", r, err)

	return nil
}

// parseWhenWolfram uses sanandak/wolfram and wolframAlpha to parse `when`
func (entry *Entry) parseWhenWolfram() (err error) {
	wtime, err := wolfram.QueryWolfram(entry.whenStr)
	if err == nil {
		entry.when = wtime
	} else {
		entry.when = time.Now()
		entry.whenErr = true
	}
	//fmt.Printf("r: %+v %+v\n", wtime, err)

	return nil
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
