package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token    = ""
	response = map[string]func(string) string{
		"!추가": addUsers,
	}
)

func init() {
	Token = os.Getenv("TOKEN")
}

func addUsers(str string) string {
	return str
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func containsKey(key string) bool {
	_, ok := response[key]
	return ok
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ㅎㅇ" {
		s.ChannelMessageSend(m.ChannelID, "ㅎㅇ")
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	split := strings.SplitN(m.Content, " ", 2)

	command := split[0]
	content := split[1]

	if containsKey(command) {
		f := response[command]
		s.ChannelMessageSend(m.ChannelID, f(content))
	}
}
