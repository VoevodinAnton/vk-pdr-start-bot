package main

import (
	"bytes"
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
var cronExpression = "0 */12 * * *"
var entryID cron.EntryID

func main() {
	vk := api.NewVK(VKToken)

	jakartaTime, _ := time.LoadLocation("Europe/Samara")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	lp, err := longpoll.NewLongPoll(vk, 212384138)
	if err != nil {
		panic(err)
	}

	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		PeerId := obj.Message.PeerID
		command := getMessageObjectText(obj)

		if len(command) > 1 {
			runCommand(scheduler, vk, PeerId, command[1])
		}
	})
	lp.Run()
}

func getMessageObjectText(obj events.MessageNewObject) []string {
	return strings.Fields(obj.Message.Text)
}

//убрать returnы если получится
func runCommand(scheduler *cron.Cron, vk *api.VK, PeerId int, command string) {
	if command == "start" {
		if len(scheduler.Entries()) == 0 {
			entryID, _ = scheduler.AddFunc(cronExpression, func() { writeMessage(vk, PeerId, "@botv_pidor старт") })
		}
		scheduler.Start()
		log.Print("Job is started")
	} else if command == "stop" {
		scheduler.Stop()
		scheduler.Remove(entryID)
		log.Print("Job is stopped")
	} else if command == "next" {
		if len(scheduler.Entries()) > 0 {
			time := scheduler.Entries()[0].Next.Format("2006-01-02 15:04:05")
			writeMessage(vk, PeerId, "Следующий запуск будет в "+time)
		} else {
			writeMessage(vk, PeerId, "Pidor Bot Helper не запущен")
		}
	} else {
		writeMessage(vk, PeerId, unknownCommand())
	}
}

func writeMessage(vk *api.VK, PeerId int, message string) {
	_, err := vk.MessagesSend(api.Params{"peer_id": PeerId, "random_id": 0, "message": message})
	if err != nil {
		log.Fatal("Error sending message:", err)
	}
}

func unknownCommand() string {
	var buffer bytes.Buffer
	buffer.WriteString("Такой команды не существует\n\n")
	buffer.WriteString("Команды:\n")
	buffer.WriteString("@botv_pidor_helper start\n")
	buffer.WriteString("@botv_pidor_helper stop\n")
	buffer.WriteString("@botv_pidor_helper next")
	return buffer.String()
}
