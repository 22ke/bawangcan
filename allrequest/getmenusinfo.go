package allrequest

import (
	"bawangcan/config"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Getmenuinfos(config config.Config) (map[int]string, []int) {
	uri, _ := url.Parse("http://127.0.0.1:8080")
	act := make(map[int]string)
	var actid []int
	cityId := Getcityid(config)
	println(cityId)
	types_food := config.Menu
	for _, i := range types_food {
		for page := 1; page <= 15; page++ {
			data := map[string]interface{}{
				"cityId": cityId,
				"type":   i,
				"mode":   "",
				"page":   page,
			}
			//println(data["type"])
			bytesData, _ := json.Marshal(data)
			reader := bytes.NewReader(bytesData)
			var client *http.Client
			if config.Debug == true {
				client = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(uri),
					},
				}
			} else {
				client = &http.Client{}
			}

			r, err := http.NewRequest("POST", "http://m.dianping.com/activity/static/pc/ajaxList", ioutil.NopCloser(reader))
			if err != nil {
				println(err.Error())
			}
			r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
			r.Header.Set("Referer", "http://s.dianping.com/event/shanghai")
			r.Header.Set("Origin", "http://s.dianping.com")
			r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
			r.Header.Set("Accept-Encoding", "gzip, deflate")
			r.Header.Set("Content-Type", "application/json;charset=UTF-8")
			r.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
			rs, err := client.Do(r)
			if err != nil {
				println(err.Error())
			}
			rs = checkgzip(rs)
			l := &list{}
			content, _ := ioutil.ReadAll(rs.Body)
			//string(content)
			json.Unmarshal(content, l)
			num := len(l.Data.Details)
			for i := 0; i < num; i++ {
				actid = append(actid, l.Data.Details[i].OfflineActivityId)
				act[l.Data.Details[i].OfflineActivityId] = l.Data.Details[i].ActivityTitle
			}
			if l.Data.HasNext == false {
				break
			}
		}
	}
	return act, actid
}
