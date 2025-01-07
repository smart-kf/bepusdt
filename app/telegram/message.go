package telegram

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/v03413/bepusdt/app/config"
	"github.com/v03413/bepusdt/app/help"
	"github.com/v03413/bepusdt/app/model"
)

func getNotifyId() (string, bool) {
	var targetId string
	botConfig := config.GetTgBot()
	if botConfig.GroupId != "" {
		targetId = botConfig.GroupId
	}
	if targetId == "" && botConfig.AdminId != "" {
		targetId = botConfig.AdminId
	}
	if targetId == "" {
		return "", false
	}
	return targetId, true
}
func SendTradeSuccMsg(order model.TradeOrders) {
	var targetId, ok = getNotifyId()
	if !ok {
		return
	}
	chatId, err := strconv.ParseInt(targetId, 10, 64)
	if err != nil {
		return
	}

	tradeType := "USDT"
	tradeUnit := `USDT.TRC20`
	if order.TradeType == model.OrderTradeTypeTronTrx {
		tradeType = "TRX"
		tradeUnit = "TRX"
	}

	text := `
#收款成功 #订单交易 #` + tradeType + `
---
` + "```" + `
🚦商户订单：%v
💰请求金额：%v CNY(%v)
💲支付数额：%v ` + tradeUnit + `
✅收款地址：%s
⏱️创建时间：%s
️🎯️支付时间：%s
` + "```" + `
`
	text = fmt.Sprintf(
		text,
		order.OrderId,
		order.Money,
		order.TradeRate,
		order.Amount,
		help.MaskAddress(order.Address),
		order.CreatedAt.Format(time.DateTime),
		order.UpdatedAt.Format(time.DateTime),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonURL(
					"📝查看交易明细",
					"https://tronscan.org/#/transaction/"+order.TradeHash,
				),
			},
		},
	}

	_, _ = botApi.Send(msg)
}

func SendOtherNotify(text string) {
	var targetId, ok = getNotifyId()
	if !ok {
		return
	}
	chatId, err := strconv.ParseInt(targetId, 10, 64)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, _ = botApi.Send(msg)
}

func SendWelcome(version string) {
	text := `
👋 如果您看到此消息，说明机器人已经启动成功

📌当前版本：` + version + `
📝发送命令 /start 可以开始使用
---
`
	msg := tgbotapi.NewMessage(0, text)

	SendMsg(msg)
}
