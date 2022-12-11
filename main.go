package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitmyth/pmouse/messages"
	"net"
	"net/http"
	"strings"

	"github.com/go-vgo/robotgo"

	_ "embed"

	"github.com/gorilla/websocket"
	qrcode "github.com/skip2/go-qrcode"
)

//go:embed websockets.html
var websocketHTML string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func localNetAddr() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			if strings.Contains(a.String(), "192.168") {
				return a.String(), nil
			}
		}
	}
	return "", errors.New("not found localAddress")
}

func main() {
	screenWidth, screenHeight := robotgo.GetScreenSize()

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var event messages.Event
			json.Unmarshal(msg, &event)

			switch event.Kind {
			case messages.MouseMoveRelative:
				x, y := robotgo.GetMousePos()
				fmt.Println(x, y)

				if x+event.X < 0 {
					continue
				}
				if x+event.X > screenWidth {
					continue
				}
				if y+event.Y < 0 {
					continue
				}
				if y+event.Y > screenHeight {
					continue
				}
				robotgo.MoveRelative(event.X, event.Y)

			case messages.ClickLeft:
				robotgo.Click("left", false)
			case messages.ClickRight:
				robotgo.Click("right", false)
			case messages.KeyDown:
				robotgo.KeyDown(event.Key)
			case messages.KeyUp:
				robotgo.KeyUp(event.Key)
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// http.ServeFile(w, r, "websockets.html")
		fmt.Fprint(w, websocketHTML)
	})

	host, _ := localNetAddr()
	host = strings.Split(host, "/")[0]
	qrurl := fmt.Sprintf("%s%s:%d", "http://", host, 8080)
	fmt.Printf("connect url: %s\n", qrurl)
	qr, _ := qrcode.New(qrurl, qrcode.Medium)
	println(qr.ToSmallString(false))
	http.ListenAndServe(":8080", nil)
}
