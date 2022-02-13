package consensus

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/svartvalp/fqw-consensus/internal/dto"
	"github.com/svartvalp/fqw-consensus/internal/log"
	"github.com/svartvalp/fqw-consensus/internal/rpc"
	"github.com/svartvalp/fqw-consensus/internal/store"
)

type CMState string

const (
	Follower  CMState = "Follower"
	Candidate CMState = "Candidate"
	Leader    CMState = "Leader"
	Dead      CMState = "Dead"
)

type Module interface {
	RequestVote(query dto.RequestVoteQuery) (dto.RequestVoteResponse, error)
	AppendEntries(query dto.AppendEntriesQuery) (dto.AppendEntriesResponse, error)
	GetState() interface{}
	Start()
}

type module struct {
	log         log.Log
	store       store.Store
	proxy       rpc.Proxy
	LocalID     int64   `json:"local_id,omitempty"`
	CurrentTerm int64   `json:"current_term,omitempty"`
	VotedFor    *int64  `json:"voted_for,omitempty"`
	State       CMState `json:"state,omitempty"`
	Votes       int     `json:"votes"`

	mu *sync.Mutex

	ElectionTimeout *time.Time `json:"election_timeout,omitempty"`
	tick            *time.Ticker

	commitIndex int64

	nextIndex  map[int64]int64
	matchIndex map[int64]int64
}

func (m *module) RequestVote(query dto.RequestVoteQuery) (dto.RequestVoteResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Println("incoming request vote")

	l := m.log.LastLog()
	var lastLogIndex int64
	var lastLogTerm int64
	if l != nil {
		lastLogIndex = l.Index
		lastLogTerm = l.Term
	}

	if query.Term > m.CurrentTerm ||
		(query.Term == m.CurrentTerm && (m.VotedFor == nil || *m.VotedFor == query.CandidateID) &&
			(query.LastLogTerm > lastLogTerm || (query.LastLogTerm == lastLogTerm && query.LastLogIndex >= lastLogIndex))) {
		fmt.Println(fmt.Sprintf("vote for %v, become follower", query.CandidateID))
		m.State = Follower
		m.CurrentTerm = query.Term
		m.VotedFor = &query.CandidateID
		m.resetElectionTimeout()

		return dto.RequestVoteResponse{
			Term:    m.CurrentTerm,
			Granted: true,
		}, nil
	}

	return dto.RequestVoteResponse{
		Term:    m.CurrentTerm,
		Granted: false,
	}, nil
}

func (m *module) AppendEntries(query dto.AppendEntriesQuery) (dto.AppendEntriesResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if query.Term < m.CurrentTerm {
		return dto.AppendEntriesResponse{
			Term:    m.CurrentTerm,
			Success: false,
		}, nil
	}

	if query.Term > m.CurrentTerm {
		m.CurrentTerm = query.Term
	}

	if m.State != Follower {
		m.State = Follower
	}

	m.resetElectionTimeout()

	lastLog := m.log.LastLog()
	prevLog := m.log.GetLog(query.PrevLogIndex)
	if query.PrevLogIndex > 0 {
		if prevLog == nil {
			return dto.AppendEntriesResponse{
				Term:          m.CurrentTerm,
				Success:       false,
				ConflictIndex: query.PrevLogIndex,
				ConflictTerm:  -1,
			}, nil
		}
		if prevLog.Term != query.PrevLogTerm {
			return dto.AppendEntriesResponse{
				Term:          m.CurrentTerm,
				Success:       false,
				ConflictIndex: prevLog.Index,
				ConflictTerm:  prevLog.Term,
			}, nil
		}
	}
	if lastLog != nil && query.PrevLogIndex < lastLog.Index {
		m.log.DeleteFrom(query.PrevLogIndex)
	}
	if len(query.Entries) > 0 {
		err := m.log.StoreLogs(query.Entries)
		if err != nil {
			panic(err)
		}
	}
	if m.commitIndex < query.LeaderCommit {
		m.commitIndex = query.LeaderCommit
		m.restore()
	}

	return dto.AppendEntriesResponse{
		Term:    m.CurrentTerm,
		Success: true,
	}, nil
}

func (m *module) GetState() interface{} {
	return m
}

func (m *module) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resetElectionTimeout()
	go m.RunElectionTimer()
	m.electSelf()
}

func (m *module) electSelf() {
	m.State = Candidate
	m.CurrentTerm += 1
	m.resetElectionTimeout()
	m.VotedFor = &m.LocalID
	m.Votes = 1

	for id, cl := range m.proxy.GetAllClients() {
		go func(id int64, cl rpc.Client) {
			if m.LocalID == id {
				return
			}
			m.mu.Lock()
			l := m.log.LastLog()
			var lastLogIndex int64
			var lastLogTerm int64
			if l != nil {
				lastLogIndex = l.Index
				lastLogTerm = l.Term
			}
			m.mu.Unlock()
			res, err := cl.RequestVote(dto.RequestVoteQuery{
				Term:         m.CurrentTerm,
				CandidateID:  m.LocalID,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			})
			if err != nil {
				fmt.Printf("request vote failed with err: %v \n", err)
				return
			}

			m.mu.Lock()
			defer m.mu.Unlock()

			if res.Term > m.CurrentTerm {
				m.State = Follower
				m.CurrentTerm = res.Term
				m.resetElectionTimeout()
				return
			}

			if res.Granted == true {
				fmt.Printf("vote granted from id: %v \n", id)
				m.Votes += 1
				if m.Votes*2 >= m.proxy.ClientsCount() && m.State != Leader {
					m.startLeader()
					return
				}
			}
		}(id, cl)
	}
	fmt.Println("electSelf")
}

func (m *module) resetElectionTimeout() {
	t := time.Now().Add(time.Duration(rand.Intn(600))*time.Millisecond + 3*time.Second)
	m.ElectionTimeout = &t
}

func (m *module) RunElectionTimer() {
	for {
		<-m.tick.C
		m.mu.Lock()

		if m.State != Candidate && m.State != Follower {
			fmt.Println("no need to election timer")
			m.mu.Unlock()
			continue
		}

		if m.ElectionTimeout.Before(time.Now()) {
			m.electSelf()
			m.mu.Unlock()
		} else {
			fmt.Printf("to election timeout %v", m.ElectionTimeout.Sub(time.Now()))
			fmt.Println()
			m.mu.Unlock()
		}
	}
}

func (m *module) startLeader() {
	fmt.Println("start leader")
	m.State = Leader
	m.nextIndex = make(map[int64]int64)
	m.matchIndex = make(map[int64]int64)
	l := m.log.LastLog()
	var lastIndex int64
	if l != nil {
		lastIndex = l.Index
	}
	for id, _ := range m.proxy.GetAllClients() {
		m.nextIndex[id] = lastIndex
		m.matchIndex[id] = 0
	}

	go func(heartBeat time.Duration) {
		m.sendAppendEntries()
		ticker := time.NewTicker(heartBeat)
		defer ticker.Stop()

		for {
			<-ticker.C
			m.mu.Lock()
			if m.State != Leader {
				m.mu.Unlock()
				return
			}
			m.mu.Unlock()
			m.sendAppendEntries()
		}
	}(500 * time.Millisecond)
}

func (m *module) restore() {

}

func (m *module) sendAppendEntries() {
	for id, cl := range m.proxy.GetAllClients() {
		go func(id int64, cl rpc.Client) {
			if id == m.LocalID {
				return
			}
			m.mu.Lock()
			nextIndex := m.nextIndex[id]
			var prevLogIndex int64
			var prevLogTerm int64
			if nextIndex > 0 {
				prevLogIndex = nextIndex - 1
				prevLogTerm = m.log.GetLog(prevLogIndex).Term
			}
			entries := m.log.GetFrom(nextIndex)
			m.mu.Unlock()
			_, err := cl.AppendEntries(dto.AppendEntriesQuery{
				Term:         m.CurrentTerm,
				LeaderID:     m.LocalID,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: m.commitIndex,
			})
			if err != nil {
				fmt.Printf("append entries failed: %v \n", err)
				return
			}
			m.mu.Lock()
			defer m.mu.Unlock()

			// logic with append entries response

		}(id, cl)
	}
}

func New(
	localID int64,
	proxy rpc.Proxy,
	log log.Log,
	store store.Store,
) Module {
	return &module{
		log:     log,
		store:   store,
		proxy:   proxy,
		LocalID: localID,
		State:   Follower,
		mu:      &sync.Mutex{},
		tick:    time.NewTicker(500 * time.Millisecond),
	}
}
