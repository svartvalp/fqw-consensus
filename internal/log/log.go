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
	GetAll() []Entry
}

type log struct {
	logs []Entry
}

func (l *log) GetAll() []Entry {
	return l.logs
}

func (l *log) GetFrom(index int64) []Entry {
	if index == 0 {
		return make([]Entry, 0)
	}
	logInd := -1
	for i, log := range l.logs {
		if log.Index == index {
			logInd = i
		}
	}
	if logInd == -1 {
		return make([]Entry, 0)
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
	var lastLogIndex int64
	if len(l.logs) > 0 {
		lastLog := l.logs[len(l.logs)-1]
		lastLogIndex = lastLog.Index
	}
	for _, entry := range entries {
		if entry.Index < lastLogIndex {
			return errors.New("entry with low index found")
		}
		l.logs = append(l.logs, entry)
		lastLogIndex = entry.Index
	}
	return nil
}

func New() Log {
	return &log{
		logs: make([]Entry, 0),
	}
}
