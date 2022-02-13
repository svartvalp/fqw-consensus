package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/svartvalp/fqw-consensus/internal/dto"
)

type Proxy interface {
	GetClient(id int64) (Client, error)
	GetAllClients() map[int64]Client
	ClientsCount() int
}

type Client interface {
	RequestVote(query dto.RequestVoteQuery) (dto.RequestVoteResponse, error)
	AppendEntries(query dto.AppendEntriesQuery) (dto.AppendEntriesResponse, error)
}

type proxy struct {
	clients map[int64]Client
}

func (p *proxy) ClientsCount() int {
	return len(p.clients)
}

func (p *proxy) GetClient(id int64) (Client, error) {
	cl, ok := p.clients[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found client by id %v", id))
	}
	return cl, nil
}

func (p *proxy) GetAllClients() map[int64]Client {
	return p.clients
}

func NewProxy(conns map[int64]string) Proxy {
	clients := make(map[int64]Client)
	for id, host := range conns {
		clients[id] = NewClient(host)
	}
	return &proxy{clients: clients}
}

func NewClient(host string) Client {
	return &client{host: host}
}

type client struct {
	host string
}

func (c *client) RequestVote(query dto.RequestVoteQuery) (dto.RequestVoteResponse, error) {
	b, err := json.Marshal(query)
	if err != nil {
		return dto.RequestVoteResponse{}, err
	}
	httpRes, err := http.Post(c.host+"/request_vote", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return dto.RequestVoteResponse{}, err
	}
	if httpRes.StatusCode != 200 {
		return dto.RequestVoteResponse{}, errors.New(fmt.Sprintf("staus code %v, not 200", httpRes.Status))
	}
	b, err = ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return dto.RequestVoteResponse{}, err
	}
	var response dto.RequestVoteResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return dto.RequestVoteResponse{}, err
	}
	return response, nil
}

func (c *client) AppendEntries(query dto.AppendEntriesQuery) (dto.AppendEntriesResponse, error) {
	b, err := json.Marshal(query)
	if err != nil {
		return dto.AppendEntriesResponse{}, err
	}
	httpRes, err := http.Post(c.host+"/append_entries", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return dto.AppendEntriesResponse{}, err
	}
	if httpRes.StatusCode != 200 {
		return dto.AppendEntriesResponse{}, errors.New(fmt.Sprintf("staus code %v, not 200", httpRes.Status))
	}
	b, err = ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return dto.AppendEntriesResponse{}, err
	}
	var response dto.AppendEntriesResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return dto.AppendEntriesResponse{}, err
	}
	return response, nil
}
