package main

import (
	"flag"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/crowdmob/goamz/sqs"
	"github.com/andrew-d/go-termutil"
)

func main() {
	cmd := "receive"
	access := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	var b []byte
	var err error
	// If stdin input
	if !termutil.Isatty(os.Stdin.Fd()) {
		b, err = ioutil.ReadAll(os.Stdin)
		if err != nil { panic(err) }
		cmd = "send"
	}

	region := flag.String("r", "", "region")
	qName := flag.String("q", "", "queue name")
	del := flag.Bool("d", false, "delete message")
	maxNumberOfMessages := flag.Int("mN", 1, "maximum messages")
	flag.Parse()
	if *del { cmd = "delete" }

	c, err := sqs.NewFrom(access, secret, *region)
	if err != nil { panic(err) }
	q, err := c.GetQueue(*qName)
	if err != nil { panic(err) }

	if cmd == "send" {
		resp, err := q.SendMessage(string(b))
		if err != nil { panic(err) }
		b, err := json.Marshal(resp)
		if err != nil { panic(err) }
		os.Stdout.Write(b)
	} else if cmd == "receive" {
		resp, err := q.ReceiveMessage(*maxNumberOfMessages)
		if err != nil { panic(err) }
		b, err := json.Marshal(resp)
		if err != nil { panic(err) }
		os.Stdout.Write(b)
	} else if cmd == "delete" {
		m := &sqs.Message{ReceiptHandle: string(b)}
		resp, err := q.DeleteMessage(m)
		if err != nil { panic(err) }
		b, err := json.Marshal(resp)
		if err != nil { panic(err) }
		os.Stdout.Write(b)
	}
}
