package service

import (
	"collab-node-platform-backend/config"
	"collab-node-platform-backend/db"
	"collab-node-platform-backend/model"
	"collab-node-platform-backend/utils"
	"errors"
)

type LoginResult struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

func Login(username, password string) (*LoginResult, error) {
	var user model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("用户名不存在")
	}
	if !utils.CheckPassword(user.Password, password) {
		return nil, errors.New("密码错误")
	}
	token, err := utils.GenerateJWT(user.UserID, config.AppConfig.JwtSecret, config.AppConfig.JwtExpire)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}
	return &LoginResult{
		Token:  token,
		UserID: user.UserID,
	}, nil
}

func Register(username, password string) (string, error) {
	var count int64
	db.DB.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return "", errors.New("用户名已存在")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", errors.New("密码加密失败")
	}
	user := model.User{
		UserID:   utils.GenerateUUID(),
		Username: username,
		Password: hash,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return "", errors.New("注册失败")
	}
	return user.UserID, nil
}
