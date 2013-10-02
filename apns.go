package apns

import (
	"crypto/tls"
	"errors"
	"net"
)

type Client struct {
	conn net.Conn
	cert tls.Certificate
	addr string
}

func DialApple(sandbox bool, cert tls.Certificate) (*Client, error) {
	if sandbox {
		return Dial("gateway.sandbox.push.apple.com:2195", cert)
	}
	return Dial("gateway.push.apple.com:2195", cert)
}

func Dial(addr string, cert tls.Certificate) (*Client, error) {
	client := &Client{cert: cert, addr: addr}
	conn, err := tls.Dial("tcp", addr, &tls.Config{Certificates: []tls.Certificate{cert}})
	client.conn = conn
	return client, err
}

func (c *Client) Send(notification *Notification) error {
	if len(notification.Payload) > 256 {
		return errors.New("apns: payload must not be longer than 256 bytes")
	}
	if len(notification.Token) != 32 {
		return errors.New("apns: device token must be exactly 32 bytes")
	}
	_, err := c.conn.Write(notification.Bytes())
	return err
}
