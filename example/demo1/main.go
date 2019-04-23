package main

import (
	"fmt"
	"github.com/cn-ygf/yin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	fmt.Println("hello world")
	r := yin.Default()
	r.GET("/test", func(c yin.Context) {
		c.HTML(200, "<h1>test</h1>")
	})
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Printf("demo get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			r.Close()
			log.Printf("demo exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
