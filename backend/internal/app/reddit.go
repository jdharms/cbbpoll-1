package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func usernameFromRedditToken(token string) (name string, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/v1/me", nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "cbbpoll_backend/0.1.0")

	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	var content []byte
	content, err = ioutil.ReadAll(resp.Body)
	data := make(map[string]interface{})
	err = json.Unmarshal(content, &data)
	if err != nil {
		return
	}

	name, ok := data["name"].(string)
	if !ok {
		return
	}

	return name, nil
}