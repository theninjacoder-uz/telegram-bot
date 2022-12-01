package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tgbot/constants"
	"tgbot/models"
	"tgbot/storage"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// 2 get language and updated language and sended phone number
// 1 send language options
// 4 stir and phone saved to db
// 5 send salary data options and set state 5
// 6 send salary info
type HandlerService struct {
	storage storage.IUserStorage
	bot     *tgbotapi.BotAPI
}

func New(storage storage.IUserStorage, bot *tgbotapi.BotAPI) *HandlerService {
	return &HandlerService{
		storage: storage,
		bot:     bot,
	}
}

// This handler is called everytime telegram sends us a webhook event
func (h *HandlerService) GlobalHandler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	var update *tgbotapi.Update
	fmt.Println("update: ", &update)
	if err := utils.ParseBody(req, &update); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	// Check messafe
	if update.Message != nil {
		chatID := update.Message.Chat.ID

		// Get user
		user, err := h.storage.CheckUser(chatID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.storage.Create(&storage.User{
				ChatID:   chatID,
				Language: constants.RUS,
			})
			user = &storage.User{
				ChatID:   chatID,
				Language: constants.RUS,
				State:    0, // After start button
			}
		} else if err != nil {
			h.sendErrorMsgWithLang(chatID, constants.INTERNAL_ERROR_UZB, constants.INTERNAL_ERROR_ENG, constants.INTERNAL_ERROR_RUS, user.Language)
		}

		// Receive PHONE or STIR
		if user.State == 2 {

			if update.Message.Contact == nil {
				h.sendErrorMsgWithLang(chatID, constants.TEXT_PHONE_UZB, constants.TEXT_PHONE_ENG, constants.TEXT_PHONE_RUS, user.Language)
				return
			}

			phone := update.Message.Contact.PhoneNumber

			if p, ok := models.IsPhone(phone); ok {
				phone = p
				h.storage.SavePhone(chatID, p)
			} else if models.IsTin(phone) {
				h.storage.SaveTin(chatID, phone)
			} else {
				h.sendErrorMsgWithLang(chatID, constants.INCORRECT_PHONE_UZB, constants.INCORRECT_PHONE_ENG, constants.INCORRECT_PHONE_RUS, user.Language)
				return
			}

			ok, err := h.storage.GetUser(phone, phone)
			h.bot.Send(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))

			if !ok {
				if err != nil {
					h.sendErrorMsgWithLang(chatID, constants.INTERNAL_ERROR_UZB, constants.INTERNAL_ERROR_ENG, constants.INTERNAL_ERROR_RUS, user.Language)
					return
				}
				// h.storage.UpdateState(chatID, 2)
				h.sendErrorMsgWithLang(chatID, constants.USER_NOT_FOUND_UZB, constants.USER_NOT_FOUND_ENG, constants.USER_NOT_FOUND_RUS, user.Language)
				return
			}
			h.storage.UpdateState(chatID, 3)
			h.storage.UpdateVerify(chatID, true)
			h.sendMainMenu(chatID, user.Language)
			return
		}

		// msg := strings.ToLower(update.Message.Text)
		msg := update.Message.Text
		switch msg {
		case "/start":

			if user.State == 0 {
				h.sendLanguageOptions(*update, user.Language, "Здравствуйте! ")
				h.storage.UpdateState(chatID, 1) // send lang options
				return
			} else if user.IsVerified {
				h.sendMainMenu(user.ChatID, user.Language)
			} else {
				h.sendErrorMsgWithLang(chatID, constants.SIGN_UP_UZB, constants.SIGN_UP_ENG, constants.SIGN_UP_RUS, user.Language)
				h.replySendPhone(chatID, user.Language)
				h.storage.UpdateState(chatID, 2)
			}
		case "◀️":
			h.bot.Send(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))
			h.storage.UpdateState(chatID, 3)
			h.sendMainMenu(chatID, user.Language)

		case constants.SETTINGS_ENG:
			h.sendSettingsMenu(chatID, constants.ENG)
			h.storage.UpdateState(chatID, 5)
		case constants.SETTINGS_UZB:
			h.sendSettingsMenu(chatID, constants.UZB)
			h.storage.UpdateState(chatID, 5)
		case constants.SETTINGS_RUS:
			h.sendSettingsMenu(chatID, constants.RUS)
			h.storage.UpdateState(chatID, 5)

		case constants.SALARY_INFO_ENG:
			h.sendDateOptions(chatID, utils.ParseNow(), constants.ENG)
			h.storage.UpdateState(chatID, 4)
		case constants.SALARY_INFO_UZB:
			h.sendDateOptions(chatID, utils.ParseNow(), constants.UZB)
			h.storage.UpdateState(chatID, 4)
		case constants.SALARY_INFO_RUS:
			h.sendDateOptions(chatID, utils.ParseNow(), constants.RUS)
			h.storage.UpdateState(chatID, 4)

		case constants.CHANGE_LANG_ENG:
			h.sendLanguageOptions(*update, constants.ENG, "")
			h.storage.UpdateState(chatID, 6)
		case constants.CHANGE_LANG_UZB:
			h.sendLanguageOptions(*update, constants.UZB, "")
			h.storage.UpdateState(chatID, 6)
		case constants.CHANGE_LANG_RUS:
			h.sendLanguageOptions(*update, constants.RUS, "")
			h.storage.UpdateState(chatID, 6)

		case constants.CHANGE_PHONE_ENG:
			h.replySendPhone(chatID, constants.ENG)
			h.storage.UpdateState(chatID, 2)
		case constants.CHANGE_PHONE_RUS:
			h.replySendPhone(chatID, constants.RUS)
			h.storage.UpdateState(chatID, 2)
		case constants.CHANGE_PHONE_UZB:
			h.replySendPhone(chatID, constants.UZB)
			h.storage.UpdateState(chatID, 2)
		default:
			if user.State == 4 {
				date := update.Message.Text
				validDate, err := models.IsValidDate(date)
				if err != nil {
					h.sendErrorMsgWithLang(chatID, constants.INCORRECT_DATE_UZB, constants.INCORRECT_DATE_ENG, constants.INCORRECT_DATE_RUS, user.Language)
					return
				}
				salary, err := h.storage.GetSalaryInfo(validDate, user.PhoneNumber, user.Tin)
				if err != nil {
					h.sendErrorMsgWithLang(chatID, constants.INTERNAL_ERROR_UZB, constants.INTERNAL_ERROR_ENG, constants.INTERNAL_ERROR_RUS, user.Language)
					return
				}
				h.bot.Send(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))

				if user.Language == constants.UZB {
					h.bot.Send(tgbotapi.NewMessage(chatID, salary.Uzb))
				} else if user.Language == constants.ENG {
					h.bot.Send(tgbotapi.NewMessage(chatID, salary.Eng))
				} else {
					h.bot.Send(tgbotapi.NewMessage(chatID, salary.Rus))
				}

			}

		}

	} else if update.CallbackQuery != nil {
		chatID := update.CallbackQuery.Message.Chat.ID
		data := update.CallbackQuery.Data

		user, err := h.storage.CheckUser(chatID)
		if err != nil {
			h.sendErrorMsgWithLang(chatID, constants.INTERNAL_ERROR_UZB, constants.INTERNAL_ERROR_ENG, constants.INTERNAL_ERROR_RUS, user.Language)

			return
		}

		if user.State == 1 {
			err := models.ValidateLang(data)
			if err != nil {
				h.sendErrorMsgWithLang(chatID, err.Error(), err.Error(), err.Error(), user.Language)
				return
			}
			h.bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))
			h.storage.UpdateUserLang(chatID, data, 2)

			h.replySendPhone(chatID, data)

		}

		if user.State == 4 && strings.HasPrefix(data, "salary-info:") {
			fmt.Printf("salary info: %v", data)
			salary, err := h.storage.GetSalaryInfo(data[12:], user.PhoneNumber, user.Tin)
			if err != nil {
				h.sendErrorMsgWithLang(chatID, constants.INTERNAL_ERROR_UZB, constants.INTERNAL_ERROR_ENG, constants.INTERNAL_ERROR_RUS, user.Language)
				return
			}
			h.bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))

			if user.Language == constants.UZB {
				h.bot.Send(tgbotapi.NewMessage(chatID, salary.Uzb))
			} else if user.Language == constants.ENG {
				h.bot.Send(tgbotapi.NewMessage(chatID, salary.Eng))
			} else {
				h.bot.Send(tgbotapi.NewMessage(chatID, salary.Rus))
			}

			// h.storage.UpdateState(chatID, 6)
		}
		if user.State == 4 && (strings.HasPrefix(data, "prev-salary:") || strings.HasPrefix(data, "next-salary:")) {

			h.bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))
			h.sendDateOptions(chatID, data[12:], user.Language)
		}
		if user.State == 6 {
			h.bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))
			h.storage.UpdateUserLang(chatID, data, 3)
			h.sendMainMenu(chatID, data)
		}
	}
}

// Send language options
func (h *HandlerService) sendLanguageOptions(update tgbotapi.Update, lang, greeting string) error {
	newMsg := constants.SELECT_LANG_RUS

	if lang == constants.ENG {
		newMsg = constants.SELECT_LANG_ENG
	} else if lang == constants.UZB {
		newMsg = constants.SELECT_LANG_UZB
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s%s", greeting, newMsg))

	var langBtn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(constants.LANGUAGE_UZB, constants.UZB),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(constants.LANGUAGE_ENG, constants.ENG),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(constants.LANGUAGE_RUS, constants.RUS),
		),
	)
	msg.ReplyMarkup = langBtn

	_, err := h.bot.Send(msg)
	return err
}

func (h *HandlerService) sendErrorMsgWithLang(chatID int64, msgUz, msgEn, msgRu, lang string) {
	if lang == constants.UZB {
		h.bot.Send(tgbotapi.NewMessage(chatID, msgUz))
	} else if lang == constants.ENG {
		h.bot.Send(tgbotapi.NewMessage(chatID, msgEn))
	} else {
		h.bot.Send(tgbotapi.NewMessage(chatID, msgRu))
	}
}

func (h *HandlerService) replySendPhone(chatID int64, lang string) {

	textMain := constants.SEND_PHON_NUMBER_RUS
	textBtn := constants.PHONE_RUS
	if lang == constants.UZB {
		textMain = constants.SEND_PHON_NUMBER_UZB
		textBtn = constants.PHONE_UZB
	} else if lang == constants.ENG {
		textMain = constants.SEND_PHON_NUMBER_ENG
		textBtn = constants.PHONE_ENG
	}

	msg := tgbotapi.NewMessage(chatID, textMain)
	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(textBtn),
		),
	)
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *HandlerService) sendMainMenu(chatID int64, lang string) error {

	msg := tgbotapi.NewMessage(chatID, constants.MENU_RUS)
	salaryInfo := constants.SALARY_INFO_RUS
	settings := constants.SETTINGS_RUS
	if lang == constants.UZB {
		msg = tgbotapi.NewMessage(chatID, constants.MENU_UZB)
		salaryInfo = constants.SALARY_INFO_UZB
		settings = constants.SETTINGS_UZB

	} else if lang == constants.ENG {
		msg = tgbotapi.NewMessage(chatID, constants.MENU_ENG)
		salaryInfo = constants.SALARY_INFO_ENG
		settings = constants.SETTINGS_ENG
	}
	var btns = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(salaryInfo),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(settings),
		),
	)
	msg.ReplyMarkup = btns

	_, err := h.bot.Send(msg)
	return err
}

func (h *HandlerService) sendDateOptions(chatID int64, date, lang string) error {
	text := constants.PERIOD_RUS
	if lang == constants.UZB {
		text = constants.PERIOD_UZB

	} else if lang == constants.ENG {
		text = constants.PERIOD_ENG
	}

	fmt.Printf("keldi: %v", text)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = utils.GetDateInlineBtns(date)

	_, err := h.bot.Send(msg)
	return err
}

func (h *HandlerService) sendSettingsMenu(chatID int64, lang string) error {

	msg := tgbotapi.NewMessage(chatID, constants.SETTINGS_MENU_RUS)
	langSettings := constants.CHANGE_LANG_RUS
	phone := constants.CHANGE_PHONE_RUS
	if lang == constants.UZB {
		msg = tgbotapi.NewMessage(chatID, constants.SETTINGS_MENU_UZB)
		langSettings = constants.CHANGE_LANG_UZB
		phone = constants.CHANGE_PHONE_UZB

	} else if lang == constants.ENG {
		msg = tgbotapi.NewMessage(chatID, constants.SETTINGS_MENU_ENG)
		langSettings = constants.CHANGE_LANG_ENG
		phone = constants.CHANGE_PHONE_ENG
	}
	var btns = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(phone),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(langSettings),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("◀️"),
		),
	)
	msg.ReplyMarkup = btns

	_, err := h.bot.Send(msg)
	return err
}
