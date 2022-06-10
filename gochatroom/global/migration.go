package global

import "github.com/josexy/gochatroom/model"

func migration() {
	_ = DB.AutoMigrate(
		&model.User{},
		&model.Message{},
	)
}
