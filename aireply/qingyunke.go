package aireply

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// QYKReply 青云客回复类
type QYKReply struct{}

const (
	qykURL     = "http://api.qingyunke.com/api.php?key=free&appid=0&msg=%s"
	qykBotName = "菲菲"
)

var (
	qykMatchFace = regexp.MustCompile(`\{face:(\d+)\}(.*)`)
)

func (*QYKReply) String() string {
	return "青云客"
}

// Talk 取得回复消息
func (*QYKReply) Talk(msg string) message.Message {
	msg = strings.ReplaceAll(msg, zero.BotConfig.NickName[0], qykBotName)

	u := fmt.Sprintf(qykURL, url.QueryEscape(msg))
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return message.Message{message.Text("ERROR:", err)}
	}
	// 自定义Header
	req.Header.Set("User-Agent", web.RandUA())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.qingyunke.com")
	resp, err := client.Do(req)
	if err != nil {
		return message.Message{message.Text("ERROR:", err)}
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return message.Message{message.Text("ERROR:", err)}
	}

	replystr := gjson.Get(helper.BytesToString(bytes), "content").String()
	replystr = strings.ReplaceAll(replystr, "{face:", "[CQ:face,id=")
	replystr = strings.ReplaceAll(replystr, "{br}", "\n")
	replystr = strings.ReplaceAll(replystr, "}", "]")
	replystr = strings.ReplaceAll(replystr, qykBotName, zero.BotConfig.NickName[0])

	return message.ParseMessageFromString(replystr)
}

// TalkPlain 取得回复消息
func (*QYKReply) TalkPlain(msg string) string {
	msg = strings.ReplaceAll(msg, zero.BotConfig.NickName[0], qykBotName)

	u := fmt.Sprintf(qykURL, url.QueryEscape(msg))
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	// 自定义Header
	req.Header.Set("User-Agent", web.RandUA())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.qingyunke.com")
	resp, err := client.Do(req)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "ERROR: " + err.Error()
	}

	replystr := gjson.Get(helper.BytesToString(bytes), "content").String()
	replystr = qykMatchFace.ReplaceAllLiteralString(replystr, "")
	replystr = strings.ReplaceAll(replystr, "{br}", "\n")
	replystr = strings.ReplaceAll(replystr, qykBotName, zero.BotConfig.NickName[0])

	return replystr
}
