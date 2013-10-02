package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/cupcake/apns"
)

func main() {
	cert, err := tls.LoadX509KeyPair("push.crt", "push.key")
	if err != nil {
		log.Fatal(err)
	}
	client, err := apns.DialApple(true, cert)
	if err != nil {
		log.Fatal(err)
	}
	token, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	notification := &apns.Notification{
		Token:   token,
		Payload: []byte(fmt.Sprintf(`{"aps":{"alert":%q}}`, os.Args[2])),
	}
	err = client.Send(notification)
	if err != nil {
		log.Fatal(err)
	}
}
