package syslogalert

import (
	"fmt"
	"log"
)

type Monitor struct {
	Parse         ParseFunc
	AlertCallback func(message Message) error
	HostRules     HostRules
	TagRules      TagRules
	ContentRules  []*Rule
	hostChan      chan (Log)
	tagChan       chan (Log)
	contentChan   chan (Log)
}

func NewMonitor(parse ParseFunc, handler func(Message) error) *Monitor {
	return &Monitor{
		Parse:         parse,
		AlertCallback: handler,
		HostRules:     make(HostRules),
		TagRules:      make(TagRules),
		ContentRules:  []*Rule{},
		hostChan:      make(chan Log),
		tagChan:       make(chan Log),
		contentChan:   make(chan Log),
	}
}

func (m *Monitor) Add(p *Policy) {
	for _, rule := range *p {
		if rule.Host != "" {
			m.AddHostRule(rule)
		} else if rule.Tag != "" {
			m.AddTagRule(rule)
		} else if rule.Content != "" {
			m.AddContentRule(rule)
		}
		// No content, do nothing
	}
}

func (m *Monitor) AddHostRule(r *Rule) {
	if r.Host == "" {
		return
	}
	hostRule, foundHost := m.HostRules[r.Host]
	if !foundHost {
		hostRule = HostRule{make(TagRules), []*Rule{}}
		m.HostRules[r.Host] = hostRule
	}
	if r.Tag != "" {
		trs := hostRule.TagRules
		rules, found := trs[r.Tag]
		if !found {
			rules = []*Rule{}
			trs[r.Tag] = rules
		}
		trs[r.Tag] = append(rules, r)
	} else {
		hostRule.ContentRules = append(hostRule.ContentRules, r)
	}
	m.HostRules[r.Host] = hostRule
}

func (m *Monitor) AddTagRule(r *Rule) {
	if r.Tag != "" {
		rules, found := m.TagRules[r.Tag]
		if !found {
			rules = []*Rule{}
			m.TagRules[r.Tag] = rules
		}
		m.TagRules[r.Tag] = append(rules, r)
	}
}

func (m *Monitor) AddContentRule(r *Rule) {
	if r.Content != "" {
		m.ContentRules = append(m.ContentRules, r)
	}
}

func (m *Monitor) Check(line string) error {
	l, parseErr := m.Parse(line)
	if parseErr != nil {
		return parseErr
	}
	m.hostChan <- *l
	m.tagChan <- *l
	m.contentChan <- *l
	return nil
}

func (m *Monitor) Start() {
	go m.MonitorHostRules()
	go m.MonitorTagRules()
	go m.MonitorContentRules()
}

func (m *Monitor) MonitorHostRules() {
	for l := range m.hostChan {
		hostRule, foundRules := m.HostRules[l.Host]
		if !foundRules {
			continue
		}
		// Note: we could have both a tag and content match
		// Check tag rules
		tagRules, foundTag := hostRule.TagRules[l.Tag]
		if foundTag {
			for _, rule := range tagRules {
				if rule.MatchContent(l) {
					m.HandleMatch(rule, l)
				}
			}
		}
		// Check content rules
		for _, rule := range hostRule.ContentRules {
			if rule.MatchContent(l) {
				m.HandleMatch(rule, l)
			}
		}
	}
}

func (m *Monitor) MonitorTagRules() {
	for l := range m.tagChan {
		tagRules, foundTag := m.TagRules[l.Tag]
		if foundTag {
			for _, rule := range tagRules {
				if rule.MatchContent(l) {
					m.HandleMatch(rule, l)
				}
			}
		}
	}
}

func (m *Monitor) MonitorContentRules() {
	for l := range m.contentChan {
		for _, rule := range m.ContentRules {
			if rule.MatchContent(l) {
				m.HandleMatch(rule, l)
			}
		}
	}
}

func (m *Monitor) HandleMatch(r *Rule, l Log) {
	body := fmt.Sprintf("\tTime: %s\n\tSource: %s on %s\n\tContent: %s", l.Timestamp, l.Tag, l.Host, l.Content)
	msg := Message{r.Description, body}
	err := m.AlertCallback(msg)
	if err != nil {
		log.Print("ERROR: %s", err)
	}
}
