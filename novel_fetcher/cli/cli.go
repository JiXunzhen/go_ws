package cli

import (
	"fmt"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
)

const (
	// Hello ...
	Hello = `
简陋小说阅读器 by gayson
Usage:
	switch
	search
	list
	select
	pre
	next
	load
	exit
	help
	`
)

var (
	searcher  base.Searcher
	cataloger base.Searcher
	sectioner base.Sectioner
)

var handlers = map[string]func() error{}

func main() {
	fmt.Println(Hello)
	for {
		var cmd string
		fmt.Println("----------- Input -----------")
		fmt.Scanln(cmd)
		if handler, ok := handlers[cmd]; ok {
			if err := handler(); err != nil {
				fmt.Println("[ERROR] ", err)
			}
		} else {
			fmt.Printf("[ERROR] 无法识别的命令: %s\n", cmd)
		}
	}
}
