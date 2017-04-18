package ilang

import "net/http"
import "net/url"
//import "bytes"
//import "strings"
import "fmt"
import "errors"
//import "encoding/json"
import "io/ioutil"

func getTranslation(source, target, text string) (res string, err error) {
	//fmt.Println("translating ", text)
	uri := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&dt=bd&q=%s"
	uri = fmt.Sprintf(uri, source, target, url.QueryEscape(text))
	
	//println(uri)

	var req *http.Request
	if req, err = http.NewRequest("GET", uri, nil); err != nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")

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
	
	//println(string(body))
	
	var word string
	var appending bool
	for _, char := range body {
		if char == '"' {
			appending = !appending
			if !appending {
				break
			}
		} else {
		if appending {
			word += string(char)
		}
		}
	}
	//println(word)
	res = word
	return
	/*
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
	for _, v := range data{
		for _, v := range v.([]interface{}) {
			res += v.([]interface{})[0].(string)
			println(res)
		}
	}
	if len(data) > 2 {
		var over, ok = data[1:2][0].([]interface{})
		if !ok {
			println("error translating!")
			return
		}
		for _, v := range over {
			res += "\n" + v.([]interface{})[0].(string) + ": "
			for _, v := range v.([]interface{})[1].([]interface{}) {
				res += v.(string) + ", "
			}
			res = strings.TrimRight(res, ", ")
		}
		res = strings.ToLower(res)
	}

	return*/
}
