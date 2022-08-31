// Package index contains a trie index used by Converter to lookup Replacements.
package index

import (
	"github.com/chanced/caps/token"
)

// IndexedReplacement is a node in an Index
// created from a Replacement
type IndexedReplacement struct {
	Screaming token.Token
	Camel     token.Token
	Lower     token.Token
}

// IsEmpty reports whether or not ir is empty
func (ir IndexedReplacement) IsEmpty() bool {
	return ir.Screaming.IsEmpty()
}

func (ir IndexedReplacement) HasValue() bool {
	return !ir.IsEmpty()
}

// Index is a trie index used by Converter to lookup Replacements.
type Index struct {
	value          IndexedReplacement
	nodes          map[rune]*Index
	lastMatch      IndexedReplacement
	partialMatches []token.Token
	caser          token.Caser
}

func (idx *Index) Clone() Index {
	return Index{
		value:          idx.value,
		nodes:          idx.nodes,
		lastMatch:      idx.lastMatch,
		partialMatches: idx.partialMatches,
		caser:          idx.caser,
	}
}

// NewIndex creates a new Index of Replacements,
// internally represented as Trie
//
// reverseIndexes indicates whether or not
func New(caser token.Caser) *Index {
	return &Index{
		nodes: make(map[rune]*Index),
		caser: token.CaserOrDefault(caser),
	}
}

func (idx Index) Value() IndexedReplacement {
	return idx.value
}

func (idx Index) Contains(tok token.Token) bool {
	_, ok := idx.Get(tok)
	return ok
}

func (idx Index) HasNode(tok token.Token) bool {
	if tok.Len() == 0 {
		return false
	}
	node := &idx
	for _, r := range tok.Lower() {
		if n, ok := node.nodes[r]; ok {
			node = n
		} else {
			return false
		}
	}
	return true
}

func (idx Index) PartialMatches() []token.Token {
	return idx.partialMatches
}

func (idx Index) HasPartialMatches() bool {
	return len(idx.partialMatches) > 0
}

func (idx Index) HasMatch() bool {
	return !idx.lastMatch.IsEmpty()
}

func (idx Index) LastMatch() IndexedReplacement {
	return idx.lastMatch
}

// HasValue reports true if the current node has a value
func (idx Index) HasValue() bool {
	return !idx.value.IsEmpty()
}

// Match searches the index for the given token, returning an Index node.
//
// If t is empty, the root node is returned.
//
// # If the Index does not contain the node, an empty Index is returned
func (idx Index) Match(t token.Token) (Index, bool) {
	var ok bool
	if t.IsEmpty() {
		return idx, false
	}

	next := &idx
	for _, r := range t.Lower() {
		if next, ok = next.nodes[r]; !ok || next == nil {
			return Index{
				partialMatches: idx.partialMatches,
				lastMatch:      idx.lastMatch,
				caser:          idx.caser,
			}, false
		}
		idx = Index{
			nodes:          next.nodes,
			value:          next.value,
			lastMatch:      idx.lastMatch,
			partialMatches: idx.partialMatches,
			caser:          idx.caser,
		}
		if next.HasValue() {
			idx.lastMatch = next.value
			idx.partialMatches = nil
		} else {
			idx.partialMatches = append(idx.partialMatches, token.FromRune(idx.caser, r))
		}
	}
	return idx, true
}

// Get searches the index for the t, returning the IndexedReplacement and true
// if found.
//
// To GetForward a reversed value, use GetReverse.
func (idx *Index) Get(t token.Token) (IndexedReplacement, bool) {
	if t.Len() == 0 {
		return idx.value, idx.value.HasValue()
	}
	node := idx
	var ok bool
	for _, r := range t.Lower() {
		if node, ok = node.nodes[r]; !ok {
			return IndexedReplacement{}, false
		}
	}
	return node.value, node.value.HasValue()
}

func (idx *Index) Nodes() []Index {
	nodes := make([]Index, 0, len(idx.nodes))
	nodes = append(nodes, *idx)
	for _, node := range idx.nodes {
		nodes = append(nodes, node.Nodes()...)
	}
	return nodes
}

func (idx *Index) Values() []IndexedReplacement {
	nodes := idx.Nodes()
	values := make([]IndexedReplacement, 0, len(nodes))
	for _, n := range nodes {
		if n.HasValue() {
			values = append(values, n.value)
		}
	}
	return values
}

// Add inserts r into the Index, indexed by the lowercase variant of r.Camel AND
// r.Screaming.
//
// If the Index does not contains the IndexedReplacement
// repesentation of r, one is created, inserted, and true is returned.
// Otherwise, the previous value is replaced with a new IndexedReplacement.
//
// If idx.IsReversed is true, the IndexedReplacement is inserted into the
// Index with the key in reverse order (e.g. AnExample -> elpmaxena).
func (idx *Index) Add(camel token.Token, screaming token.Token) bool {
	ir := IndexedReplacement{
		Screaming: screaming,
		Camel:     camel,
		Lower:     token.FromString(idx.caser, camel.Lower()),
	}
	var exists bool
	var er IndexedReplacement
	var ok bool
	if er, ok = idx.Get(ir.Screaming); ok {
		exists = true
		idx.Delete(er.Camel)
		idx.Delete(er.Screaming)
	}
	if er, ok = idx.Get(ir.Camel); ok {
		exists = true
		idx.Delete(er.Screaming)
		idx.Delete(er.Camel)
	}

	key := ir.Screaming.LowerRunes()
	node := idx
	for _, r := range key {
		if _, ok = node.nodes[r]; !ok {
			node.nodes[r] = &Index{
				nodes: make(map[rune]*Index),
				caser: idx.caser,
			}
		}
		node = node.nodes[r]
	}
	node.value = ir

	return exists
}

func (idx *Index) Delete(key token.Token) bool {
	tokstr := key.String()
	_ = tokstr
	if key.IsEmpty() {
		return false
	}
	node := idx
	k := key.LowerRunes()
	nodes := make([]*Index, len(k))
	var ok bool
	var i int
	var r rune

	for i, r = range k {
		rstr := string(r)
		_ = rstr
		nodes[i] = node
		if node, ok = node.nodes[r]; !ok || node == nil {
			return false
		}
	}
	node.value = IndexedReplacement{}

	child := node
	for i := len(k) - 1; i >= 0; i-- {
		r = k[i]
		rstr := string(r)
		_ = rstr
		child = node
		node = nodes[i]
		if child.value.IsEmpty() && len(child.nodes) == 0 {
			// safe to delete
			delete(node.nodes, r)
		} else {
			break
		}
	}
	return true
}
