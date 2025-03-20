package translations

import (
	"sync"

	"github.com/OzodbekX/TuronMiniApp/volumes"
)

var (
	translations = map[string]map[string]string{
		"ru": {
			"welcome":                    "Добро пожаловать! Пожалуйста, выберите язык.",
			"login":                      "Пожалуйста, введите ваш логин.",
			"cancel":                     "⬅️ Отмена",
			"mainMenu":                   "🏠 Главное меню",
			"enterPassword":              "Введите ваш пароль:",
			"wrongParol":                 "Неправильный логин или пароль, пожалуйста, введите логин еще раз",
			"Tariffs":                    "Тарифы",
			"FAQ":                        "FAQ",
			"Application":                "Заявление",
			"Language":                   "Язык",
			"Balance":                    "Баланс",
			"Exit":                       "Выход",
			"PleaseSelectOption":         "Пожалуйста, выберите раздел из меню:",
			"pleaseEnterYourName":        "Пожалуйста, введите свое имя:",
			"enterPhone":                 "Введите номер телефона в следующем формате: +998#########",
			"pleaseShareYourPhoneNumber": "Пожалуйста, поделитесь номером телефона",
			"listOfTariffs":              "Список тарифов",
			"price":                      "Цена",
			"speedByTime":                "Скорость по времени",
			"mbs":                        "Мбит/с",
			"uzs":                        "сум",
			"sharePhoneNumber":           "Поделиться номером телефона",
			"invalidPhoneNumber":         "Неверный формат номера телефона. Укажите действительный номер: +998######### или #########.",
			"loginSuccessful":            "Авторизация успешна! Добро пожаловать в главное меню.",
			"yourBalance":                "Ваш баланс",
			"tariffName":                 "Ваш тариф",
			"subscriptionPrice":          "Цена подписки",
			"nextSubscriptionDate":       "Дата следующей подписки",
			"subscriptionPeriod":         "Период подписки",
			"from":                       "с",
			"to":                         "по",
			"subscriptionActive":         "Подписка активна",
			"active":                     "Активно",   // Active in Russian
			"inactive":                   "Неактивно", // Inactive in Russian
			"pleaseSelectCategory":       "Пожалуйста, выберите категорию:",
			"pleaseSelectFAQ":            "Пожалуйста, выберите FAQ",
			"yes":                        "✅ Да",
			"no":                         "❌ Нет",
			"doYouWantLogout":            "Вы хотите выйти из системы?",
			"promoCode":                  "Промо-код",
			"enterCode":                  "Введите код",
			"promoCodeInactive":          "Промокод неактивен",
			"promoCodeActive":            "Промокод активен",
			"promoCodeAlreadyActivated":  "Промокод уже активирован",
			"promoCodePermissionDenied":  "Доступ к промокоду запрещен",
			"promoCodeNotFound":          "Промокод не найден",
			"status":                     "Статус",
			"connectWithOperator":        "Связаться с оператором",
			"operatorMessage":            "<a href='https://t.me/turonsupport'>Оператор</a>, свяжитесь и получите ответы на свои вопросы.",
		},
		"uz": {
			"welcome":                    "Xush kelibsiz! Iltimos, tilni tanlang.",
			"login":                      "Iltimos loginingizni kiriting",
			"cancel":                     "⬅️ Bekor qilish",
			"mainMenu":                   "🏠 Asosiy menyu",
			"enterPassword":              "Iltimos, parolingizni kiriting:",
			"wrongParol":                 "Noto'g'ri login yoki parol, iltimos loginni qayta kiriting",
			"Tariffs":                    "Tariflar",
			"FAQ":                        "FAQ",
			"Application":                "Ariza",
			"Language":                   "Til",
			"Balance":                    "Balans",
			"Exit":                       "Chiqish",
			"PleaseSelectOption":         "Iltimos, menyudan bo‘limni tanlang:",
			"pleaseEnterYourName":        "Iltimos, ismingizni kiriting:",
			"enterPhone":                 "Iltimos, telefon raqamingizni quyidagi formatda kiriting: +998#########",
			"pleaseShareYourPhoneNumber": "Iltimos, telefon raqamingizni ulashing",
			"listOfTariffs":              "Tariflar ro'yxati",
			"price":                      "Narxi",
			"speedByTime":                "Vaqt bo'yicha tezlik",
			"mbs":                        "Mbit/s",
			"uzs":                        "so'm",
			"sharePhoneNumber":           "Telefon raqamini ulashing",
			"invalidPhoneNumber":         "Telefon raqami formati noto‘g‘ri. Iltimos, haqiqiy raqamni kiriting: +998######### yoki #########.",
			"loginSuccessful":            "Kirish muvaffaqiyatli! Bosh menyuga xush kelibsiz.",
			"yourBalance":                "Sizning balansingiz",
			"tariffName":                 "Sizning tarifingiz",
			"subscriptionPrice":          "Obuna narxi",
			"nextSubscriptionDate":       "Keyingi obuna sanasi",
			"subscriptionPeriod":         "Obuna davri",
			"from":                       "dan",
			"to":                         "gacha",
			"subscriptionActive":         "Obuna faolligi",
			"active":                     "Faol",    // Active in Uzbek
			"inactive":                   "No Faol", // Inactive in Uzbek
			"pleaseSelectCategory":       "Kategoriya tanlang:",
			"pleaseSelectFAQ":            "Iltimos FAQni tanlang",
			"yes":                        "✅ Ha",
			"no":                         "❌ Yo'q",
			"doYouWantLogout":            "Chiqishni xohlaysizmi.",
			"promoCode":                  "Promo kod",
			"enterCode":                  "Kodni kiriting",
			"promoCodeInactive":          "Promo kod faol emas",
			"promoCodeActive":            "Promo kod faol",
			"promoCodeAlreadyActivated":  "Promo kod allaqachon faollashtirilgan",
			"promoCodePermissionDenied":  "Promo kodga ruxsat berilmagan",
			"promoCodeNotFound":          "Promo kod topilmadi",
			"status":                     "Status",
			"connectWithOperator":        "Operator bilan bog'lanish",
			"operatorMessage":            "<a href='https://t.me/turonsupport'>Operator</a> bilan bog'laning va savollaringizga javob toping.",
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
