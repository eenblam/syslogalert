package syslogalert

import (
	"encoding/json"
	"os"
	"regexp"
	"time"
)

type Policy []*Rule

type Config struct {
	Timeout time.Duration `json:MessageBufferWaitSeconds`
	Policy  Policy
}

func GetConfig(filename string) (*Config, error) {
	config := &Config{}
	rawFile, readErr := os.Open(filename)
	if readErr != nil {
		return nil, readErr
	}
	jsonParser := json.NewDecoder(rawFile)
	parseErr := jsonParser.Decode(config)
	if parseErr != nil {
		return nil, parseErr
	}
	for _, rule := range config.Policy {
		if rule.Regex {
			r, compileErr := regexp.Compile(rule.Content)
			if compileErr != nil {
				return nil, compileErr
			}
			rule.regexp = r
		}
	}
	return config, nil
}
