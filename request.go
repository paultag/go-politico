package politico

import (
	"encoding/json"
	"net/http"
)

func News() ([]Story, error) {
	resp, err := http.Get("http://www.politico.com/feeds/ipad/headline_stories.json")
	if err != nil {
		return nil, err
	}

	stories := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&stories); err != nil {
		return nil, err
	}

	return stories.Stories.Stories, nil
}
