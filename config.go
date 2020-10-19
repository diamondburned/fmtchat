package main

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/discord"
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

	var lines = strings.Split(m.Content, "\n")
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
