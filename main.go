package main

import (
	"bytes"
	"jacob/dusk/pkg/run"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	log "github.com/Sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var token = kingpin.Flag("token", "telegram bot token").Envar("TBTOKEN").Required().String()

func main() {
	bot, err := tb.NewBot(tb.Settings{
		Token: *token,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})

	if err != nil {
		log.Fatalf("failed to created bot: %v", err)
	}

	log.Infoln("Created bot")

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Infoln(m.Text)
		offset := len("@DuskCBot ")
		if len(m.Text) > offset {
			if strings.TrimSpace(m.Text[:offset]) == "@DuskCBot" {
				done := make(chan struct{})
				stop := make(chan struct{})

				reader := strings.NewReader(m.Text[offset:])
				var writer bytes.Buffer

				go func() {
					run.Run(reader, &writer, m.Sender.Username, stop)
					close(done)
				}()

				select {
				case <-done:
					bot.Send(m.Chat, writer.String())
					return
				case <-time.NewTimer(5 * time.Second).C:
					log.Infoln("stopping eval")
					close(stop)
					bot.Send(m.Chat, "timeout: more than 5 seconds")
					return
				}
			}
		}
	})

	bot.Start()
}
