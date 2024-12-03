package main

var (
	translations = map[string]map[string]string{
		"ru": {"welcome": "Добро пожаловать! Пожалуйста, выберите язык.", "login": "Пожалуйста, отправьте свои учетные данные в формате 'имя_пользователя:пароль'."},
		"uz": {"welcome": "Xush kelibsiz! Iltimos, tilni tanlang.", "login": "Iltimos, 'foydalanuvchi_nomi:parol' formatida kiring."},
	}
)

func getTranslation(lang, key string) string {
	if text, ok := translations[lang][key]; ok {
		return text
	}
	return key // Fallback to key if translation is missing
}
