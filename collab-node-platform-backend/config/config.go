package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MysqlDSN  string
	JwtSecret string
	JwtExpire int64
	WsPath    string
}

var AppConfig *Config

func InitConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("未找到 .env 文件，使用环境变量或默认值")
	}
	AppConfig = &Config{
		Port:      getEnv("PORT", "3000"),
		MysqlDSN:  getEnv("MYSQL_DSN", "root:password@tcp(localhost:3306)/collab_node_platform?charset=utf8mb4&parseTime=True&loc=Local"),
		JwtSecret: getEnv("JWT_SECRET", "your_jwt_secret"),
		WsPath:    getEnv("WS_PATH", "/ws"),
	}
	jwtExpireStr := getEnv("JWT_EXPIRE", "86400")
	jwtExpire, err := strconv.ParseInt(jwtExpireStr, 10, 64)
	if err != nil {
		jwtExpire = 86400
	}
	AppConfig.JwtExpire = jwtExpire
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
