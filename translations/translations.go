package translations

import (
	"sync"

	"github.com/OzodbekX/TuronMiniApp/volumes"
)

var (
	translations = map[string]map[string]string{
		"ru": {
			"welcome":             "Добро пожаловать! Пожалуйста, выберите язык.",
			"login":               "Пожалуйста, введите ваш логин.",
			"cancel":              "Отмена",
			"mainMenu":            "Главное меню",
			"enterPassword":       "Введите ваш пароль:",
			"wrongParol":          "Неправильный логин или пароль, пожалуйста, введите логин еще раз",
			"Tariffs":             "Тарифы",
			"FAQ":                 "FAQ",
			"Application":         "Заявление",
			"Language":            "Язык",
			"Balance":             "Баланс",
			"Exit":                "Выход",
			"PleaseSelectOption":  "Илтимос, менюдан бўлимни танланг:",
			"pleaseEnterYourName": "Пожалуйста, введите свое имя:",
			"enterPhone":          "Введите номер телефона в следующем формате: +998#########",
		},
		"uz": {
			"welcome":             "Xush kelibsiz! Iltimos, tilni tanlang.",
			"login":               "Iltimos loginingizni kiriting",
			"cancel":              "Bekor qilish",
			"mainMenu":            "Asosiy menyu",
			"enterPassword":       "Iltimos, parolingizni kiriting:",
			"wrongParol":          "Noto'g'ri login yoki parol, iltimos loginni qayta kiriting",
			"Tariffs":             "Tariflar",
			"FAQ":                 "FAQ",
			"Application":         "Ariza",
			"Language":            "Til",
			"Balance":             "Balance",
			"Exit":                "Chiqish",
			"PleaseSelectOption":  "Iltimos, menyudan bo‘limni tanlang:",
			"pleaseEnterYourName": "Iltimos, ismingizni kiriting:",
			"enterPhone":          "Iltimos, telefon raqamingizni quyidagi formatda kiriting: +998#########",
		},
	}
)

func GetTranslation(userSessions *sync.Map, chatID int64, key string) string {
	lang := "uz"
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		lang = user.Language
	}
	if text, ok := translations[lang][key]; ok {
		return text
	}
	return key // Fallback to key if translation is missing
}
