package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"os"

	"github.com/crowdmob/goamz/sqs"
)

func Format(format string, resp interface{}) ([]byte, error) {
	if format == "json" {
		return json.Marshal(resp)
	}
	return xml.Marshal(resp)
}

func main() {
	access := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	cmd := flag.String("c", "r", "command (r=receive, s=send, d=delete)")
	qName := flag.String("q", "", "queue name")
	region := flag.String("re", "us-east-1", "region")
	format := flag.String("f", "xml", "response format (xml or json)")
	maxNumberOfMessages := flag.Int("mN", 1, "maximum messages")
	flag.Parse()

	// Check required environment variables
	if access == "" {
		panic("AWS_ACCESS_KEY_ID is undefined")
	}

	if secret == "" {
		panic("AWS_SECRET_ACCESS_KEY is undefined")
	}

	// If required flags are are not filled
	if *qName == "" || *region == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var b []byte
	var err error
	if *cmd == "s" || *cmd == "d" {
		b, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	}

	c, err := sqs.NewFrom(access, secret, *region)
	if err != nil {
		panic(err)
	}
	q, err := c.GetQueue(*qName)
	if err != nil {
		panic(err)
	}

	if *cmd == "s" {
		resp, err := q.SendMessage(string(b))
		if err != nil {
			panic(err)
		}
		b, err := Format(*format, resp)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(b)
	} else if *cmd == "r" {
		resp, err := q.ReceiveMessage(*maxNumberOfMessages)
		if err != nil {
			panic(err)
		}
		b, err := Format(*format, resp)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(b)
	} else if *cmd == "d" {
		m := &sqs.Message{ReceiptHandle: string(b)}
		resp, err := q.DeleteMessage(m)
		if err != nil {
			panic(err)
		}
		b, err := Format(*format, resp)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(b)
	} else {
		flag.PrintDefaults()
		panic("Invalid command")
	}
}
