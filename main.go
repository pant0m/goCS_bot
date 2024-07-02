package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 定义 API URL 常量
const (
	wechatBotURL   = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"
	dingDingBotURL = "https://oapi.dingtalk.com/robot/send"
)

// WeChatBody 结构体代表微信消息体的结构
type WeChatBody struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
}

// Text 结构体代表消息体中的文本内容
type Text struct {
	Content string `json:"content"`
}

// DingDingBody 结构体代表钉钉消息体的结构
type DingDingBody struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

// At 结构体代表钉钉消息中的提及信息
type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

// wechatBot 函数向微信机器人发送消息
func wechatBot(news, wechatBotKey string) {
	body := WeChatBody{
		MsgType: "text",
		Text:    Text{Content: news},
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("JSON编码失败: %v", err)
	}

	response, err := http.Post(fmt.Sprintf("%s?key=%s", wechatBotURL, wechatBotKey), "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Fatalf("发送消息到微信机器人失败: %v", err)
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("读取微信响应失败: %v", err)
	}

	log.Printf("微信状态码: %d\n", response.StatusCode)
	log.Printf("微信响应内容: %s\n", respBody)
}

// dindinBot 函数向钉钉机器人发送消息
func dindinBot(news, dindinBotKey string) {
	body := DingDingBody{
		MsgType: "text",
		Text:    Text{Content: news},
		At:      At{AtMobiles: []string{"@all"}, IsAtAll: true},
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("JSON编码失败: %v", err)
	}

	response, err := http.Post(fmt.Sprintf("%s?access_token=%s", dingDingBotURL, dindinBotKey), "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Fatalf("发送消息到钉钉机器人失败: %v", err)
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("读取钉钉响应失败: %v", err)
	}

	log.Printf("钉钉状态码: %d\n", response.StatusCode)
	log.Printf("钉钉响应内容: %s\n", respBody)
}

// CheckIPchina 函数检查提供的 IP 是否来自中国
func CheckIPchina(IPaddr string) bool {
	url := fmt.Sprintf("https://searchplugin.csdn.net/api/v1/ip/get?ip=%s", IPaddr)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("获取数据失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应体失败: %v", err)
	}

	htmlContent := string(body)
	return strings.Contains(htmlContent, "中国")
}

// generateAlertMessage 函数生成警报消息内容
func generateAlertMessage(computerName, internalIP, externalIP, userName, proCess *string) string {

	if !CheckIPchina(*externalIP) {
		log.Printf("外网 IP (%s) 不是来自中国，程序退出。\n", *externalIP)
		os.Exit(1)
	}

	content := fmt.Sprintf("新主机上线！\n主机名：%s\n内网 IP：%s\n外网 IP：%s\n用户名：%s\n进程：%s\n请及时处理！@所有人",
		*computerName, *internalIP, *externalIP, *userName, *proCess)

	return content
}

func main() {

	// 解析命令行参数
	computerName := flag.String("computername", "", "计算机名称")
	internalIP := flag.String("internalip", "", "内网 IP")
	externalIP := flag.String("externalip", "", "外网 IP")
	userName := flag.String("username", "", "用户名")
	proCess := flag.String("process", "", "进程")

	flag.Parse()

	wechatBotKey := ""
	dingdingBotKey := ""

	content := generateAlertMessage(computerName, internalIP, externalIP, userName, proCess)
	wechatBot(content, wechatBotKey)
	dindinBot(content, dingdingBotKey)
	log.Printf("警报消息内容:\n%s\n", content)
}
