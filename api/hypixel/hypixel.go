package hypixel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	BaseUrl = "https://api.hypixel.net"
)

type API struct {
	Key string `json:"key"`
}

func (api *API) Player(id string) (*Player, error) {
	response, err := http.Get(fmt.Sprintf("%s/player?key=%s&uuid=%s", BaseUrl, api.Key, id))
	if err != nil {
		return nil, err
	}

	// Read the response body.
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	// Decode the response body into a JSON object.
	var responseData *playerResponse
	err = json.Unmarshal(contents, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData.Player, nil
}
