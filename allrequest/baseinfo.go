package allrequest

import (
	"bawangcan/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func getbaseheaders() map[string]string {
	headers := map[string]string{
		"Origin":           "http://s.dianping.com",
		"Accept-Encoding":  "gzip, deflate",
		"X-Request":        "JSON",
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
		"Content-Type":     "application/json;charset=UTF-8",
		"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept":           "application/json, text/javascript",
		"Referer":          "http://s.dianping.com/event/1063463422",
		"X-Requested-With": "XMLHttpRequest",
		"Connection":       "keep-alive",
	}
	return headers
}

func Getcityid(con config.Config) int {
	headers := getbaseheaders()
	//ur , _ := url2.Parse("http://127.0.0.1:10809")
	city := con.Cityname
	client := &http.Client{
		Transport: &http.Transport{},
	}

	re, err := http.NewRequest("GET", "http://s.dianping.com/event/"+city, nil)
	if err != nil {
		println(err.Error())
	}
	addhead(re, headers)
	resp, err := client.Do(re)
	if err != nil {
		println(err.Error())
	}
	resp = checkgzip(resp)

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()
	println(resp.Status)
	//println(string(content))
	id := strings.Split(strings.Split(string(content), "\"cityId\":\"")[1], "\"")[0]
	cityid, err := strconv.Atoi(id)
	println(cityid)
	return cityid
}
