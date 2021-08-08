package model

type Users struct {
	ChatId int64
}

func (p *Users) TableName() string {
	return "users"
}
