package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
)

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Username       string `gorm:"index;column:username;size:20;not null"`
	Nickname       string `gorm:"size:20;not null"`
	AvatarId       int    `gorm:"avatar_id"`
	PasswordDigest string `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// EncryptPassword 加密密码
func (user *User) EncryptPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 验证密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
