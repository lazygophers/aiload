package main

import (
	"github.com/lazygophers/log"
)

func main() {
	err := Run()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
