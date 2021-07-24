package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token    = ""
	response = map[string]func(Request) string{
		"!도움": printHelp,
		"!추가": addUsers,
		"!제거": removeUsers,
		"!총원": printUsers,
		"!리셋": resetUsers,
		"!섞어": shuffleUsers,
		"!팀":  printTeam,
		"!전송": sendHome,
		"!상태": setState,
	}
	users    []string
	gdFlag   time.Time
	commands string
	team1    string
	team2    string
)

func setState(r Request) string {
	err := r.session.UpdateGameStatus(0, r.arg)
	if err != nil {
		log.Println("game status error", err)
	}
	return "완료!"
}

func init() {
	Token = os.Getenv("TOKEN")
	gdFlag = time.Now()
	commands = func() string {
		var commands []string
		for k := range response {
			commands = append(commands, k)
		}

		return "`" + strings.Join(commands, "`, `") + "`"
	}()
}

func printTeam(_ Request) string {
	return fmt.Sprintf("`1팀: %s`\n`2팀: %s`", team1, team2)
}

func shuffleUsers(r Request) string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})
	team1 = strings.Join(users[:len(users)/2], ", ")
	team2 = strings.Join(users[len(users)/2:], ", ")
	return sendHome(r)
}

func printHelp(_ Request) string {
	return commands
}

func resetUsers(_ Request) string {
	users = []string{}
	return "완료"
}

type Request struct {
	command string
	arg     string
	session *discordgo.Session
}

func NewRequest(s string, session *discordgo.Session) (Request, error) {
	if !strings.HasPrefix(s, "!") {
		return Request{}, errors.New("not command")
	}
	split := strings.SplitN(s, " ", 2)

	log.Println("receive command: ", s)

	if len(split) != 2 {
		split = append(split, "")
	}

	return Request{split[0], split[1], session}, nil

}

func printUsers(_ Request) string {
	if len(users) == 0 {
		return "empty!"
	}
	s := strings.Join(users, ", ")
	return fmt.Sprintf("```Total: %d\nUsers: %s```", len(users), s)
}

func addUsers(r Request) string {
	split := strings.Split(r.arg, ",")
	for _, v := range split {
		strings.Trim(v, " ")
		if !contains(v) {
			users = append(users, v)
		}
	}
	return printUsers(r)
}

func removeUsers(r Request) string {
	removeCount := 0
	split := strings.Split(r.arg, ",")
	for i, v := range split {
		strings.Trim(v, " ")
		if contains(v) {
			users = append(users[:i-removeCount], users[i+1-removeCount:]...)
			removeCount++
		}
	}
	return printUsers(r)
}

func sendHome(r Request) string {
	r.session.ChannelMessageSend("592275489100398594", printTeam(r))
	return "완료!"
}
func contains(s string) bool {
	for _, user := range users {
		if user == s {
			return true
		}
	}
	return false
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "hello world")
		if err != nil {
			panic("fail response")
		}
	})
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic(fmt.Errorf("http server error"))
		}
	}()

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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ㅎㅇ" {
		fmt.Println(time.Since(gdFlag).String())
		if time.Since(gdFlag) > 30*time.Minute {
			gdFlag = time.Now()
			s.ChannelMessageSend(m.ChannelID, "ㅎㅇ")
		}
	}

	request, err := NewRequest(m.Content, s)
	if err != nil {
		return
	}

	if f, ok := response[request.command]; ok {
		s.ChannelMessageSend(m.ChannelID, f(request))
	}
}
