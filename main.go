package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yanzay/tbot/v2"
)

const checkPeriod = 1 * time.Minute

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	cli := bot.Client()
	store := &Store{subscriptions: make(map[string]bool)}
	bot.HandleMessage("/start", func(m *tbot.Message) {
		store.Subscribe(m.Chat.ID)
		cli.SendMessage(m.Chat.ID, "subscribed")
	})
	bot.HandleMessage("/stop", func(m *tbot.Message) {
		store.Unsubscribe(m.Chat.ID)
		cli.SendMessage(m.Chat.ID, "unsubscribed")
	})
	ctx := context.Background()
	go watch(ctx, cli, store)
	log.Fatal(bot.Start())
}

func watch(ctx context.Context, cli *tbot.Client, store *Store) {
	ticker := time.NewTicker(checkPeriod)
	for {
		select {
		case <-ticker.C:
			msg, ok := checkForWindows()
			if ok && !store.IsNotified() {
				for _, sub := range store.List() {
					cli.SendMessage(sub, msg)
				}
				store.SetNotified(true)
			}
			if !ok {
				store.SetNotified(false)
			}
		case <-ctx.Done():
			return
		}
	}
}

func checkForWindows() (string, bool) {
	days, err := GetWindows("48215633", "kiev_desnianskyi")
	if err != nil {
		log.Println(err)
		return "", false
	}
	if len(days) == 0 {
		fmt.Print(".")
		return "", false
	}
	str := &strings.Builder{}
	str.WriteString("Есть слоты доставки:\n")
	for _, day := range days {
		fmt.Fprintf(str, "%s (%s)\n", day.Title, day.Date)
		for _, window := range day.Windows {
			fmt.Fprintf(str, "%s %d₴\n", window.Title, window.Price.Num0/100)
		}
	}
	return str.String(), true
}
