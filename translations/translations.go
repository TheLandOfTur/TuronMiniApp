package translations

import (
	"sync"

	"github.com/OzodbekX/TuronMiniApp/volumes"
)

var (
	translations = map[string]map[string]string{
		"ru": {
			"welcome":              "Добро пожаловать! Пожалуйста, выберите язык.",
			"login":                "Пожалуйста, введите ваш логин.",
			"cancel":               "Отмена",
			"mainMenu":             "Главное меню",
			"enterPassword":        "Введите ваш пароль:",
			"wrongParol":           "Неправильный логин или пароль, пожалуйста, введите логин еще раз",
			"Tariffs":              "Тарифы",
			"FAQ":                  "FAQ",
			"Application":          "Заявление",
			"Language":             "Язык",
			"Balance":              "Баланс",
			"Exit":                 "Выход",
			"PleaseSelectOption":   "Илтимос, менюдан бўлимни танланг:",
			"pleaseEnterYourName":  "Пожалуйста, введите свое имя:",
			"enterPhone":           "Введите номер телефона в следующем формате: +998#########",
			"listOfTariffs":        "Список тарифов",
			"price":                "Цена",
			"speedByTime":          "Скорость по времени",
			"mbs":                  "Мбит/с",
			"uzs":                  "uzs",
			"sharePhonenumber":     "Поделиться номером телефона",
			"invalidPhoneNumber":   "Неверный формат номера телефона. Укажите действительный номер: +998######### или #########.",
			"loginSuccessful":      "Авторизация успешна! Добро пожаловать в главное меню.",
			"yourBalance":          "Ваш баланс",
			"tariffName":           "Тарифное название",
			"subscriptionPrice":    "Цена подписки",
			"nextSubscriptionDate": "Дата следующей подписки",
			"subscriptionPeriod":   "Период подписки",
			"from":                 "с",
			"to":                   "по",
			"subscriptionActive":   "Подписка активна",
			"active":               "Активно",   // Active in Russian
			"inactive":             "Неактивно", // Inactive in Russian

		},
		"uz": {
			"welcome":              "Xush kelibsiz! Iltimos, tilni tanlang.",
			"login":                "Iltimos loginingizni kiriting",
			"cancel":               "Bekor qilish",
			"mainMenu":             "Asosiy menyu",
			"enterPassword":        "Iltimos, parolingizni kiriting:",
			"wrongParol":           "Noto'g'ri login yoki parol, iltimos loginni qayta kiriting",
			"Tariffs":              "Tariflar",
			"FAQ":                  "FAQ",
			"Application":          "Ariza",
			"Language":             "Til",
			"Balance":              "Balance",
			"Exit":                 "Chiqish",
			"PleaseSelectOption":   "Iltimos, menyudan bo‘limni tanlang:",
			"pleaseEnterYourName":  "Iltimos, ismingizni kiriting:",
			"enterPhone":           "Iltimos, telefon raqamingizni quyidagi formatda kiriting: +998#########",
			"listOfTariffs":        "Tariflar ro'yxati",
			"price":                "Narxi",
			"speedByTime":          "Vaqt bo'yicha tezlik",
			"mbs":                  "Mbit/s",
			"uzs":                  "uzs",
			"sharePhonenumber":     "Telefon raqamini ulashing",
			"invalidPhoneNumber":   "Telefon raqami formati noto‘g‘ri. Iltimos, haqiqiy raqamni kiriting: +998######### yoki #########.",
			"loginSuccessful":      "Kirish muvaffaqiyatli! Bosh menyuga xush kelibsiz.",
			"yourBalance":          "Sizning balansingiz",
			"tariffName":           "Tarif nomi",
			"subscriptionPrice":    "Obuna narxi",
			"nextSubscriptionDate": "Keyingi obuna sanasi",
			"subscriptionPeriod":   "Obuna davri",
			"from":                 "dan",
			"to":                   "gacha",
			"subscriptionActive":   "Obuna faolligi",
			"active":               "Faol",    // Active in Uzbek
			"inactive":             "No Faol", // Inactive in Uzbek

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
