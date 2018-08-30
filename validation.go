package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type ContentType int

const (
	WWW_FORM ContentType = iota
	JSON
	UNSUPPORTED_TYPE
)

const (
	indieAuthTokenUrl = "https://tokens.indieauth.com/token"
	indieAuthMe       = "http://colelyman.com/"
)

type IndieAuthRes struct {
	Me       string `json:"me"`
	ClientId string `json:"client_id"`
	Scope    string `json:"scope"`
	Issue    int    `json:"issued_at"`
	Nonce    int    `json:"nonce"`
}

func checkAccess(token string) (bool, error) {
	if token == "" {
		return false,
			errors.New("Token string is empty")
	}
	// form the request to check the token
	client := &http.Client{}
	req, err := http.NewRequest("GET", indieAuthTokenUrl, nil)
	if err != nil {
		return false,
			errors.New("Error making the request for checking token access")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	// send the request
	res, err := client.Do(req)
	if err != nil {
		return false,
			errors.New("Error sending the request for checking token access")
	}
	defer res.Body.Close()
	// parse the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false,
			errors.New("Error parsing the response for checking token access")
	}
	var indieAuthRes = new(IndieAuthRes)
	err = json.Unmarshal(body, &indieAuthRes)
	if err != nil {
		return false,
			errors.New("Error parsing the response into json for checking token access " + err.Error())
	}

	// verify results of the response
	if indieAuthRes.Me != indieAuthMe {
		return false,
			errors.New("Me does not match")
	}
	scopes := strings.Fields(indieAuthRes.Scope)
	postPresent := false
	for _, scope := range scopes {
		if scope == "post" || scope == "create" || scope == "update" {
			postPresent = true
			break
		}
	}
	if !postPresent {
		return false,
			errors.New("Post is not present in the scope")
	}
	return true, nil
}

func CheckAuthorization(entry *Entry, headers map[string]string) bool {
	token, ok := headers["authorization"]
	if !ok {
		token = entry.token
	} else {
		return false
	}

	if ok, err := checkAccess(token); ok {
		return true
	} else if err != nil {
		return false
	} else {
		return false
	}
}

func GetContentType(headers map[string]string) (ContentType, error) {
	if contentType, ok := headers["content-type"]; ok {
		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			return WWW_FORM, nil
		}
		if strings.Contains(contentType, "application/json") {
			return JSON, nil
		}
		return UNSUPPORTED_TYPE, errors.New("Content-type " + contentType + " is not supported, use application/x-www-form-urlencoded or application/json")
	}
	return UNSUPPORTED_TYPE, errors.New("Content-type is not provided in the request")
}
