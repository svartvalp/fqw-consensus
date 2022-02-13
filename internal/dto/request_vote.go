package dto

type RequestVoteQuery struct {
	Term         int64 `json:"term,omitempty"`
	CandidateID  int64 `json:"candidate_id,omitempty"`
	LastLogIndex int64 `json:"last_log_index,omitempty"`
	LastLogTerm  int64 `json:"last_log_term,omitempty"`
}

type RequestVoteResponse struct {
	Term    int64 `json:"term,omitempty"`
	Granted bool  `json:"granted,omitempty"`
}
