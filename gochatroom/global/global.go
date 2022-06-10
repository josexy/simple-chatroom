package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	AppConfig = new(Config)
	DB        *gorm.DB
	Redis     redis.UniversalClient
	Logger    *logrus.Logger
)
