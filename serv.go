package main

import (
	_ "github.com/lib/pq"
	"github.com/zqz/upl/server"
)

func main() {
	// os.Mkdir(tmpPath, 0744)
	// os.Mkdir(finalPath, 0744)

	s, err := server.Init("./config.json")
	if err != nil {
		s.Log("failed to start server:", err.Error())
		return
	}
	defer s.Close()

	err = s.Run()
	if err != nil {
		s.Log("error running server:", err.Error())
	}

	s.Log("finished")
}
