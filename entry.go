package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	hashids "github.com/speps/go-hashids"
)

type PostType int

const (
	EntryPost PostType = iota + 1
	CardPost
	EventPost
	CitePost
)

type Entry struct {
	Content     string   `json:"content"`
	Name        string   `json:"name"`
	Categories  []string `json:"category"`
	Type        PostType `json:"type"`
	Slug        string   `json:"mp-slug"`
	Summary     string   `json:"summary"`
	In_reply_to string   `json:"in-reply-to"`
	Like_of     string   `json:"like-of"`
	Repost_of   string   `json:"repost-of"`
	hash        string
	token       string
}

func CreateEntry(contentType ContentType, body string) (*Entry, error) {
	if contentType == WWW_FORM {
		bodyValues, err := url.ParseQuery(body)
		if err != nil {
			return nil, err
		}
		return createEntryFromURLValues(bodyValues)
	} else if contentType == JSON {
		entry := new(Entry)
		err := json.Unmarshal([]byte(body), entry)
		return entry, err
	} else if contentType == MULTIPART {
		fmt.Println("Multipart content-type was detected")
		fmt.Printf("body is: %s\n", body)
		return nil, errors.New("Multipart content-type not implemented yet")
	} else {
		return nil, errors.New("Unsupported content-type")
	}
}

func createEntryFromURLValues(bodyValues url.Values) (*Entry, error) {
	if _, ok := bodyValues["content"]; ok {
		entry := new(Entry)
		entry.Content = bodyValues["content"][0]
		entry.hash = generateHash()
		if name, ok := bodyValues["name"]; ok {
			entry.Name = name[0]
		}
		if category, ok := bodyValues["category"]; ok {
			entry.Categories = category
		} else if categories, ok := bodyValues["category[]"]; ok {
			entry.Categories = categories
		} else {
			entry.Categories = nil
		}
		if slug, ok := bodyValues["mp-slug"]; ok && len(slug) > 0 && slug[0] != "" {
			entry.Slug = slug[0] + "-" + entry.hash
		} else {
			entry.Slug = entry.hash
		}
		if summary, ok := bodyValues["summary"]; ok {
			entry.Summary = summary[0]
		}
		if inReplyTo, ok := bodyValues["in-reply-to"]; ok {
			entry.In_reply_to = inReplyTo[0]
		}
		if likeOf, ok := bodyValues["like-of"]; ok {
			entry.Like_of = likeOf[0]
		}
		if repostOf, ok := bodyValues["repost-of"]; ok {
			entry.Repost_of = repostOf[0]
		}
		if token, ok := bodyValues["access_token"]; ok {
			entry.token = "Bearer " + token[0]
		}

		return entry, nil
	}
	return nil,
		errors.New("Error parsing the entry from URL Values")
}

func WriteEntry(entry *Entry) (string, error) {
	path, file := WriteHugoPost(entry)
	err := CommitEntry(path, file)
	if err != nil {
		return "", err
	}
	return "/micro/" + entry.Slug, nil
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
