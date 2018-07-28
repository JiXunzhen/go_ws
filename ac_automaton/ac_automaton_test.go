package ac_automaton

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	patterns := []string{"hjab", "hjac", "abc", "bcg", "ja", "saf"}

	automaton := NewAutomaton(patterns)
	source := "a1sabcghjabcf"
	fmt.Println(automaton.GetMatched(source))
}
