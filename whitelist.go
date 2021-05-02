package main

func checkWhiteList(chatId string) bool {
	chatIds := map[string]bool{
		"1771439892":         true,
		"https://paypal.com": true,
	}
	if chatIds[chatId] {
		return true
	}
	return false
}
