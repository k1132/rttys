package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func PushPlusGet(url string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func PushPlusPost(url string, data interface{}, contentType string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}

func PushPlusMsg(br *broker,title string, devid_str string) {
	//var title = "go_test"
	//var devid_str = "gotestcontent"
	cfg := br.cfg
	//连接字符串
	if cfg.PushString != "" {
		if cfg.PushId != "" {
			if cfg.PushId == devid_str {
				connString := fmt.Sprintf("%s%s-%s", cfg.PushString,title, devid_str)
				connString = strings.Replace(connString, " ", "", -1)
				fmt.Println(connString)
				PushPlusGet(connString)
			}
		} else {
			connString := fmt.Sprintf("%s%s-%s", cfg.PushString,title, devid_str)
			connString = strings.Replace(connString, " ", "", -1)
			fmt.Println(connString)
			PushPlusGet(connString)
		}
	}

}
