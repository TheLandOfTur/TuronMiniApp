package conversations

import (
	"sync"

	"github.com/OzodbekX/TuronMiniApp/handlers/chat"

	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var lastMessageIDs sync.Map // To track the last message sent by the bot

func HandleUpdateConversation(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	session, _ := userSessions.LoadOrStore(chatID, &volumes.UserSession{State: volumes.SELECT_LANGUAGE})
	user := session.(*volumes.UserSession)
	switch user.State {
	case volumes.SELECT_LANGUAGE:
		handleLanguage(bot, update, userSessions)
	case volumes.LOGIN:
		handleLogin(bot, update, userSessions)
	case volumes.PASSWORD:
		handlePassword(bot, update, userSessions)
	case volumes.CHOOSE_USER_TYPE:
		identifyUserType(bot, update, userSessions)
	case volumes.SUBMIT_PHONE:
		handlePhoneNumber(bot, update, userSessions)
	case volumes.LOG_OUT:
		handleLogOut(bot, update, userSessions)
	case volumes.CHANGE_LANGUAGE:
		onchangeLanguage(bot, update, userSessions)
	case volumes.ACTIVATE_PROMOCODE:
		checkActivePromoCode(bot, update, userSessions)

	//case volumes.ENTER_FULL_NAME:
	//	handleFullNameInput(bot, update, userSessions)
	case volumes.CHOOSE_REGIONS:
		if update.Message.Text != "" {
			handleRegionWrite(bot, update, userSessions)
		}
	case volumes.CHOOSE_DISTRICTS:
		if update.Message.Text != "" {
			handeDistrictWrite(bot, update, userSessions)
		}
	case volumes.SELECT_CATEGORY, volumes.SELECT_FAQ:
		chat.HandleChatConversation(bot, update, userSessions, user)
	case volumes.USER_CABINET:
		handleSuccessfulMessageState(bot, update, userSessions)
	}
}
