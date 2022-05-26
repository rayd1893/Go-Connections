package v1

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

type Event struct {
	Message      map[string]interface{} `json:"message,omitempty"`
	Subscription string                 `json:"subscription,omitempty"`
}

func (e *Event) GetDecodedData() (map[string]interface{}, error) {
	data, isPresented := e.Message["data"]
	if !isPresented {
		return nil, errors.New("data in message is not presented")
	}

	dec, err := base64.StdEncoding.DecodeString(data.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	decoded := make(map[string]interface{})
	if err := json.Unmarshal(dec, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}
