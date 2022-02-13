package log

import (
	"errors"
)

type Log interface {
	LastLog() *Entry
	GetLog(index int64) *Entry
	StoreLogs(entries []Entry) error
	DeleteFrom(index int64)
	GetFrom(index int64) []Entry
}

type log struct {
	logs []Entry
}

func (l *log) GetFrom(index int64) []Entry {
	if index == 0 {
		return make([]Entry, 0)
	}
	var logInd int
	for i, log := range l.logs {
		if log.Index == index {
			logInd = i
		}
	}
	return l.logs[logInd:]
}

func (l *log) DeleteFrom(index int64) {
	var logInd int
	for i, log := range l.logs {
		if log.Index == index {
			logInd = i
		}
	}
	l.logs = l.logs[:logInd]
}

func (l *log) LastLog() *Entry {
	if len(l.logs) == 0 {
		return nil
	}
	return &l.logs[len(l.logs)-1]
}

func (l *log) GetLog(index int64) *Entry {
	for _, log := range l.logs {
		if log.Index == index {
			return &log
		}
	}
	return nil
}

func (l *log) StoreLogs(entries []Entry) error {
	lastLog := l.logs[len(l.logs)-1]
	for _, entry := range entries {
		if entry.Index < lastLog.Index {
			return errors.New("entry with low index found")
		}
		l.logs = append(l.logs, entry)
		lastLog = entry
	}
	return nil
}

func New() Log {
	return &log{}
}
