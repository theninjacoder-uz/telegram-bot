package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tgbot/configs"
	"tgbot/models"
	"time"

	"gorm.io/gorm"
)

type IUserStorage interface {
	// CreateTable(ctx context.Context, req *entity.Table) error
	// DeleteTable(ctx context.Context, req *entity.Table) error
	CheckUser(chatID int64) (*User, error)
	UpdateUserLang(chatID int64, lang string, state int) error
	SavePhone(chatID int64, phone string) error
	SaveTin(chatID int64, tin string) error
	UpdateState(chatID int64, state int) error
	Create(user *User) error
	GetUser(tin, phone string) (bool, error)
	GetSalaryInfo(date, phone, tin string) (*models.SalaryInfo, error)
	UpdateVerify(chatID int64, status bool) error
}

type userRepo struct {
	db *gorm.DB
}

// New ...
func New(db *gorm.DB) IUserStorage {
	return &userRepo{db: db}
}

func (r *userRepo) CheckUser(chatID int64) (*User, error) {
	var user *User
	tx := r.db.Take(&user, "chat_id=?", chatID)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (r *userRepo) UpdateUserLang(chatID int64, lang string, state int) error {
	tx := r.db.Model(&User{}).Where("chat_id=?", chatID).Updates(User{Language: lang, State: state})
	if tx.Error != nil || tx.RowsAffected != 1 {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *userRepo) SavePhone(chatID int64, phone string) error {
	tx := r.db.Model(&User{}).Where("chat_id=?", chatID).Updates(User{PhoneNumber: phone})
	if tx.Error != nil || tx.RowsAffected != 1 {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *userRepo) SaveTin(chatID int64, tin string) error {
	tx := r.db.Model(&User{}).Where("chat_id=?", chatID).Updates(User{Tin: tin})
	if tx.Error != nil || tx.RowsAffected != 1 {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *userRepo) UpdateState(chatID int64, state int) error {
	tx := r.db.Model(&User{}).Where("chat_id=?", chatID).Update("state", state)
	if tx.Error != nil || tx.RowsAffected != 1 {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *userRepo) UpdateVerify(chatID int64, status bool) error {
	tx := r.db.Model(&User{}).Where("chat_id=?", chatID).Update("is_verified", status)
	if tx.Error != nil || tx.RowsAffected != 1 {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *userRepo) Create(user *User) error {
	tx := r.db.Create(user)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *userRepo) GetUser(tin, phone string) (bool, error) {
	reqUrl := configs.Config().RpcHost + "/check_user/api/v1/users?phone=" + phone + "&tin=" + tin
	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.SetBasicAuth(configs.Config().RpcAuthLogin, configs.Config().RpcAuthPassword)

	fmt.Printf("Zafar's api req: %v\n", req)

	client := &http.Client{
		Timeout: 2000 * time.Millisecond,
	}

	result, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("serverda xatolik")
	}

	var ok models.CheckModel
	data, _ := io.ReadAll(result.Body)
	_ = json.Unmarshal([]byte(data), &ok)
	fmt.Printf("daTA %v\n", ok)

	return ok.Ok, nil
}

func (r *userRepo) GetSalaryInfo(date, phone, tin string) (*models.SalaryInfo, error) {
	if tin == "" {
		tin = phone
	} else if phone == "" {
		phone = tin
	}
	reqUrl := fmt.Sprintf("%s/get_report?tin=%s&phone=%s&period=%s", configs.Config().RpcHost, tin, phone, date)
	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.SetBasicAuth(configs.Config().RpcAuthLogin, configs.Config().RpcAuthPassword)

	fmt.Printf("Zafar's api req: %v\n", req)

	client := &http.Client{
		Timeout: 5000 * time.Millisecond,
	}

	result, err := client.Do(req)
	fmt.Printf("result %v\n", result.Body)
	fmt.Printf("error %v\n", err)
	var salary models.SalaryInfo
	if err != nil {
		return nil, fmt.Errorf("serverda xatolik")
	}
	data, _ := io.ReadAll(result.Body)
	// fmt.Printf("data %v\n", string(data[:]))

	json.Unmarshal([]byte(data), &salary)
	// fmt.Printf("daTA %v\n", &salary)

	return &salary, nil
}
