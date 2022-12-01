package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tgbot/constants"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseBody(req *http.Request, body interface{}) error {
	return json.NewDecoder(req.Body).Decode(body)
}

func ParseNow() string {
	return fmt.Sprintf("%04d-%02d", time.Now().Year(), int(time.Now().Month()))
}

func GetDateInlineBtns(date, lang string) *tgbotapi.InlineKeyboardMarkup {

	nextBtnName := constants.NEXT_DATE_ENG
	prevBtnName := constants.PREVIOUS_DATE_ENG

	if lang == constants.UZB {
		nextBtnName = constants.NEXT_DATE_UZB
		prevBtnName = constants.PREVIOUS_DATE_UZB
	} else if lang == constants.RUS {
		nextBtnName = constants.NEXT_DATE_RUS
		prevBtnName = constants.PREVIOUS_DATE_RUS
	}

	year, _ := strconv.Atoi(date[:4])
	month, _ := strconv.Atoi(date[5:])

	nextMonth := month
	nextYear := year + 1

	prevMonth := month
	prevYear := year - 1

	currentMonth := month + 1
	currentYear := year - 1
	if currentMonth > 12 {
		currentMonth = 1
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
			tgbotapi.NewInlineKeyboardButtonData(prevBtnName, fmt.Sprintf("prev-salary:%04d-%02d", prevYear, prevMonth)),
		))
	} else {

		matr = append(matr, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(prevBtnName, fmt.Sprintf("prev-salary:%04d-%02d", prevYear, prevMonth)),
			tgbotapi.NewInlineKeyboardButtonData(nextBtnName, fmt.Sprintf("next-salary:%04d-%02d", nextYear, nextMonth)),
		))
	}
	btns.InlineKeyboard = matr

	return &btns
}
