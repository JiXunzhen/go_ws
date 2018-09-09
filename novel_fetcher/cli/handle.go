package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/JiXunzhen/go_ws/novel_fetcher/sources"
)

func printSection(s base.Sectioner) error {
	text, err := sectioner.GetText(nil)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println(text)
	fmt.Println()
	return nil
}

func handleHelp() error {
	fmt.Println(Hello)
	return nil
}

func handleExit() error {
	fmt.Println("再您🐎的见！")
	os.Exit(0)
	return nil
}

func handleSwitch() error {
	fmt.Println("可选来源:")
	for idx, name := range sources.SearcherNames {
		fmt.Println("\t", idx, name)
	}
	var val string
	fmt.Scanln(&val)
	index, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if index > len(sources.SearcherNames) {
		return fmt.Errorf("边 界 溢 出: %d", index)
	}
	searcher = sources.SearcherMap[sources.SearcherNames[index]]
	return nil
}

func handleSearch() error {
	fmt.Print("小说名称: ")
	var name string
	fmt.Scanln(&name)
	if name == cataloger.GetBookName() {
		fmt.Println("我 搜 我 自 己")
	}
	catalog, err := searcher.Search(nil, name)
	if err != nil {
		return err
	}
	cataloger = catalog
	fmt.Println("成功")
	return nil
}

func handleList() error {
	step := 0
	for _, section := range cataloger.List() {
		fmt.Printf("%d.%s\t", section.GetIndex(), section.GetName())
		step++
		if step == 2 {
			step = 0
			fmt.Println()
		}
	}
	return nil
}

func handleSelect() error {
	fmt.Print("章节编号: ")
	var val string
	fmt.Scanln(&val)
	index, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if index > cataloger.Count() {
		return fmt.Errorf("边 界 溢 出: %d", index)
	}
	section, err := cataloger.Get(nil, index)
	if err != nil {
		return err
	}
	if section == nil {
		return fmt.Errorf("获取章节有误，请尝试刷新")
	}
	sectioner = section
	return printSection(sectioner)
}

func handlePre() error {
	pre := sectioner.GetPre()
	if pre == nil {
		return fmt.Errorf("没有上一章了")
	}
	sectioner = pre
	return printSection(sectioner)
}

func handleNext() error {
	next := sectioner.GetNext()
	if next == nil {
		return fmt.Errorf("没有下一章了")
	}
	sectioner = next
	return printSection(sectioner)
}

func handleLoad() error {
	fmt.Print("输入欲加载范围(空格分开): ")
	var strStart, strEnd string
	fmt.Scanln(&strStart, &strEnd)

	start, err := strconv.Atoi(strStart)
	if err != nil {
		return err
	}
	end, err := strconv.Atoi(strEnd)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		cataloger.Load(ctx, start, end)
		time.Sleep(5 * time.Second)
		cancel()
	}()
	return nil
}

func handleFlush() error {
	return cataloger.Flush(nil, false)
}
