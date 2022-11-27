package models

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

const (
	MIN_YEAR = 1980
)

type CheckModel struct {
	Ok   bool `json:"ok"`
	Down bool `json:"down"`
}

type SalaryInfo struct {
	Eng string `json:"en,omitempty"`
	Rus string `json:"ru,omitempty"`
	Uzb string `json:"uz,omitempty"`
}

func ValidateLang(lang string) error {
	if lang == "UZB" || lang == "RUS" || lang == "ENG" {
		return nil
	}
	return fmt.Errorf("%v error: ", lang)
}

var uzbPhoneRegex = regexp.MustCompile(`[+]{0,1}99{1}[0-9]{10}$`)

// IsPhoneValid validates phone number for Uzbekistan
func IsPhone(p string) bool {
	return uzbPhoneRegex.MatchString(p)
}

func IsTin(v string) bool {
	for _, c := range v {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return len(v) >= 9 && len(v) <= 11
}

func IsValidDate(date string) (string, error) {

	if len(date) < 6 || len(date) > 7 {
		return "", fmt.Errorf("error")
	}

	year, err := strconv.Atoi(date[:4])
	if err != nil || year > time.Now().Year() || year < MIN_YEAR {
		return "", fmt.Errorf("error")
	}
	month, err := strconv.Atoi(date[5:])
	fmt.Printf("month %d", month)
	if err != nil || month > 12 || month < 1 {
		return "", fmt.Errorf("error")
	}

	return fmt.Sprintf("%04d-%02d", year, month), nil
}
