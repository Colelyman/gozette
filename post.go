package main

import (
	"bytes"
	"time"
)

const (
	timeZone   = "MST"
	timeFormat = time.RFC822
)

func writeTomlHugoHeader(entry *Entry) string {
	var buff bytes.Buffer

	location, _ := time.LoadLocation(timeZone)
	t := time.Now().In(location).Format(timeFormat)
	// write the front matter in toml format
	buff.WriteString("+++\n")
	if len(entry.Name) == 0 {
		buff.WriteString("title = \"\"\n")
	} else {
		buff.WriteString("title = \"" + entry.Name + "\"\n")
	}
	buff.WriteString("date = \"" + t + "\"\n")
	buff.WriteString("categories = [\"Micro\"]\n")
	buff.WriteString("tags = [")
	for i, tag := range entry.Categories {
		buff.WriteString("\"" + tag + "\"")
		if i < len(entry.Categories)-1 {
			buff.WriteString(", ")
		}
	}
	buff.WriteString("]\n")
	buff.WriteString("slug = \"" + entry.Slug + "\"\n")
	buff.WriteString("+++\n")

	return buff.String()
}

func WriteHugoPost(entry *Entry) (string, string) {
	var buff bytes.Buffer

	buff.WriteString(writeTomlHugoHeader(entry))

	if len(entry.In_reply_to) > 0 {
		buff.WriteString("↪️ replying to: " + entry.In_reply_to + "\n")
	}
	if len(entry.Like_of) > 0 {
		buff.WriteString("👍: " + entry.Like_of + "\n")
	}
	if len(entry.Repost_of) > 0 {
		buff.WriteString("🔁 repost of: " + entry.Repost_of + "\n")
	}
	if len(entry.Content) > 0 {
		buff.WriteString(entry.Content + "\n")
	}

	path := entry.Slug + ".md"

	return "site/content/micro/" + path, buff.String()
}
