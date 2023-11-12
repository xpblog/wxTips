package main

import (
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	qrcode "github.com/skip2/go-qrcode"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wxTips/weather"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = MsgHandler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl //ConsoleQrCode

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 定时任务
	Timer(self)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

// 二维码打印到控制台
func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

// 定时提醒，
// 每周五5：30提醒写周报；每月5号(非周末)提醒批考勤
func Timer(self *openwechat.Self) {
	go func(self *openwechat.Self) {
		fmt.Print("开始执行定时。。。。。")

		for true {
			select {
			case t := <-time.After(1 * time.Minute):

				// 凌晨1点获取天气，如果天气异常，则提醒
				weatherTime := t.Format("15:04")
				if weatherTime == "01:00" {
					isSend, msg := getWeather("101010100")
					if isSend {
						friends, _ := self.Friends()
						friend := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.NickName == "鹏程" })
						friend.SendText(msg)
					}
				}

				weekday := t.Weekday()
				// 周五
				if weekday == 5 {
					nowTime := t.Format("15:04")
					if nowTime == "17:30" {
						friends, _ := self.Friends()
						friend := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.NickName == "鹏程" })
						friend.SendText("公主请写周报...")
					}
				}

				day := t.Day()
				if (day == 5 && (weekday != 6 || weekday != 0)) ||
					(day == 6 && weekday == 1) ||
					(day == 7 && weekday == 2) {
					friends, _ := self.Friends()
					friend := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.NickName == "鹏程" })
					friend.SendText("公主请提醒批考勤...")
				}
			}
		}
	}(self)
}

func MsgHandler(msg *openwechat.Message) {
	if msg.IsText() && msg.Content == "ping" {
		msg.ReplyText("ping")
	}
}

func getWeather(addrCode string) (bool, string) {
	url := "http://t.weather.itboy.net/api/weather/city/" + addrCode
	resp, _ := http.Get(url)
	text, _ := ioutil.ReadAll(resp.Body)

	var weather weather.WeatherRes
	json.Unmarshal(text, &weather)

	today_weather := weather.Data.Forecast[0]
	fmt.Println(today_weather)
	high, low, type_, fl :=
		strings.Trim(strings.ReplaceAll(strings.ReplaceAll(today_weather.High, "高温", ""), "℃", ""), " "),
		strings.Trim(strings.ReplaceAll(strings.ReplaceAll(today_weather.Low, "低温", ""), "℃", ""), " "),
		today_weather.Type,
		strings.Trim(strings.ReplaceAll(today_weather.Fl, "级", ""), " ")

	var msg = ""
	highi, _ := strconv.Atoi(high)
	if highi > 30 {
		msg = "天气超过30度了"
	}
	lowi, _ := strconv.Atoi(low)
	if lowi < 0 {
		msg = msg + "那天气是相当的冷"
	}
	if strings.Contains(type_, "雨") {
		msg = msg + " 有雨，带伞"
	}
	if strings.Contains(type_, "雪") {
		msg = msg + " 有雪，穿厚衣"
	}
	fli, _ := strconv.Atoi(fl)
	if fli > 5 {
		msg = msg + " 风超级大的"
	}

	if len(msg) > 0 {
		return true, msg
	}
	return false, ""
}
