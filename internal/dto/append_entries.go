package dto

import (
	"github.com/svartvalp/fqw-consensus/internal/log"
)

type AppendEntriesQuery struct {
	Term         int64       `json:"term,omitempty"`
	LeaderID     int64       `json:"leader_id,omitempty"`
	PrevLogIndex int64       `json:"prev_log_index,omitempty"`
	PrevLogTerm  int64       `json:"prev_log_term,omitempty"`
	Entries      []log.Entry `json:"entries,omitempty"`
	LeaderCommit int64       `json:"leader_commit,omitempty"`
}

type AppendEntriesResponse struct {
	Term    int64 `json:"term,omitempty"`
	Success bool  `json:"success,omitempty"`

	ConflictIndex int64 `json:"conflict_index"`
	ConflictTerm  int64 `json:"conflict_term"`
}
