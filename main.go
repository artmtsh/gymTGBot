package main

import "github.com/artmtsh/gymTGBot/gym"

func main() {
	app := gym.NewApp()

	// здесь настраиваешь Telegram-библиотеку
	// и в цикле:
	//  - получаешь userID и text
	//  - вызываешь reply := app.HandleText(userID, text)
	//  - отправляешь reply обратно пользователю
}
