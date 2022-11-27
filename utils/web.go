package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseBody(req *http.Request, body interface{}) error {
	return json.NewDecoder(req.Body).Decode(body)
}

func ParseNow() string {
	return fmt.Sprintf("%04d-%02d", time.Now().Year(), int(time.Now().Month()))
}

func GetDateInlineBtns(date string) *tgbotapi.InlineKeyboardMarkup {
	year, _ := strconv.Atoi(date[:4])
	month, _ := strconv.Atoi(date[5:])

	nextMonth := month
	nextYear := year + 1

	prevMonth := month
	prevYear := year - 1

	currentMonth := month + 1
	currentYear := year - 1
	if currentMonth > 12 {
		currentYear++
	}

	var btns tgbotapi.InlineKeyboardMarkup
	var matr [][]tgbotapi.InlineKeyboardButton
	for row := 0; row < 3; row++ {
		var arr []tgbotapi.InlineKeyboardButton
		for col := 0; col < 4; col++ {
			arr = append(arr, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%04d-%02d", currentYear, currentMonth), fmt.Sprintf("salary-info:%04d-%02d", currentYear, currentMonth)))
			currentMonth++
			if currentMonth > 12 {
				currentMonth = 1
				currentYear++
			}
		}

		matr = append(matr, arr)

	}

	if nextYear > time.Now().Year() {
		matr = append(matr, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("prev-salary:%04d-%02d", prevYear, prevMonth)),
		))
	} else {

		matr = append(matr, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("prev-salary:%04d-%02d", prevYear, prevMonth)),
			tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("next-salary:%04d-%02d", nextYear, nextMonth)),
		))
	}
	btns.InlineKeyboard = matr

	return &btns
}
