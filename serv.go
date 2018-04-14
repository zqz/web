package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/zqz/upl/server"
)

func main() {
	// os.Mkdir(tmpPath, 0744)
	// os.Mkdir(finalPath, 0744)

	l := log.New(os.Stdout, "", log.LstdFlags)
	s, err := server.Init("./config.json", l)
	if err != nil {
		l.Println("failed to start server:", err.Error())
		return
	}
	defer s.Close()

	err = s.Run()
	if err != nil {
		l.Println("error running server:", err.Error())
	}

	l.Println("finished")
}
