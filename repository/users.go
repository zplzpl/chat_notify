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

func (p *Repository) GetChatList(ctx context.Context, scope string, offset, limit int64) ([]*model.Users, error) {

	var list []*model.Users

	q := p.db
	switch scope {
	case "sub":
		q.Where("is_subscribed = 1")
	case "unsub":
		q.Where("is_subscribed = 0")
	}

	err := q.Order("id desc").Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
