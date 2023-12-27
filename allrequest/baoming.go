package allrequest

import (
	"bawangcan/config"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type enroolResult struct {
	Code int `json:"code"`
	Data Da  `json:"data"`
}

type Da struct {
	Desc string `json:"desc"`
}

type list struct {
	Code int  `json:"code"`
	Data Data `json:"data"`
}

type Data struct {
	HasNext bool     `json:"hasNext"`
	Details []Detail `json:"detail"`
}

type Detail struct {
	OfflineActivityId int    `json:"offlineActivityId"`
	ActivityTitle     string `json:"activityTitle"`
	DetailUrl         string `json:"detailUrl"`
}

func Baoming(act map[int]string, actid []int, con config.Config, f *os.File) {
	//act := make(map[int]string)
	//var actid []int
	headers := getbaseheaders()
	f.WriteString("-----" + time.Now().String() + "-----\r\n")
	println("-----" + time.Now().String() + "-----")
	f.WriteString("共 " + strconv.Itoa(len(actid)) + "条霸王餐，开始报名 \r\n")
	println("共 ", len(actid), "条霸王餐，开始报名")
	for _, id := range actid {
		client := &http.Client{}
		re, err := http.NewRequest("GET", "http://m.dianping.com/astro-plat/freemeal/bwcDetail?offlineActivityId="+strconv.Itoa(id), nil)
		re = addhead(re, headers)
		resp, err := client.Do(re)
		if err != nil {
			println(err.Error())
		}
		resp = checkgzip(resp)
		content, _ := ioutil.ReadAll(resp.Body)
		respBody := string(content)
		brid := strings.Split(strings.Split(respBody, string("\"dpShopId\":"))[1], ",")[0]

		client = &http.Client{}
		if brid == "" {
			return
		}
		data := map[string]interface{}{
			"offlineActivityId": strconv.Itoa(id),
			"branchId":          brid,
		}
		bytesData, _ := json.Marshal(data)
		reader := bytes.NewReader(bytesData)
		r, err := http.NewRequest("POST", "http://m.dianping.com/mobile/dinendish/apply/doApplyActivity", ioutil.NopCloser(reader))
		if err != nil {
			println(err.Error())
		}
		addhead(r, headers)
		r.Header.Add("Cookie", "dper="+con.Dper)
		res, err := client.Do(r)
		if err != nil {
			println(err.Error())
		}
		var result enroolResult
		res = checkgzip(res)
		content, _ = ioutil.ReadAll(res.Body)
		//respBody := string(content)
		json.Unmarshal(content, &result)
		s := "活动名称 ：" + act[id] + " 活动ID： " + strconv.Itoa(id) + " 活动状态： " + result.Data.Desc + "\r\n"
		f.WriteString(s)

		println("活动名称 ：", act[id], " 活动ID： ", id, " 活动状态： ", result.Data.Desc)
		if result.Data.Desc == "服务忙，请重试" {
			time.Sleep(3 * time.Second)
		}
		if result.Data.Desc == "请先登录" {
			println("未登录或登录失败")
			break
		}
		time.Sleep(time.Duration(con.Internal) * time.Second)
	}
	println("--------------结束--------------------\r\n")
	f.WriteString("--------------结束--------------------\r\n")
}

func addhead(r *http.Request, a map[string]string) *http.Request {
	for k, v := range a {
		r.Header.Add(k, v)
	}
	return r
}

func checkgzip(resp *http.Response) *http.Response {
	var e error
	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, e = gzip.NewReader(resp.Body)
		if e != nil {
			println(e.Error())
		}
	}
	return resp
}
