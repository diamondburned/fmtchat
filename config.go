package main

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mitchellh/go-wordwrap"
)

var messageColumns = []string{
	`[`,
	`{{.Offset}}`,
	`]`,
	`{{.Author.Username}}`,
	`{{.SplitContent 0 4}}`,
}

type Message struct {
	discord.Message
	LastTime time.Time
}

func (m Message) Offset() string {
	if m.LastTime.IsZero() {
		return m.Timestamp.Format(time.Kitchen)
	}

	dura := m.Timestamp.Time().Sub(m.LastTime)
	return "+" + dura.String()
}

func (m Message) SplitContent(indentLevels ...int) string {
	if len(indentLevels) == 0 {
		indentLevels = []int{0}
	}

	var content = m.Content
	if columns > 0 {
		content = wordwrap.WrapString(content, columns)
	}

	var lines = strings.Split(content, "\n")
	for i, line := range lines {
		var ilevel int
		if i < len(indentLevels) {
			ilevel = indentLevels[i]
		} else {
			ilevel = indentLevels[len(indentLevels)-1]
		}

		lines[i] = strings.Repeat("\t", ilevel) + line
	}

	return strings.Join(lines, "\n")
}
