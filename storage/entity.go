package storage

type User struct {
	ID          int64  `gorm:"id;primaryKey"`
	ChatID      int64  `gorm:"chat_id"`
	PhoneNumber string `gorm:"phone_number"`
	Language    string `gorm:"language"`
	Tin         string `gorm:"tin"`
	IsVerified  bool   `gorm:"is_verified"`
	State       int    `gorm:"state"`
}
