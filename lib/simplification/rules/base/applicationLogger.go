package base

import (
	"strings"
	"sync"
)

type ApplicationLogger interface {
	LogApplication(ruleName, description, before, after string)
	Applications() []RuleApplication
	Clear()
	String(verbose bool) string
}

type BasicApplicationLogger struct {
	applications []RuleApplication
	lock         sync.RWMutex
}

func NewBasicApplicationLogger() *BasicApplicationLogger {
	return &BasicApplicationLogger{
		applications: make([]RuleApplication, 0),
	}
}

func (l *BasicApplicationLogger) LogApplication(ruleName, description, before, after string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.applications = append(l.applications, RuleApplication{
		Name:        ruleName,
		Description: description,
		Before:      before,
		After:       after,
	})
}

func (l *BasicApplicationLogger) Applications() []RuleApplication {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.applications
}

func (l *BasicApplicationLogger) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.applications = make([]RuleApplication, 0)
}

func (l *BasicApplicationLogger) String(verbose bool) string {
	l.lock.RLock()
	defer l.lock.RUnlock()

	result := strings.Builder{}
	for _, record := range l.applications {
		if verbose {
			result.WriteString(record.VerboseString())
		} else {
			result.WriteString(record.String())
		}
	}

	return result.String()
}
