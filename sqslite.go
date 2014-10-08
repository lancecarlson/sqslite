package main

import (
	"fmt"
	"flag"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/crowdmob/goamz/sqs"
)

var Usage = func(flags *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "Usage of %s [command]:\n", os.Args[0])
	flags.PrintDefaults()
	os.Exit(0)
}

func main() {
	cmd := ""
        if len(os.Args) > 1 {
                cmd = os.Args[1]
        }

	access := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	region := flags.String("r", "", "region")
	qName := flags.String("q", "", "queue name")
	maxNumberOfMessages := flags.Int("mN", 1, "maximum messages")

	if len(os.Args) > 1 {
		flags.Parse(os.Args[2:])
	}

	// Get Stdin
	var b []byte
	stat, err := os.Stdin.Stat()
	if err != nil { panic(err) }
	if stat.Size() > 0 {
		b, err = ioutil.ReadAll(os.Stdin)
		if err != nil { panic(err) }
		cmd = "send"
	} else {
		cmd = "receive"
	}

	c, err := sqs.NewFrom(access, secret, *region)
	if err != nil { panic(err) }
	q, err := c.GetQueue(*qName)
	if err != nil { panic(err) }

	if cmd == "send" {
		fmt.Println(string(b))
	} else if cmd == "receive" {
		resp, err := q.ReceiveMessage(*maxNumberOfMessages)
		if err != nil { panic(err) }
		b, err := json.Marshal(resp)
		if err != nil { panic(err) }
		os.Stdout.Write(b)
	}
}