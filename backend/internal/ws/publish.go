package ws

import "encoding/json"

func BroadcastSeatEvent(hub *Hub, event SeatEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	hub.Broadcast <- payload
	return nil
}
