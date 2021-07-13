package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	dper string
	menu []int
	cityname string
}

type enroolResult struct{
	Code int `json:"code"`
	Data Da  `json:"data"`
}
type Da struct {
	Desc string `json:"desc"`
}

type list struct {
	Code  int   	`json:"code"`
	Data   Data     `json:"data"`
}
type Data struct {
	HasNext bool  `json:"hasNext"`
	Details  []Detail  `json:"detail"`
}

type Detail struct {
	OfflineActivityId int `json:"offlineActivityId"`
	ActivityTitle string `json:"activityTitle"`
	DetailUrl     string `json:"detailUrl"`
}

const MT = "luke_tb"
var Con  Config
func main() {
	f ,_ := os.OpenFile("log.txt" , os.O_APPEND,0644)
	defer f.Close()
	makeconfig()

	headers := map[string]string{
		"Origin": "http://s.dianping.com",
		"Accept-Encoding": "gzip, deflate",
		"X-Request": "JSON",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
		"Content-Type": "application/json;charset=UTF-8",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept": "application/json, text/javascript",
		"Referer": "http://s.dianping.com/event/1063463422",
		"X-Requested-With": "XMLHttpRequest",
		"Connection": "keep-alive",
	}

	cityId := getcityid(headers)
	act := make(map[int]string)
	var actid []int
	dper := Con.dper
	types_food := []int{1}
	for i := range types_food {
		for page := 1 ; page<=15;page++{
			data := map[string]interface{}{
				"cityId" : cityId,
				"type":i,
				"mode": "" ,
				"page": page,
			}
			bytesData, _ := json.Marshal(data)
			reader := bytes.NewReader(bytesData)

			client := &http.Client{}
			r ,err := http.NewRequest("POST" , "http://m.dianping.com/activity/static/pc/ajaxList",ioutil.NopCloser(reader))
			if err !=nil{
				println(err.Error())
			}
			r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
			r.Header.Set("Referer","http://s.dianping.com/event/shanghai")
			r.Header.Set("Origin" ,"http://s.dianping.com")
			r.Header.Set("User-Agent" , "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
			r.Header.Set("Accept-Encoding","gzip, deflate")
			r.Header.Set("Content-Type", "application/json;charset=UTF-8")
			r.Header.Set("Accept-Language","zh-CN,zh;q=0.9")
			rs , err := client.Do(r)
			if err != nil {
				println(err.Error())
			}
			rs = checkgzip(rs)
			l := &list{}
			content, _ := ioutil.ReadAll(rs.Body)
			//string(content)
			json.Unmarshal(content, l)
			num := len(l.Data.Details)
			for i := 0;i<num;i++{
				actid = append(actid, l.Data.Details[i].OfflineActivityId)
				act[l.Data.Details[i].OfflineActivityId]=l.Data.Details[i].ActivityTitle
			}
			if l.Data.HasNext == false{
				break
			}
		}
	}

	//报名


	f.WriteString("-----"+time.Now().String()+"-----\r\n")
	println("-----"+time.Now().String()+"-----")
	f.WriteString("共 " + strconv.Itoa(len(actid))+"条把霸王餐，开始报名 \r\n")
	println("共 " , len(actid),"条把霸王餐，开始报名")
	for _ , id := range actid{
		client := &http.Client{}
		re ,err := http.NewRequest("GET" , "http://m.dianping.com/astro-plat/freemeal/bwcDetail?offlineActivityId="+strconv.Itoa(id) ,nil)
		re = addhead(re , headers)
		resp ,err:= client.Do(re)
		if err != nil {
			println(err.Error())
		}
		resp = checkgzip(resp)
		content, _ := ioutil.ReadAll(resp.Body)
		respBody := string(content)
		brid := strings.Split(strings.Split(respBody , string("\"dpShopId\":"))[1],",")[0]

		client = &http.Client{}
		if brid == ""  {
			return
		}
		data := map[string]interface{}{
			"offlineActivityId": strconv.Itoa(id),
			"branchId":brid,
		}
		bytesData, _ := json.Marshal(data)
		reader := bytes.NewReader(bytesData)
		r , err := http.NewRequest("POST" , "http://m.dianping.com/mobile/dinendish/apply/doApplyActivity" ,ioutil.NopCloser(reader))
		if err != nil {
			println(err.Error())
		}
		addhead(r,headers)
		r.Header.Add("Cookie" ,"dper="+dper)
		res ,err:= client.Do(r)
		if err != nil {
			println(err.Error())
		}
		var result enroolResult
		res = checkgzip(res)
		content, _ = ioutil.ReadAll(res.Body)
		//respBody := string(content)
		json.Unmarshal(content ,&result )
		s := "活动名称 ：" + act[id]+" 活动ID： " + strconv.Itoa(id) +" 活动状态： " + result.Data.Desc +"\r\n"
		f.WriteString(s)

		println("活动名称 ：" , act[id]," 活动ID： " , id ," 活动状态： " , result.Data.Desc)
		if result.Data.Desc == "服务忙，请重试"{
			time.Sleep(3*time.Second)
		}
	}
	f.WriteString("--------------结束--------------------\r\n\r\n\r\n")
}

func addhead(r *http.Request , a map[string]string) *http.Request{
	for k,v := range a {
		r.Header.Add(k, v)
	}
	return r
}

func makeconfig() {

	L := lua.NewState()

	conf := L.CreateTable(0, 32)

	L.SetGlobal("luke", conf)

	luke := L.GetGlobal("luke").(lua.LValue)
	ud := L.NewUserData()
	ud.Value = &Con

	mt := L.NewTypeMetatable(MT)
	L.SetField(mt, "__index", L.NewFunction(Get))
	L.SetField(mt, "__newindex", L.NewFunction(Set))
	L.SetMetatable(ud, mt)

	L.SetField(luke, "config", ud)

	L.DoFile("config.lua")
}

func Get(L *lua.LState) int {
	return 1
}

func Set(L *lua.LState) int {
	name := L.CheckString(2)
	//println(name)
	switch name {
	case "dper":
		//println(L.CheckString(3))
		Con.dper = L.CheckString(3)
	case "cityname":
		Con.cityname = L.CheckString(3)
	case "menu":
		menu := L.CheckString(3)
		me := strings.Split(menu , ",")
		for i := range me{
			//println(me[i])
			a , _ :=strconv.Atoi(me[i])
			Con.menu = append(Con.menu,a)
		}
	}
	return 0
}

func getcityid(headers map[string]string) int {
	//ur , _ := url2.Parse("http://127.0.0.1:10809")
	city := Con.cityname
	client := &http.Client{
		Transport: &http.Transport{
		},
	}

	re , err:= http.NewRequest("GET" , "http://s.dianping.com/event/"+city , nil)
	if err != nil {
		println(err.Error())
	}
	addhead(re , headers)
	resp , err := client.Do(re)
	if err != nil {
		println(err.Error())
	}
	resp = checkgzip(resp)

	content , err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()
	println(resp.Status)
	//println(string(content))
	id := strings.Split(strings.Split(string(content) ,"\"cityId\":\"")[1] , "\"")[0]
	cityid , err := strconv.Atoi(id)
	println(cityid)
	return cityid
}

func checkgzip(resp *http.Response) *http.Response{
	var e error
	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, e = gzip.NewReader(resp.Body)
		if e!= nil {
			println(e.Error())
		}
	}
	return resp
}