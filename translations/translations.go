package translations

import (
	"sync"

	"github.com/OzodbekX/TuronMiniApp/volumes"
)

var (
	translations = map[string]map[string]string{
		"ru": {
			"welcome":                     "Добро пожаловать! Пожалуйста, выберите язык.",
			"login":                       "Пожалуйста, введите ваш логин.",
			"cancel":                      "⬅️ Отмена",
			"mainMenu":                    "🏠 Главное меню",
			"enterPassword":               "Введите ваш пароль:",
			"wrongParol":                  "Неправильный логин или пароль, пожалуйста, введите логин еще раз",
			"Tariffs":                     "Тарифы",
			"FAQ":                         "FAQ",
			"Application":                 "Заявление",
			"Language":                    "Язык",
			"Balance":                     "Баланс",
			"Exit":                        "Выход",
			"GoBack":                      "Назад",
			"PleaseSelectOption":          "Пожалуйста, выберите раздел из меню:",
			"pleaseEnterYourName":         "Пожалуйста, введите свое имя:",
			"enterPhone":                  "Введите номер телефона в следующем формате: +998#########",
			"pleaseShareYourPhoneNumber":  "Пожалуйста, поделитесь номером телефона",
			"listOfTariffs":               "Список тарифов",
			"price":                       "Цена",
			"speedByTime":                 "Скорость по времени",
			"mbs":                         "Мбит/с",
			"uzs":                         "сум",
			"sharePhoneNumber":            "Поделиться номером телефона",
			"invalidPhoneNumber":          "Неверный формат номера телефона. Укажите действительный номер: +998######### или #########.",
			"loginSuccessful":             "Авторизация успешна! Добро пожаловать в главное меню.",
			"yourBalance":                 "Ваш баланс",
			"tariffName":                  "Ваш тариф",
			"subscriptionPrice":           "Цена подписки",
			"nextSubscriptionDate":        "Дата следующей подписки",
			"subscriptionPeriod":          "Период подписки",
			"from":                        "с",
			"to":                          "по",
			"subscriptionActive":          "Подписка активна",
			"active":                      "Активно",   // Active in Russian
			"inactive":                    "Неактивно", // Inactive in Russian
			"pleaseSelectCategory":        "Пожалуйста, выберите категорию:",
			"pleaseSelectFAQ":             "Пожалуйста, выберите FAQ",
			"yes":                         "✅ Да",
			"no":                          "❌ Нет",
			"doYouWantLogout":             "Вы хотите выйти из системы?",
			"promoCode":                   "Промо-код",
			"enterCode":                   "Введите код",
			"promoCodeInactive":           "Промокод неактивен",
			"promoCodeActive":             "Промокод активен",
			"promoCodeAlreadyActivated":   "Промокод уже активирован",
			"promoCodePermissionDenied":   "Доступ к промокоду запрещен",
			"promoCodeNotFound":           "Промокод не найден",
			"status":                      "Статус",
			"connectWithOperator":         "Связаться с оператором",
			"operatorMessage":             "<a href='https://t.me/turonsupport'>Оператор</a>, свяжитесь и получите ответы на свои вопросы.",
			"chooseRole":                  "Пожалуйста, выберите кто вы:",
			"abonent":                     "📱 Абонент",
			"user":                        "👤 Пользователь",
			"pleaseSelectYurDistrict":     "Пожалуйста, выберите ваш город:",
			"pleaseSelectYurRegion":       "Пожалуйста, выберите свой район:",
			"enterFullName":               "📝 Пожалуйста, введите своё имя:",
			"errorFetchingData":           "⚠️ Ошибка при получении данных. Попробуйте позже.",
			"enterAdditionalPhone":        "Введите дополнительный телефон",
			"sendApplication":             "📨 Отправить заявку",
			"placeEnterFullName":          "Пожалуйста, введите Ваше полное имя (не менее 3 букв).",
			"wrongPhoneNumber":            "❗Неверный формат номера телефона. Введите правильный номер.",
			"applicationSentSuccessfully": "✅ Ваша заявка успешно отправлена! Оператор свяжется с вами в ближайшее время.",
			"failedToSendApp":             "❌ Произошла ошибка при отправке заявки. Пожалуйста, попробуйте позже.",
		},
		"uz": {
			"welcome":                     "Xush kelibsiz! Iltimos, tilni tanlang.",
			"login":                       "Iltimos loginingizni kiriting",
			"cancel":                      "⬅️ Bekor qilish",
			"mainMenu":                    "🏠 Asosiy menyu",
			"enterPassword":               "Iltimos, parolingizni kiriting:",
			"wrongParol":                  "Noto'g'ri login yoki parol, iltimos loginni qayta kiriting",
			"Tariffs":                     "Tariflar",
			"FAQ":                         "FAQ",
			"Application":                 "Ariza",
			"Language":                    "Til",
			"Balance":                     "Balans",
			"Exit":                        "Chiqish",
			"GoBack":                      "Ortga",
			"PleaseSelectOption":          "Iltimos, menyudan bo‘limni tanlang:",
			"pleaseEnterYourName":         "Iltimos, ismingizni kiriting:",
			"enterPhone":                  "Iltimos, telefon raqamingizni quyidagi formatda kiriting: +998#########",
			"pleaseShareYourPhoneNumber":  "Iltimos, telefon raqamingizni ulashing",
			"listOfTariffs":               "Tariflar ro'yxati",
			"price":                       "Narxi",
			"speedByTime":                 "Vaqt bo'yicha tezlik",
			"mbs":                         "Mbit/s",
			"uzs":                         "so'm",
			"sharePhoneNumber":            "Telefon raqamini ulashing",
			"invalidPhoneNumber":          "Telefon raqami formati noto‘g‘ri. Iltimos, haqiqiy raqamni kiriting: +998######### yoki #########.",
			"loginSuccessful":             "Kirish muvaffaqiyatli! Bosh menyuga xush kelibsiz.",
			"yourBalance":                 "Sizning balansingiz",
			"tariffName":                  "Sizning tarifingiz",
			"subscriptionPrice":           "Obuna narxi",
			"nextSubscriptionDate":        "Keyingi obuna sanasi",
			"subscriptionPeriod":          "Obuna davri",
			"from":                        "dan",
			"to":                          "gacha",
			"subscriptionActive":          "Obuna faolligi",
			"active":                      "Faol",    // Active in Uzbek
			"inactive":                    "No Faol", // Inactive in Uzbek
			"pleaseSelectCategory":        "Kategoriya tanlang:",
			"pleaseSelectFAQ":             "Iltimos FAQni tanlang",
			"yes":                         "✅ Ha",
			"no":                          "❌ Yo'q",
			"doYouWantLogout":             "Chiqishni xohlaysizmi.",
			"promoCode":                   "Promo kod",
			"enterCode":                   "Kodni kiriting",
			"promoCodeInactive":           "Promo kod faol emas",
			"promoCodeActive":             "Promo kod faol",
			"promoCodeAlreadyActivated":   "Promo kod allaqachon faollashtirilgan",
			"promoCodePermissionDenied":   "Promo kodga ruxsat berilmagan",
			"promoCodeNotFound":           "Promo kod topilmadi",
			"status":                      "Status",
			"connectWithOperator":         "Operator bilan bog'lanish",
			"operatorMessage":             "<a href='https://t.me/turonsupport'>Operator</a> bilan bog'laning va savollaringizga javob toping.",
			"chooseRole":                  "Iltimos, kimligingizni tanlang:",
			"abonent":                     "📱 Abonent",
			"user":                        "👤 Foydalanuvchi",
			"pleaseSelectYurDistrict":     "Iltimos shaharingizni tanlang:",
			"pleaseSelectYurRegion":       "Iltimos tumaningizni tanlang:",
			"enterFullName":               "📝 Iltimos, ismingizni kiriting:",
			"errorFetchingData":           "⚠️ Ma'lumotlarni olishda xatolik. Keyinroq urinib ko‘ring.",
			"enterAdditionalPhone":        "Qo'shimcha telefonni kiriting",
			"sendApplication":             "📨 Ariza yuborish",
			"placeEnterFullName":          "Iltimos, to‘liq ism kiriting (kamida 3 ta harf).",
			"wrongPhoneNumber":            "❗Telefon raqam formati noto‘g‘ri. Iltimos, to‘g‘ri raqam kiriting.",
			"applicationSentSuccessfully": "✅ Arizangiz muvaffaqiyatli yuborildi! Tez orada operator siz bilan bog‘lanadi.",
			"failedToSendApp":             "❌ Arizani yuborishda xatolik yuz berdi. Iltimos, keyinroq urinib ko‘ring.",
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
