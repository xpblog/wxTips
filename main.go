package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	qrcode "github.com/skip2/go-qrcode"
	"time"
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
