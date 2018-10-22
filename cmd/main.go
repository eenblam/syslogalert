package main

import (
	"fmt"
	"io"
	"log"
	"os"

	sa "github.com/eenblam/syslog-alert"
	"github.com/papertrail/go-tail/follower"
)

func main() {
	config, configErr := sa.GetConfig("config.json")
	if configErr != nil {
		log.Fatalf("Could not get config file: %s", configErr)
	}
	mailer, mailErr := sa.NewMailer("smtp.json")
	if mailErr != nil {
		log.Fatal("Could not configure SMTP: %s", mailErr)
	}
	bm := sa.NewBufferedMailer(config.Timeout, mailer)
	go bm.Start()
	m := sa.NewMonitor(sa.ParseLog, bm.SendMessage)
	m.Add(&config.Policy)
	m.Start()

	t, err := follower.New("/var/log/syslog-ng/syslog.log", follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for line := range t.Lines() {
		checkErr := m.Check(line.String())
		if checkErr != nil {
			log.Printf("WARN: %s", checkErr)
		}
	}
	if t.Err() != nil {
		fmt.Fprintln(os.Stderr, t.Err())
	}
}

func printAlert(header, body string) error {
	fmt.Printf("%s: \n%s\n", header, body)
	return nil
}
