package usecase

import (
	"strings"
)

func buildMailAddress(name string, senders []string) map[string][]string {
	newSenders := []string{}
	for _, v := range senders {
		if strings.Trim(v, " ") != "" {
			newSenders = append(newSenders, v)
		}
	}
	_senders := make(map[string][]string)
	_senders[name] = newSenders
	return _senders
}
