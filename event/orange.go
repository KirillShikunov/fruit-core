package event

import "encoding/json"

type OrangeEvent struct {
	Size int `json:"size"`
}

func (s *OrangeEvent) Name() string {
	return "fruit.orange"
}

func (s *OrangeEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, s)
}
