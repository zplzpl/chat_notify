package repository

import (
	"context"

	"chat_notify/repository/model"
	"github.com/jinzhu/gorm"
)

func (p *Repository) GetChatCnt(ctx context.Context) (int64, error) {

	var cnt struct{ Cnt int64 }

	err := p.db.Select("COUNT(*) AS cnt").Model(model.Users{}).Scan(&cnt).Error

	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return cnt.Cnt, nil
}

func (p *Repository) GetChatList(ctx context.Context, offset, limit int64) ([]*model.Users, error) {

	var list []*model.Users
	err := p.db.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
