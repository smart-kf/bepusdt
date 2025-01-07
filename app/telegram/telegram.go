package telegram

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/v03413/bepusdt/app/config"
)

var (
	botApi *tgbotapi.BotAPI
	err    error
)

func InitBot(token string) {
	botApi, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic("TG Bot NewBotAPI Error:" + err.Error())

		return
	}

	// 注册命令
	_, err = botApi.Request(
		tgbotapi.NewSetMyCommands(
			[]tgbotapi.BotCommand{
				{Command: "/" + cmdGetId, Description: "获取ID"},
				{Command: "/" + cmdStart, Description: "开始使用"},
				{Command: "/" + cmdUsdt, Description: "实时汇率"},
				{Command: "/" + cmdWallet, Description: "钱包信息"},
				{Command: "/" + cmdOrder, Description: "最近订单"},
			}...,
		),
	)
	if err != nil {
		panic("TG Bot Request Error:" + err.Error())

		return
	}

	fmt.Println("Bot UserName: ", botApi.Self.UserName)
}

func GetBotApi() *tgbotapi.BotAPI {
	return botApi
}

func SendMsg(msg tgbotapi.MessageConfig) {
	bot := config.GetTgBot()
	if !bot.Enable {
		return
	}
	if msg.ChatID != 0 {
		_, _ = botApi.Send(msg)

		return
	}

	botConfig := config.GetTgBot()

	chatId, err := strconv.ParseInt(botConfig.AdminId, 10, 64)
	if err == nil {
		msg.ChatID = chatId
		_, _ = botApi.Send(msg)
	}
}

func DeleteMsg(msgId int) {
	botConfig := config.GetTgBot()
	if !botConfig.Enable {
		return
	}
	chatId, err := strconv.ParseInt(botConfig.AdminId, 10, 64)
	if err == nil {
		_, _ = botApi.Send(tgbotapi.NewDeleteMessage(chatId, msgId))
	}
}

func EditAndSendMsg(msgId int, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) {
	botConfig := config.GetTgBot()
	if !botConfig.Enable {
		return
	}
	chatId, err := strconv.ParseInt(botConfig.AdminId, 10, 64)
	if err == nil {
		_, _ = botApi.Send(tgbotapi.NewEditMessageTextAndMarkup(chatId, msgId, text, replyMarkup))
	}
}
