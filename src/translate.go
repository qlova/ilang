package main

import "net/http"
import "net/url"
import "bytes"
import "strings"
import "fmt"
import "errors"
import "encoding/json"
import "io/ioutil"

func getTranslation(source, target, text string) (res string, err error) {
	uri := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&dt=bd&q=%s"
	uri = fmt.Sprintf(uri, source, target, url.QueryEscape(text))

	var req *http.Request
	if req, err = http.NewRequest("GET", uri, nil); err != nil {
		return
	}
	req.Header.Add("User-Agent", "")

	hc := new(http.Client)
	var resp *http.Response
	if resp, err = hc.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		err = errors.New("Google Translate: " + resp.Status)
		return
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	// Fixes bad JSON
	var prev rune
	body = bytes.Map(func(r rune) rune {
		if r == 44 && prev == 44 {
			return 32
		}
		prev = r
		return r
	}, body)

	var data []interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	// Concatenates output
	// Hold on tight, here we go. Aaah!!!
	for _, v := range data[:1] {
		for _, v := range v.([]interface{}) {
			res += v.([]interface{})[0].(string)
		}
	}
	if len(data) > 2 {
		for _, v := range data[1:2][0].([]interface{}) {
			res += "\n" + v.([]interface{})[0].(string) + ": "
			for _, v := range v.([]interface{})[1].([]interface{}) {
				res += v.(string) + ", "
			}
			res = strings.TrimRight(res, ", ")
		}
		res = strings.ToLower(res)
	}

	return
}
