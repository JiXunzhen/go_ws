package cli

import (
	"fmt"
	"os"

	source "github.com/JiXunzhen/go_ws/novel_fetcher/sources"
)

func handleHelp() error {
	fmt.Println(Hello)
	return nil
}

func handleExit() error {
	os.Exit(0)
	return nil
}

func handleSwitch() error {
	fmt.Println("可选来源:")
	var index int
	for idx, name := range source.SearcherNames {
		fmt.Println(idx, name)
	}
	fmt.Scanln(&index)
	if index < len(source.SearcherNames) {
		searcher = source.SearcherMap[source.SearcherNames[index]]
	}

	return fmt.Errorf("边 界 溢 出: %d", index)
}
