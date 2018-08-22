package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"time"

	hashids "github.com/speps/go-hashids"
)

type Entry struct {
	Content string
	tags    []string
	slug    string
	hash    string
}

func CreateEntry(bodyValues url.Values) (string, error) {
	if _, ok := bodyValues["content"]; ok {
		fmt.Println(bodyValues)
		entry := new(Entry)
		entry.Content = bodyValues["content"][0]
		entry.hash = generateHash()
		if tags, ok := bodyValues["category"]; ok {
			entry.tags = tags
		} else {
			entry.tags = nil
		}
		if slug, ok := bodyValues["mp-slug"]; ok && len(slug) > 0 && slug[0] != "" {
			entry.slug = slug[0] + "-" + entry.hash
		} else {
			entry.slug = entry.hash
		}
		fmt.Printf("Hash value is %+v\n", entry.hash)

		// construct the post
		path, file, _ := writePost(entry)
		err := CommitEntry(path, file)
		if err != nil {
			return "", err
		}

		return "/micro/" + entry.slug, err
	}
	return "",
		errors.New("Content in response body is missing")
}

func generateHash() string {
	hd := hashids.NewData()
	hd.Salt = "do you want to know a secret?"
	h, _ := hashids.NewWithData(hd)
	t := []int{time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
	}
	id, _ := h.Encode(t)

	return id
}

func writePost(entry *Entry) (string, string, error) {
	var buff bytes.Buffer

	location, _ := time.LoadLocation("MST")
	t := time.Now().In(location).Format(time.RFC822)
	// write the front matter in toml format
	buff.WriteString("+++\n")
	buff.WriteString("title = \"\"\n")
	buff.WriteString("date = \"" + t + "\"\n")
	buff.WriteString("categories = [\"Micro\"]\n")
	buff.WriteString("tags = [")
	for i, tag := range entry.tags {
		buff.WriteString("\"" + tag + "\"")
		if i < len(entry.tags)-1 {
			buff.WriteString(", ")
		}
	}
	buff.WriteString("]\n")
	buff.WriteString("slug = \"" + entry.slug + "\"\n")
	buff.WriteString("+++\n")

	// write the content
	buff.WriteString(entry.Content + "\n")

	fmt.Printf("Length of slug is %d, with value %s.\n", len(entry.slug), entry.slug)
	// path := strings.Replace(entry.slug, " ", "-", -1) + ".md"
	path := entry.slug + ".md"
	fmt.Printf("path is %+v\n", path)

	return "site/content/micro/" + path, buff.String(), nil
}
