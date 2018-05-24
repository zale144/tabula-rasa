package websocket

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"fmt"
)

var (
	upgrader = websocket.Upgrader{}
	// TODO move map to a separate goroutine
	notifReceivers = make(map[string]*websocket.Conn) // rename to avoid conflict in the same package
)

// the notification struct
type Notification struct {
	ReceiverID	string 	`json:"receiverId"`
	Type		string	`json:"type"`
	Message		string	`json:"message"`
	Blink		bool	`json:"blink"`
}

// handling the websocket connection
func HandleNotification(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	userID := ps.ByName("accountId")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	// add notification receiver to the map
	notifReceivers[userID] = ws
	return nil
}

// func for sending notification to clients
func SendNotification(receiverId, message, typ string, blink bool) error {
	notif := Notification{
		ReceiverID:		receiverId,
		Type:			typ,
		Message:		message,
		Blink:			blink,
	}
	// instantiate the receiver websocket
	receiverWS, find := notifReceivers[notif.ReceiverID]
	if !find {
		err := fmt.Errorf("Client %v is not connected", notif.ReceiverID)
		return err
	}
	// respond to the client
	err := receiverWS.WriteJSON(notif)
	if err != nil {
		receiverWS.Close()
		delete(notifReceivers, notif.ReceiverID)
		return err
	}
	return nil
}