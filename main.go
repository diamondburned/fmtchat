package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
)

var tmpls = make([]*template.Template, len(messageColumns))

var funcs = template.FuncMap{}

// Flags.
var (
	columns   = uint(80)
	limit     = uint(100)
	channelID = uint64(0)
	beforeID  = uint64(0)
	afterID   = uint64(0)
)

func init() {
	flag.UintVar(&columns, "cols", columns, "maximum width per column before wrapping")
	flag.UintVar(&limit, "l", limit, "maximum number of messages to fetch")
	flag.Uint64Var(&channelID, "channelID", 0, "channel ID")
	flag.Uint64Var(&beforeID, "beforeID", 0, "message to fetch before, optional")
	flag.Uint64Var(&afterID, "afterID", 0, "message to fetch after, optional")

	for i, column := range messageColumns {
		tmpls[i] = template.Must(
			template.New(strconv.Itoa(i)).Funcs(funcs).Parse(column),
		)
	}
}

func main() {
	flag.Parse()

	var token = os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("Missing $TOKEN")
	}

	c := api.NewClient(token)

	var msgs []discord.Message
	var err error

	switch {
	case beforeID > 0:
		msgs, err = c.MessagesBefore(
			discord.ChannelID(channelID), discord.MessageID(beforeID), limit,
		)
	case afterID > 0:
		msgs, err = c.MessagesAfter(
			discord.ChannelID(channelID), discord.MessageID(afterID), limit,
		)
	default:
		msgs, err = c.Messages(
			discord.ChannelID(channelID), limit,
		)
	}

	if err != nil {
		log.Fatalln("Failed to get messages:", err)
	}

	if len(msgs) == 0 {
		return
	}

	log.Println("Found", len(msgs), "messages.")

	tabber := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)
	buffer := bytes.Buffer{} // scratch buffer for template execution

	var lastTime time.Time

	// Iterate backwards, that is from earliest to latest.
	for i := len(msgs) - 1; i != 0; i-- {
		var message = Message{
			Message:  msgs[i],
			LastTime: lastTime,
		}

		for i, tmpl := range tmpls {
			if err := tmpl.Execute(&buffer, message); err != nil {
				log.Fatalln("Template failed:", err)
			}

			if i != len(tmpls)-1 {
				buffer.WriteByte('\t')
			} else {
				buffer.WriteByte('\n')
			}

			buffer.WriteTo(tabber)
			buffer.Reset()
		}

		lastTime = message.Timestamp.Time()
	}

	if err := tabber.Flush(); err != nil {
		log.Println("Tabwriter Flush failed:", err)
	}
}
