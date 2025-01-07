package handlers

import (
	"fmt"
	"maze-conquest-api/module/webrtc"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Stream(c *fiber.Ctx) error {
	suuid := c.Params("suuid")
	if suuid == "" {
		c.Status(400)
		return nil
	}

	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		ws = "wss"
	}

	webrtc.RoomsLock.Lock()
	if _, ok := webrtc.Streams[suuid]; ok {
		webrtc.RoomsLock.Unlock()
		return c.Render("stream", fiber.Map{
			"StreamWebsocketAddr": fmt.Sprintf("%s://%s/api/v1/stream/%s/websocket", ws, c.Hostname(), suuid),
			"ChatWebsocketAddr":   fmt.Sprintf("%s://%s/api/v1/stream/%s/chat/websocket", ws, c.Hostname(), suuid),
			"ViewerWebsocketAddr": fmt.Sprintf("%s://%s/api/v1/stream/%s/viewer/websocket", ws, c.Hostname(), suuid),
			"Type":                "stream",
		}, "layouts/main")
	}
	webrtc.RoomsLock.Unlock()

	return c.Render("stream", fiber.Map{
		"NoStream": "true",
		"Leave":    "true",
	}, "layouts/main")
}

func StreamWebsocket(c *websocket.Conn) {
	// suuid := c.Params("suuid")
	// if suuid == "" {
	// 	return
	// }

	// webrtc.RoomsLock.Lock()
	// if stream, ok := webrtc.Streams[suuid]; ok {
	// 	webrtc.RoomsLock.Unlock()
	// 	webrtc.StreamConn(c, stream.Peers)
	// 	return
	// }
	// webrtc.RoomsLock.Unlock()
}

func StreamViewerWebsocket(c *websocket.Conn) {
	suuid := c.Params("suuid")
	if suuid == "" {
		return
	}

	webrtc.RoomsLock.Lock()
	if stream, ok := webrtc.Streams[suuid]; ok {
		webrtc.RoomsLock.Unlock()
		viewerConn(c, stream.Peers)
		return
	}
	webrtc.RoomsLock.Unlock()
}

func viewerConn(c *websocket.Conn, p *webrtc.Peers) {
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
