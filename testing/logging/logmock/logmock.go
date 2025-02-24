package logmock

import (
	"fmt"
	"sync"

	"github.com/sjdaws/pkg/logging"
)

type LogMock struct {
	journal     []map[string]string
	lastLevel   string
	lastMessage string
	mutex       *sync.Mutex
}

func New() *LogMock {
	return &LogMock{
		journal:     make([]map[string]string, 0),
		lastLevel:   "",
		lastMessage: "",
		mutex:       &sync.Mutex{},
	}
}

func (l *LogMock) Debug(message any, replacements ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.append("debug", message, replacements...)
}

func (l *LogMock) Error(message any, replacements ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.append("error", message, replacements...)
}

func (l *LogMock) Fatal(message any, replacements ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.append("fatal", message, replacements...)

	panic("fatal log received")
}

func (l *LogMock) GetAllLogs() []map[string]string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.journal
}

func (l *LogMock) GetLastLevel() string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.lastLevel
}

func (l *LogMock) GetLastMessage() string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.lastMessage
}

func (l *LogMock) Info(message any, replacements ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.append("info", message, replacements...)
}

func (l *LogMock) SetDepth(_ int) logging.Logger {
	return l
}

func (l *LogMock) SetVerbosity(_ logging.Verbosity) logging.Logger {
	return l
}

func (l *LogMock) Warn(message any, replacements ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.append("warn", message, replacements...)
}

func (l *LogMock) append(level string, message any, replacements ...any) {
	combined := fmt.Sprintf(fmt.Sprintf("%v", message), replacements...)

	l.lastLevel = level
	l.lastMessage = combined
	l.journal = append(l.journal, map[string]string{level: combined})
}
