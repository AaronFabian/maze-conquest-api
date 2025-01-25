package handlers

import (
	"crypto/sha256"
	"fmt"
	"maze-conquest-api/module/chat"
	"maze-conquest-api/module/webrtc"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func RoomCreate(ctx *fiber.Ctx) error {
	return ctx.Redirect("/api/v1/room/" + uuid.New().String())
}

func Room(ctx *fiber.Ctx) error {
	uuid := ctx.Params("uuid")
	if uuid == "" {
		ctx.Status(400)
		return nil
	}

	ws := "ws"
	if os.Getenv("MODE") == "prod" {
		ws = "wss"
	}

	uuid, suuid, _ := createOrGetRoom(uuid)
	return ctx.JSON(fiber.Map{
		"RoomWebsocketAddr":   fmt.Sprintf("%s://%s/api/v1/room/%s/websocket", ws, ctx.Hostname(), uuid),
		"RoomLink":            fmt.Sprintf("%s://%s/api/v1/room/%s", ctx.Protocol(), ctx.Hostname(), uuid),
		"ChatWebsocketAddr":   fmt.Sprintf("%s://%s/api/v1/room/%s/chat/websocket", ws, ctx.Hostname(), uuid),
		"ViewerWebsocketAddr": fmt.Sprintf("%s://%s/api/v1/room/%s/viewer/websocket", ws, ctx.Hostname(), uuid),
		"StreamLink":          fmt.Sprintf("%s://%s/api/v1/stream/%s", ctx.Protocol(), ctx.Hostname(), suuid),
		"Type":                "room",
	})
}

func RoomWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	_, _, room := createOrGetRoom(uuid)
	webrtc.RoomConn(c, room.Peers)
}

func createOrGetRoom(uuid string) (string, string, *webrtc.Room) {
	webrtc.RoomsLock.Lock()
	defer webrtc.RoomsLock.Unlock()

	h := sha256.New()
	h.Write([]byte(uuid))
	suuid := fmt.Sprintf("%x", h.Sum(nil))

	// If Room already there
	if room := webrtc.Rooms[uuid]; room != nil {
		if len(room.Hub.Clients) > 0 {
			if _, ok := webrtc.Streams[suuid]; !ok {
				webrtc.Streams[suuid] = room
			}
			return uuid, suuid, room
		} else {
			delete(webrtc.Rooms, uuid)
		}
	}

	// If room not found / Create new room
	hub := chat.NewHub(uuid)
	p := &webrtc.Peers{}
	// p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &webrtc.Room{
		Peers: p,
		Hub:   hub,
	}

	webrtc.Rooms[uuid] = room
	webrtc.Streams[suuid] = room

	go hub.Run()
	return uuid, suuid, room
}

func RoomViewerWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	webrtc.RoomsLock.Lock()
	if peer, ok := webrtc.Rooms[uuid]; ok {
		webrtc.RoomsLock.Unlock()
		roomViewerConn(c, peer.Peers)
		return
	}
	webrtc.RoomsLock.Unlock()
}

func roomViewerConn(c *websocket.Conn, p *webrtc.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		<-ticker.C
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
	}
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
