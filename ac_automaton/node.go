package ac_automaton

type MatchType int

const (
	NO_MATCH MatchType = iota
	MATCH
	FAIL_MATCH
)

type Node struct {
	key   rune
	match MatchType
	next  map[rune]*Node
	fail  *Node
	pre   *Node
}

func NewNode(key rune, pre *Node, fail *Node) *Node {
	return &Node{
		key:   key,
		match: NO_MATCH,
		next:  map[rune]*Node{},
		fail:  fail,
		pre:   pre,
	}
}

func (n *Node) IsMatched() bool {
	return n.match == MATCH || n.match == FAIL_MATCH
}

func (n *Node) GetMatchedPattern() string {
	if n.match == MATCH {
		return n.GetPattern()
	} else if n.match == FAIL_MATCH {
		return n.fail.GetMatchedPattern()
	}
	return ""
}

func (n *Node) GetPattern() string {
	pattern := make([]rune, 0, 10)
	pattern = append(pattern, n.key)
	for pre := n.pre; pre != nil; pre = pre.pre {
		pattern = append(pattern, pre.key)
	}

	for i, j := 0, len(pattern)-1; i < j; i, j = i+1, j-1 {
		pattern[i], pattern[j] = pattern[j], pattern[i]
	}

	return string(pattern)
}
