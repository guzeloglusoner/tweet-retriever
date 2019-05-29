package main

import "encoding/json"

// Message struct is used for storing the event as a Json object.
type Message struct {
	Username string      `json:"user"`
	Data     interface{} `json:"data"`
	Date     string      `json:"date"`
	Type     string      `json:"type"`
}

// UnMarshalMessage function is used to create an Event type from input
// Unmarshales input into Event type
func UnMarshalMessage(input []byte) (*Message, error) {
	event := new(Message)
	err := json.Unmarshal(input, event)
	return event, err
}

// Marshal method is used for marshaling event type
func (e *Message) marshal() []byte {
	output, _ := json.Marshal(e)
	return output
}
