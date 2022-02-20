package consensus

import (
	"time"
)

type State struct {
	Store           map[string]string `json:"store"`
	LocalID         int64             `json:"local_id,omitempty"`
	CurrentTerm     int64             `json:"current_term,omitempty"`
	VotedFor        *int64            `json:"voted_for,omitempty"`
	State           CMState           `json:"state,omitempty"`
	Votes           int               `json:"votes"`
	ElectionTimeout *time.Time        `json:"election_timeout,omitempty"`
	CommitIndex     int64             `json:"commit_index"`
}
