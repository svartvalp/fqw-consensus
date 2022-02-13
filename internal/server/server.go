package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/svartvalp/fqw-consensus/internal/consensus"
	"github.com/svartvalp/fqw-consensus/internal/dto"
)

type Server struct {
	r *gin.Engine
	m consensus.Module
}

func New(m consensus.Module) *Server {
	r := gin.Default()
	return &Server{r: r, m: m}
}

func (s *Server) Start(addr string) error {
	s.r.POST("/request_vote", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		var req dto.RequestVoteQuery
		err = json.Unmarshal(b, &req)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		fmt.Printf("%#v", req)
		res, err := s.m.RequestVote(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}

		c.JSON(200, res)
	})
	s.r.POST("/append_entries", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		var req dto.AppendEntriesQuery
		err = json.Unmarshal(b, &req)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		fmt.Printf("%#v", req)
		res, err := s.m.AppendEntries(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}

		c.JSON(200, res)
	})
	s.r.GET("/state", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		state := s.m.GetState()
		c.JSON(200, state)
	})

	s.m.Start()

	err := s.r.Run(addr)
	if err != nil {
		return err
	}
	return nil
}
