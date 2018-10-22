package syslogalert

import (
	"log"
	"regexp"
)

type Rule struct {
	Host        string
	Tag         string
	Content     string
	Description string
	Regex       bool
	//TODO be sure regexp is ignored in json
	regexp *regexp.Regexp
}

// MatchContent returns true if parts["contents"] matches the rule.
func (r *Rule) MatchContent(l Log) bool {
	if r.Regex && r.regexp != nil {
		return r.regexp.MatchString(l.Content)
	} else if r.Regex && r.regexp == nil {
		log.Print("WARN: Regex true but regexp is nil")
	}
	return r.Content == l.Content
}

type HostRules map[string]HostRule

type HostRule struct {
	TagRules     TagRules
	ContentRules []*Rule
}

type TagRules map[string][]*Rule
