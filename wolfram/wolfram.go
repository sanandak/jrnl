// Package wolfram will query WOlframAlpha for the time from the `when` string
// log in to developer.wolfram.com and create a new app
// store app id in env variable WOLFRAMAPPID
package wolfram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	// possible return forms of time from wolfram
	wolframTimeLayout = [...]string{
		"Monday, January 2, 2006 at 3:04 pm MST",
		"3:04 pm MST | Monday, January 2, 2006",
		"3:04:05 pm MST | Monday, January 2, 2006",
		"Monday, January 2, 2006",
	}
	// need to parse "x months y days ago/from now" and "x days ago/from now"
)

// borrowed from github.com/Krognol/go-wolfram

// QueryResp is the top level response from WA
type QueryResp struct {
	Res QueryRes `json:"queryresult"`
}

// QueryRes is the actual query result - consists of pods and other things
type QueryRes struct {
	Pods    []Pod `json:"pods"`
	NumPods int   `json:"numpods"`
}

// Pod is an object (up to NumPod Pods) with the answers
// the key Primary tells us which pod is, well, the primary one
type Pod struct {
	//The subpod elements of the pod
	SubPods []SubPod `json:"subpods"`
	//Marks the pod that displays the closest thing to a simple "answer" that Wolfram|Alpha can provide
	Primary bool   `json:"primary"`
	Title   string `json:"title"`
}

// SubPod - there is one subpod per pod, and the plaintext has the answer in...plaintext
type SubPod struct {
	//Textual representation of the subpod
	Plaintext string `json:"plaintext"`
}

var (
	baseURL = "https://api.wolframalpha.com/v2/query"
	appid   string
)

// QueryWolfram will send a natural language time query to wolfram alpha and
// parse the returned JSON string for the `primary` pod
// return time (or time.Now()) for error
func QueryWolfram(query string) (wtime time.Time, err error) {
	appid = os.Getenv("WOLFRAMAPPID")
	if len(appid) == 0 {
		log.Printf("No Wolfram App ID")
		return time.Now(), errors.New("no wolfram app id")
	}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Println("new request ", err)
		return time.Now(), err
	}
	// new query
	q := req.URL.Query()
	q.Add("appid", os.Getenv("WOLFRAMAPPID"))
	q.Add("output", "JSON")
	q.Add("input", query)

	req.URL.RawQuery = q.Encode()
	//fmt.Println(req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client req", err)
		return time.Now(), err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(resp.Status)
	//fmt.Println(string(respBody))

	data := &QueryResp{}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		log.Println("err unmarshalling", respBody)
		return time.Now(), err
	}
	// look for primary pod
	for _, p := range data.Res.Pods {
		//fmt.Println(i, p)
		if p.Primary {
			wTimeStr := p.SubPods[0].Plaintext
			//fmt.Println("got it", wTimeStr)
			for _, layout := range wolframTimeLayout {
				wtime, err = time.Parse(layout, wTimeStr)
				if err == nil {
					//fmt.Println("got it!!", wtime, err)
					return wtime, nil
				}
			}
		}
	}

	// no luck in the primary, look in the "Input interpretation" pod
	for _, p := range data.Res.Pods {
		//fmt.Println(i, p)
		if p.Title == "Input interpretation" {
			wTimeStr := p.SubPods[0].Plaintext
			//fmt.Println("ii got it", wTimeStr)
			for _, layout := range wolframTimeLayout {
				wtime, err = time.Parse(layout, wTimeStr)
				if err == nil {
					//fmt.Println("ii got it!!", wtime, err)
					return wtime, nil
				}
			}
			if err != nil {
				log.Println("err parsing returned time", p.SubPods, err)
				return time.Now(), errors.New("err parsing when")
			}
		}
	}

	//fmt.Printf("%+v\n", data.Res)
	return time.Now(), errors.New("no primary answer")
}

func parseAgo(agoStr string) (wtime time.Time, err error) {
	nmon := matchPart(agoStr, "months")
	nday := matchPart(agoStr, "days")
	nweek := matchPart(agoStr, "weeks")
	fmt.Println(agoStr, nmon, nday, nweek)

	return time.Now(), nil
}

func matchPart(str string, val string) int {
	matcher := fmt.Sprintf("(\\d+)\\s+(?:\\b%s?\\b)", val)
	fmt.Println(matcher)
	theRE := regexp.MustCompile(matcher)

	var theVal int
	var err error
	if v := theRE.FindAllStringSubmatch(str, 1); len(v) > 0 {
		//fmt.Printf("ago: %v mon: %+v\n", agoStr, mon[0][1])
		if len(v[0]) > 0 {
			theVal, err = strconv.Atoi(v[0][1])
			if err != nil {
				theVal = 0
			}
		}
	}
	fmt.Println(theVal)
	return theVal
}
