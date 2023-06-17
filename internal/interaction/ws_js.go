//go:build js && wasm
// +build js,wasm

package interaction

import (
	"context"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type JSWebSocket struct {
	ConnType int
	conn     *websocket.Conn
	sendConn *websocket.Conn
}

func (w *JSWebSocket) LocalAddr() string {
	return ""
}

func NewWebSocket(connType int) *JSWebSocket {
	return &JSWebSocket{ConnType: connType}
}

func (w *JSWebSocket) Close() error {
	return w.conn.Close(websocket.StatusNormalClosure, "Actively close the conn have old conn")
}

func (w *JSWebSocket) WriteMessage(messageType int, message []byte) error {
	w.setSendConn(w.conn)
	return w.conn.Write(context.Background(), websocket.MessageType(messageType), message)
}
func (w *JSWebSocket) setSendConn(sendConn *websocket.Conn) {
	w.sendConn = sendConn
}
func (w *JSWebSocket) ReadMessage() (int, []byte, error) {
	messageType, b, err := w.conn.Read(context.Background())
	return int(messageType), b, err
}

func (w *JSWebSocket) SetReadTimeout(timeout int) error {
	return nil
}

func (w *JSWebSocket) SetWriteTimeout(timeout int) error {
	return nil
}

func (w *JSWebSocket) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, httpResp, err := websocket.Dial(ctx, urlStr, nil)
	if err == nil {
		w.conn = conn
	}
	return httpResp, err
}

func (w *JSWebSocket) IsNil() bool {
	if w.conn != nil {
		return false
	}
	return true
}

func (w *JSWebSocket) SetConnNil() {
	w.conn = nil
}
func (w *JSWebSocket) CheckSendConnDiffNow() bool {
	return w.sendConn == w.conn
}
