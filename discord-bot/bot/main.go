package bot

import (
	"fmt"
	"log"
	"math/rand/v2"
	"vivalchemy/discord-bot/config"

	"github.com/bwmarrin/discordgo"
)

var (
	botId        string
	golangQuotes = []string{
		"The key to making programs fast is to make them do practically nothing. — Rob Pike",
		"Go is about software engineering, not just programming. — Rob Pike",
		"Clear is better than clever. — Rob Pike",
		"The biggest source of mistakes in C and C++ is using null pointers. — Ken Thompson",
		"Go is pragmatic, simple, and productive. — Andrew Gerrand",
		"The goal of Go is to make programming fun again. — Rob Pike",
		"Less is exponentially more. — Rob Pike",
		"Go is simple, but simple doesn’t mean easy. — William Kennedy",
		"Concurrency is not parallelism. — Rob Pike",
		"Go is what happens when you rethink systems programming from scratch. — Anonymous",
		"Don’t communicate by sharing memory, share memory by communicating. — Go Concurrency Principle",
		"Go’s simplicity is its power. — Francesc Campoy",
		"Channels orchestrate; mutexes serialize. — Rob Pike",
		"Goroutines are cheap, but they are not free. — Dave Cheney",
		"A little copying is better than a little dependency. — Go proverb",
		"Gofmt’s style is no one’s favorite, but gofmt is everyone’s favorite. — Rob Pike",
		"Go makes it easy to write correct programs, but hard to write incorrect ones. — Anonymous",
		"Errors are just values. — Dave Cheney",
		"Code is read more than it is written. — Rob Pike",
		"The only way to get a Go program to crash is to explicitly panic. — Dave Cheney",
		"A well-designed Go program feels like poetry. — Anonymous",
		"Go enforces good habits through its design, not just its rules. — Anonymous",
		"Go is for people who want to spend their time thinking about their problem domain, not their programming language. — Anonymous",
		"Make the zero value useful. — Go Design Philosophy",
		"If performance matters, measure. — Rob Pike",
		"Reflection in Go is powerful, but should be used sparingly. — Anonymous",
		"Simple code is easier to maintain than clever code. — Anonymous",
		"The best optimization is the one you don’t need. — Rob Pike",
		"In Go, interfaces are satisfied implicitly, not explicitly. — Anonymous",
		"Go is designed to scale from small scripts to large distributed systems. — Anonymous",
		"Idiomatic Go is simple Go. — Francesc Campoy",
		"Go has no generics, but you rarely need them. — Rob Pike (before generics were introduced)",
		"Error handling in Go is explicit and predictable. — Anonymous",
		"A Go program should look like it was written by a single person, no matter how many people worked on it. — Rob Pike",
		"The garbage collector lets you focus on solving problems, not managing memory. — Anonymous",
		"Go was designed for speed, simplicity, and safety. — Anonymous",
		"The worst Go code I’ve ever seen is still readable. — Anonymous",
		"Write code that is easy to delete, not easy to extend. — Dave Cheney",
		"Go developers don’t ask 'which framework should I use?' — Go is the framework. — Anonymous",
		"If Go didn’t exist, someone would have to invent it. — Anonymous",
		"Go doesn’t have inheritance, but it has composition, which is better. — Anonymous",
		"Go’s standard library is one of its greatest strengths. — Anonymous",
		"Don’t fight Go’s garbage collector, work with it. — Anonymous",
		"Readable code is maintainable code. — Anonymous",
		"The Go community prefers convention over configuration. — Anonymous",
		"In Go, less boilerplate means more productivity. — Anonymous",
		"Go’s error handling is verbose, but explicit is better than implicit. — Anonymous",
		"Learning Go is easy; mastering its simplicity is hard. — Anonymous",
		"Go is opinionated in the best way possible. — Anonymous",
		"Go is not just a language, it’s a philosophy. — Anonymous",
	}
)

func Start() {
	bot, err := discordgo.New("Bot " + config.Config.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	// start a websocket connection to listen to all the messages
	err = bot.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	user, err := bot.User("@me")
	if err != nil {
		log.Fatal("error getting current user,", err)
	}

	botId = user.ID

	bot.AddHandler(messageHandler)

	bot.Identify.Intents = discordgo.IntentGuildMessages
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore the messages written by the bot itself
	if m.Author.ID == botId {
		return
	}
	log.Println("message received", m.Content)

	// if m.Content == "ping" || m.Content == "pong" {
	s.ChannelMessageSend(m.ChannelID, golangQuotes[rand.IntN(len(golangQuotes))])
	// }
}
