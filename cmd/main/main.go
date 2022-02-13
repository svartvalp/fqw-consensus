package main

import (
	"os"
	"strconv"

	"github.com/svartvalp/fqw-consensus/internal/consensus"
	"github.com/svartvalp/fqw-consensus/internal/log"
	"github.com/svartvalp/fqw-consensus/internal/rpc"
	"github.com/svartvalp/fqw-consensus/internal/server"
	"github.com/svartvalp/fqw-consensus/internal/store"
)

func main() {
	args := os.Args
	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		panic(err)
	}
	host := args[2]

	proxy := rpc.NewProxy(map[int64]string{
		1: "http://localhost:8081",
		2: "http://localhost:8082",
		3: "http://localhost:8083",
	})
	logs := log.New()
	st := store.New()

	module := consensus.New(
		id,
		proxy,
		logs,
		st,
	)
	srv := server.New(module)
	err = srv.Start(host)
	if err != nil {
		panic(err)
	}
}
