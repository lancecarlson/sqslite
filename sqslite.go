package main

import (
	"flag"
	"os"
	"io/ioutil"
	"encoding/xml"
	"encoding/json"
	"github.com/crowdmob/goamz/sqs"
	"github.com/andrew-d/go-termutil"
)

func Format(format string, resp interface {}) ([]byte, error) {
	if format == "json" {
		return json.Marshal(resp)
	} else {
		return xml.Marshal(resp)
	}
}

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

	qName := flag.String("q", "", "queue name")
	region := flag.String("r", "", "region (ie: us-east-1)")
	del := flag.Bool("d", false, "delete message (send only)")
	format := flag.String("f", "xml", "response format (xml or json)")
	maxNumberOfMessages := flag.Int("mN", 1, "maximum messages")
	flag.Parse()
	if *qName == "" || *region == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *del { cmd = "delete" }

	c, err := sqs.NewFrom(access, secret, *region)
	if err != nil { panic(err) }
	q, err := c.GetQueue(*qName)
	if err != nil { panic(err) }

	if cmd == "send" {
		resp, err := q.SendMessage(string(b))
		if err != nil { panic(err) }
		b, err := Format(*format, resp)
		if err != nil { panic(err) }
		os.Stdout.Write(b)
	} else if cmd == "receive" {
		resp, err := q.ReceiveMessage(*maxNumberOfMessages)
		if err != nil { panic(err) }
		b, err := Format(*format, resp)
		if err != nil { panic(err) }
		os.Stdout.Write(b)
	} else if cmd == "delete" {
		m := &sqs.Message{ReceiptHandle: string(b)}
		resp, err := q.DeleteMessage(m)
		if err != nil { panic(err) }
		b, err := Format(*format, resp)
		if err != nil { panic(err) }
		os.Stdout.Write(b)
	}
}
