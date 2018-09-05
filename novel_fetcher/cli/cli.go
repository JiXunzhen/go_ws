package main

import (
	"fmt"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/JiXunzhen/go_ws/novel_fetcher/sources"
)

const (
	// Hello ...
	Hello = `
简陋小说阅读器 by gayson
Usage:
	exit
	help

	switch
	search
	list
	select
	pre
	next
	load
	flush
	`
)

var (
	searcher  base.Searcher
	cataloger base.Cataloger
	sectioner base.Sectioner
)

func main() {
	fmt.Println(Hello)
	sources.Init(nil)
	searcher = sources.SearcherMap["笔趣阁"]
	cataloger = base.NilCatalog
	sectioner = base.NilSection

	var handlers = map[string]func() error{
		"exit":   handleExit,
		"help":   handleHelp,
		"switch": handleSwitch,
		"search": handleSearch,
		"list":   handleList,
		"select": handleSelect,
		"pre":    handlePre,
		"next":   handleNext,
		"load":   handleLoad,
		"flush":  handleFlush,
	}

	for {
		fmt.Printf("-- 来源: %s -- 书名: %s -- 章节: %d, %s -- 请输入命令:\n", searcher.Name(), cataloger.GetBookName(), sectioner.GetIndex(), sectioner.GetName())
		var cmd string
		fmt.Scanln(&cmd)
		if handler, ok := handlers[cmd]; ok {
			if err := handler(); err != nil {
				fmt.Println("[ERROR] ", err)
			}
		} else {
			fmt.Printf("[ERROR] 无法识别的命令: %s\n", cmd)
		}
	}
}
