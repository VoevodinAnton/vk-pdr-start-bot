package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	longpoll "github.com/SevereCloud/vksdk/v2/longpoll-bot"
	cron "github.com/robfig/cron/v3"
)

var VKToken = "890b57a42c1d670f6969f1336cf931c5d1dd6e3ad46538d0be86761544976acfc7a3679bf50cf483be251"

func main() {
	vk := api.NewVK(VKToken)

	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	lp, err := longpoll.NewLongPoll(vk, 212384138)
	if err != nil {
		panic(err)
	}

	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		println(obj.Message.Text)
		PeerId := obj.Message.PeerID
		//println(PeerId)
		if strings.Split(obj.Message.Text, " ")[1] == "start" {
			log.Print("Job is started")
			defer scheduler.Stop()
			scheduler.AddFunc("0 */12 * * *", func() { startPidorBot(vk, PeerId) })
			go scheduler.Start()
		}
	})
	lp.Run()
}

func startPidorBot(vk *api.VK, PeerId int) {
	_, err := vk.MessagesSend(api.Params{"peer_id": PeerId, "random_id": 0, "message": "@botv_pidor старт"})
	if err != nil {
		log.Fatal("Error sending message:", err)
	}
	log.Print("Bot is started")
}
