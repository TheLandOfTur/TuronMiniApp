package conversations

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func writeUserApplications(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map, applications []volumes.UserApplications) {
	for _, app := range applications {
		t, _ := time.Parse("2006-01-02 15:04:05", app.CreatedAt)

		response := fmt.Sprintf(
			"📄 *%s*\n\n"+
				"👤 *%s*: %s\n"+
				"📞 *%s*: %s\n"+
				"🌆 *%s*: %s\n"+
				"🏘 *%s*: %s\n"+
				//"📂 *%s*: %s\n"+
				"📅 *%s*: %s\n",
			//"📈 *%s*: %s\n",
			translations.GetTranslation(userSessions, chatID, "applicationTitle"),
			translations.GetTranslation(userSessions, chatID, "fullName"), app.FullName,
			translations.GetTranslation(userSessions, chatID, "phoneNumber"), app.TelegramPhoneNumber,
			translations.GetTranslation(userSessions, chatID, "city"), app.CityName,
			translations.GetTranslation(userSessions, chatID, "district"), app.DistrictName,
			//translations.GetTranslation(userSessions, chatID, "applicationType"), app.RequestCategory,
			translations.GetTranslation(userSessions, chatID, "createDate"), t.Format("2006.01.02 15:04"),
			//translations.GetTranslation(userSessions, chatID, "statusApplication"), app.RequestStatus,
		)

		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown" // ✅ make titles bold
		bot.Send(msg)
	}
}

func handleSuccessfulMessageState(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	msg := update.Message
	text := strings.TrimSpace(msg.Text)
	sessionData, ok := userSessions.Load(chatID)
	if !ok {
		log.Printf("[WARN] No session found for chatID: %d", chatID)
		return
	}
	user := sessionData.(*volumes.UserSession)
	switch text {
	case translations.GetTranslation(userSessions, chatID, "exit1"):
		langKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "abonent")),
				tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "user")),
			),
		)

		selectUserTypeMessage := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "chooseRole"))
		selectUserTypeMessage.ReplyMarkup = langKeyboard
		bot.Send(selectUserTypeMessage)
		user.State = volumes.CHOOSE_USER_TYPE
		return

	case translations.GetTranslation(userSessions, chatID, "myApplications"):
		// 📋 Show user's applications

		applications, err := server.MyApplications(user, update.Message.From.ID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "❌ "+translations.GetTranslation(userSessions, chatID, "errorFetchingApplications")))
			return
		}
		log.Println("applications", applications)

		if len(applications) == 0 {
			bot.Send(tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "noApplications")))
			return
		}
		writeUserApplications(bot, chatID, userSessions, applications)

		// State remains SUCCESSFUL_MESSAGE
		return

	default:
		// Unknown input — optional fallback
		bot.Send(tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "unknownCommand")))
	}
}
