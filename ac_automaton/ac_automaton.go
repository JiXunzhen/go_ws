package ac_automaton

type Automaton struct {
	patterns      []string
	caseSensitive bool

	root *Node
}

func NewAutomaton(patterns []string) *Automaton {
	automaton := &Automaton{
		patterns:      patterns,
		caseSensitive: false,
		root:          NewNode(0, nil, nil),
	}

	automaton.BuildTrie()
	automaton.GenerateFails()
	return automaton
}

func (a *Automaton) SetPatterns(patterns []string) {
	a.patterns = patterns
	a.BuildTrie()
}

func (a *Automaton) AddPattern(pattern string) {
	a.patterns = append(a.patterns, pattern)
	a.insertTrie(pattern)
}

func (a *Automaton) BuildTrie() {
	for _, pattern := range a.patterns {
		a.insertTrie(pattern)
	}
}

func (a *Automaton) insertTrie(pattern string) {
	cur := a.root
	for _, key := range pattern {
		// repeated pattern
		if cur.match == MATCH {
			break
		}

		if n, ok := cur.next[key]; ok {
			cur = n
		} else {
			n = NewNode(key, cur, a.root)
			cur.next[key] = n
			cur = n
		}
	}
	cur.match = MATCH
	cur.fail = a.root
}

func (a *Automaton) GenerateFails() {
	for _, n := range a.root.next {
		a.buildNodeFail(n)
	}
}

func (a *Automaton) buildNodeFail(node *Node) {
	for key, n := range node.next {
		for cur := node; cur != a.root; cur = cur.fail {
			if nextFail, ok := cur.fail.next[key]; ok {
				n.fail = nextFail

				// update match info
				if n.match == NO_MATCH && nextFail.IsMatched() {
					n.match = FAIL_MATCH
				}
				break
			}
		}
		a.buildNodeFail(n)
	}
}

func (a *Automaton) IsMatch(source string) bool {
	cur := a.root

	for _, key := range source {
		for {
			if n, ok := cur.next[key]; ok {
				cur = n
				if n.IsMatched() {
					return true
				}
				break
			}

			if cur == a.root {
				break
			}

			cur = cur.fail
		}
	}
	return false
}

func (a *Automaton) GetMatched(source string) (matched []string) {
	cur := a.root
	for _, key := range source {
		for {
			if n, ok := cur.next[key]; ok {
				cur = n
				if n.IsMatched() {
					matched = append(matched, n.GetMatchedPattern())
				}
				break
			}

			if cur == a.root {
				break
			}

			cur = cur.fail
		}
	}
	return
}
