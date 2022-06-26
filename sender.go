package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/panjf2000/ants"
	"go.uber.org/atomic"
)

type Sender struct {
	workerNum    int
	telegramBot  *tgbotapi.BotAPI
	pool         *ants.PoolWithFunc
	successTotal *atomic.Int64
}

type SendEntry struct {
	ChatId  int64
	Message string
}

func NewSender(workerNum int, telegramBot *tgbotapi.BotAPI) (*Sender, error) {

	s := &Sender{
		workerNum:    workerNum,
		telegramBot:  telegramBot,
		successTotal: atomic.NewInt64(0),
	}

	pool, err := ants.NewPoolWithFunc(workerNum, s.process)
	if err != nil {
		return nil, err
	}
	s.pool = pool

	return s, nil
}

func (s *Sender) handleSendEntry(e *SendEntry) error {

	msg := tgbotapi.NewMessage(e.ChatId, e.Message)

	var err error
	for i := 0; i < 10; i++ {

		_, err = s.telegramBot.Send(msg)
		if err != nil {
			log.Println("telegramBot.Send err: ", err.Error())
			time.Sleep(10 * time.Duration(i) * time.Millisecond)
			continue
		}

		break

	}

	return err
}

func (s *Sender) process(payload interface{}) {

	sendEntry, ok := payload.(*SendEntry)
	if !ok {
		return
	}

	if err := s.handleSendEntry(sendEntry); err != nil {
		log.Println("send entry handle err: ", err.Error(), "chat_id: ", sendEntry.ChatId)
		return
	}

	s.successTotal.Inc()

}

func (p *Sender) Invoke(payload interface{}) error {
	return p.pool.Invoke(payload)
}

func (p *Sender) Close() error {

	p.pool.Tune(0)
	p.pool.Release()

	return nil
}
