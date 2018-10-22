package syslogalert

import (
	"errors"
	"fmt"
	"time"

	"github.com/eenblam/syslogparser/rfc3164"
)

type ParseFunc func(string) (*Log, error)

type Log struct {
	Timestamp time.Time
	//Priority // meh
	//Facility // meh
	//Severity // meh
	// RFC 3164 requires these fields to all be readable strings.
	Host    string
	Tag     string
	Content string
}

func ParseLog(line string) (*Log, error) {
	buff := []byte(line)
	p := rfc3164.NewParser(buff)
	err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read line: %s", err)
	}
	parts := p.Dump()
	// Extract individual parts
	maybeTimestamp, found := parts["timestamp"]
	if !found {
		return nil, errors.New("Log has no timestamp")
	}
	timestamp, parseOk := maybeTimestamp.(time.Time)
	if !parseOk {
		return nil, errors.New("Could not extract log timestamp")
	}
	maybeHost, found := parts["hostname"]
	if !found {
		return nil, errors.New("Log has no host")
	}
	host, parseOk := maybeHost.(string)
	if !parseOk {
		return nil, errors.New("Could not coerce log hostname to string")
	}
	maybeTag, found := parts["tag"]
	if !found {
		return nil, errors.New("Log has no tag")
	}
	tag, parseOk := maybeTag.(string)
	if !parseOk {
		return nil, errors.New("Could not coerce log tag to string")
	}
	maybeContent, found := parts["content"]
	if !found {
		return nil, errors.New("Log has no content")
	}
	content, parseOk := maybeContent.(string)
	if !parseOk {
		return nil, errors.New("Could not coerce log content to string")
	}
	return &Log{timestamp, host, tag, content}, nil
}
