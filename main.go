package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"chat_notify/repository"
)

func main() {

	// flag params
	configFilename := flag.String("cfg", "", "config filename,didn't need input file extend name.")
	notifyFilename := flag.String("notify", "", "notify filename.")
	workNum := flag.Int("worker", 5, "send worker num. default: 5")
	flag.Parse()

	if len(*configFilename) == 0 {
		panic(errors.New("config filename is empty"))
	}

	if len(*notifyFilename) == 0 {
		panic(errors.New("notify filename is empty"))
	}

	log.Println("config file: ", *configFilename, "notify file: ", *notifyFilename)

	// read notify
	buf, err := ioutil.ReadFile(fmt.Sprintf("./notify/%s", *notifyFilename))
	if err != nil {
		panic(err)
	}

	if len(buf) == 0 {
		log.Println("notify file content is empty")
		return
	}

	log.Println("your send notify content: ")
	log.Println(string(buf))

	// config init
	upCfg := make(chan struct{})
	err = WatchConfig(upCfg, *configFilename)
	if err != nil {
		panic(err)
	}
	defer close(upCfg)

	// os signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer close(signals)

	// run loop
	if err := runLoop(upCfg, signals, string(buf), *workNum); err != nil {
		log.Println("run error: ", err.Error())
		panic(err)
	} else {
		os.Exit(0)
	}

}

func runLoop(upCfg chan struct{}, signals chan os.Signal, message string, workerNum int) error {

	// load config
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	_ = cfg

	// context
	ctx := context.Background()

	// init db
	if err := InitDB(cfg.MysqlDB); err != nil {
		return err
	}

	// init repo
	repo := repository.NewRepository(DB)

	// get chat cnt
	cnt, err := repo.GetChatCnt(ctx)
	if err != nil {
		return err
	}

	log.Println("chat count: ", cnt)

	if cnt == 0 {
		return nil
	}

	// init telegram bot
	if err := InitTelegramBot(cfg.Telegram.BotKey); err != nil {
		return err
	}

	// init sender
	sender, err := NewSender(workerNum, telegramBot)
	if err != nil {
		return err
	}

	if sender == nil {
		return nil
	}

	log.Println("init sender worker success, worker num: ", workerNum)

	// for ...loop chat
	var offset, limit int64 = 0, 100
	log.Println("start run sender...")
	var num int
	for {

		list, err := repo.GetChatList(ctx, offset, limit)
		if err != nil {
			continue
		}

		if len(list) == 0 {
			break
		}

		log.Println("send offset: ", offset, "limit: ", limit)

		num = num + len(list)

		offset = offset + limit

		for _, item := range list {
			_ = sender.Invoke(&SendEntry{
				ChatId:  item.ChatId,
				Message: message,
			})
		}

	}

	// close sender
	if err := sender.Close(); err != nil {
		log.Println("sender close err: ", err.Error())
	}

	successTotal := sender.successTotal.Load()

	log.Println("send success total: ", successTotal)
	log.Println("send failed total: ", int64(num)-successTotal)
	log.Println("send total: ", num)

	return nil
}
