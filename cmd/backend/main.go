package main

import (
	_ "github.com/lib/pq"
	"github.com/zqz/upl/server"
)

func main() {
	s, err := server.Init("./config.json")
	if err != nil {
		s.Log("failed to start server:", err.Error())
		return
	}
	defer s.Close()
	s.Log("zqz backend: booting")

	err = s.Run()
	if err != nil {
		s.Log("error running server:", err.Error())
	}

	s.Log("finished")
}
