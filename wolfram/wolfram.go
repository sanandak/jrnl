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
	"time"
)

var (
	// possible return forms of time from wolfram
	wolframTimeLayout = [...]string{
		"3:04 pm MST | Monday, January 2, 2006",
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
	Primary bool `json:"primary"`
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
	// look for the "Primary" pod
	for _, p := range data.Res.Pods {
		//fmt.Println(i, p)
		if p.Primary {
			wTimeStr := p.SubPods[0].Plaintext
			fmt.Println("got it", p, wTimeStr)
			for _, layout := range wolframTimeLayout {
				wtime, err = time.Parse(layout, wTimeStr)
				if err == nil {
					break
				}
			}
			// TODO - parse answers of form x months 7 days ago
			if err != nil {
				log.Println("err parsing returned time", wtime, err)
				return time.Now(), err
			}
			//fmt.Println(wtime, err)
			return wtime, nil
		}
	}
	//fmt.Printf("%+v\n", data.Res)
	return time.Now(), errors.New("no primary answer")
}
