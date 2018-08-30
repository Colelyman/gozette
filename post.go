package main

import (
	"bytes"
	"time"
)

const (
	timeZone   = "MST"
	timeFormat = time.RFC822
)

func writeTomlHugoHeader(entry *Entry, buff bytes.Buffer) {
	location, _ := time.LoadLocation(timeZone)
	t := time.Now().In(location).Format(timeFormat)
	// write the front matter in toml format
	buff.WriteString("+++\n")
	buff.WriteString("title = \"\"\n")
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
}

func WriteHugoPost(entry *Entry) (string, string) {
	var buff bytes.Buffer

	writeTomlHugoHeader(entry, buff)

	// write the content
	buff.WriteString(entry.Content + "\n")

	path := entry.Slug + ".md"

	return "site/content/micro/" + path, buff.String()
}
