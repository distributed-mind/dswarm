// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package dht

import (
	"log"
	"sync"
	"time"
	"io"
	"github.com/gorilla/websocket"
)

const (
	connectionTimeout = 30 * time.Second
)

type (
	conn struct {
        websocket *websocket.Conn
		write sync.Mutex
		read sync.Mutex
		r  io.Reader
		w  io.WriteCloser
    }
)


func (c *conn) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	err := c.websocket.SetReadDeadline(time.Now().Add(connectionTimeout))
	if err != nil {
		log.Println("read error: ", err)
	}
	c.read.Lock()
	defer c.read.Unlock()
	n := 0
	for n < len(p) {
		if c.r == nil {
			_, c.r, err = c.websocket.NextReader()
			if err != nil {
				return 0, err
			}
			n = 0
		}
		r := 0
		r, err = c.r.Read(p[n:])
		n += r
		if err == io.EOF {
			c.r = nil
			break
		} else if err != nil {
			log.Println("read error: ", err)
			break
		}
	}
	return n, err
}

func (c *conn) Write(p []byte) (int, error) {
	c.write.Lock()
	defer c.write.Unlock()
	// c.read.Lock()
	// defer c.read.Unlock()
	err := c.websocket.SetWriteDeadline(time.Now().Add(connectionTimeout))
	if err != nil {
		log.Println("write error: ", err)
	}
	err = c.websocket.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), err
}

func (c *conn) Close() error {
	c.write.Lock()
	c.read.Lock()
	// defer c.write.Unlock()
	// defer c.read.Unlock()
	return c.websocket.Close()
}