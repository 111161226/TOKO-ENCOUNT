package handler

import (
	"github.com/111161226/TOKO-ENCOUNT/db/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type webSocketPublisher struct {
	userIdConnectionPool map[string]map[*websocket.Conn]struct{}
}

type newMessageNotificationBody struct {
	RoomId  string         `json:"roomId"`
	Message *model.Message `json:"message"`
}

type newMessageNotification struct {
	Type string                     `json:"type"`
	Body newMessageNotificationBody `json:"body"`
}

func newWebSocketPublisher() *webSocketPublisher {
	return &webSocketPublisher{
		userIdConnectionPool: map[string]map[*websocket.Conn]struct{}{},
	}
}

//notify new received message
func (ws *webSocketPublisher) NotifyNewMessage(userIds []string, roomId string, message *model.Message) error {
	n := newMessageNotification{
		Type: "NEW_MESSAGE",
		Body: newMessageNotificationBody{
			RoomId:  roomId,
			Message: message,
		},
	}

	if roomId == "0" { //case : open chat
		for id, connections := range ws.userIdConnectionPool { // notify all loginned
			if id == message.UserId {
				//don't notify to posted user
				continue
			}

			for connection := range connections {
				err := websocket.JSON.Send(connection, n)
				if err != nil {
					return err
				}
			}
		}
	} else {
		for _, userId := range userIds { //notify designated user
			connections, ok := ws.userIdConnectionPool[userId]
			if !ok {
				continue
			}

			for connection := range connections {
				err := websocket.JSON.Send(connection, n)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (h *Handler) GetWebSocket(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	userId := sess.UserId

	//websocket setting
	server := websocket.Server{Handler: func(connection *websocket.Conn) {
		defer connection.Close()

		connections, ok := h.ws.userIdConnectionPool[userId]
		if !ok {
			connections = map[*websocket.Conn]struct{}{}
			h.ws.userIdConnectionPool[userId] = connections
		}
		connections[connection] = struct{}{}

		for {
			message := ""
			err = websocket.Message.Receive(connection, &message)
			if err != nil {
				break
			}
		}
		delete(connections, connection)
		if len(connections) == 0 {
			delete(h.ws.userIdConnectionPool, userId)
		}
	}}

	server.ServeHTTP(c.Response(), c.Request())
	return nil
}
