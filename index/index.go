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

type indexedValue struct {
	ForwardPath  IndexedReplacement
	ReversedPath IndexedReplacement
}

func (iv indexedValue) isEmpty() bool {
	return iv.ForwardPath.IsEmpty() && iv.ReversedPath.IsEmpty()
}

// IsEmpty reports whether or not ir is empty
func (ir IndexedReplacement) IsEmpty() bool {
	return ir.Screaming.IsEmpty()
}

func (ir IndexedReplacement) HasValue() bool {
	return !ir.IsEmpty()
}

// Index is a double trie (forward and backward indexed) of token.Token.
type Index struct {
	value          indexedValue
	nodes          map[rune]*Index
	lastMatch      IndexedReplacement
	partialMatches []token.Token
}

// NewIndex creates a new Index of Replacements,
// internally represented as Trie
//
// reverseIndexes indicates whether or not
func New() *Index {
	return &Index{
		nodes: make(map[rune]*Index),
	}
}

func (idx Index) ForwardValue() IndexedReplacement {
	return idx.value.ForwardPath
}

func (idx Index) FowardIsEmpty() bool {
	return idx.value.ForwardPath.IsEmpty() && len(idx.nodes) == 0
}

func (idx Index) ContainsForward(tok token.Token) bool {
	_, ok := idx.GetForward(tok)
	return ok
}

func (idx Index) ReverseValue() IndexedReplacement {
	return idx.value.ReversedPath
}

func (idx Index) ReverseIsEmpty() bool {
	return idx.value.ReversedPath.IsEmpty() && len(idx.nodes) == 0
}

func (idx Index) ContainsReverse(tok token.Token) bool {
	_, ok := idx.GetReverse(tok)
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

func (idx Index) HasReverseValue() bool {
	return idx.value.ReversedPath.HasValue()
}

func (idx Index) HasForwardValue() bool {
	return idx.value.ForwardPath.HasValue()
}

// MatchForward searches the index for the given token, returning an Index node.
//
// If t is empty, the root node is returned.
//
// # If the Index does not contain the node, an empty Index is returned
//
// Note: Use MatchReverse if seeking a reversed value
func (idx Index) MatchForward(t token.Token) (Index, bool) {
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
			}, false
		}
		idx = Index{
			nodes:          next.nodes,
			value:          next.value,
			lastMatch:      idx.lastMatch,
			partialMatches: idx.partialMatches,
		}
		if next.HasForwardValue() {
			idx.lastMatch = next.value.ForwardPath
			idx.partialMatches = nil
		} else {
			idx.partialMatches = append(idx.partialMatches, token.FromRune(r))
		}
	}
	return idx, true
}

// MatchReverse attempts to find the Index at the reversed path of t, returning
// an Index containing partial matches and the last match found, if it exists.
//
// Note: the value itself is not reversed, but the path is. For example, a
// search for "nsoj" would return an Index with the LastReverseMatch of
// "JSON"/"Json" (assuming it exists)
func (idx Index) MatchReverse(t token.Token) (Index, bool) {
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
			}, false
		}
		idx = Index{
			nodes:          next.nodes,
			value:          next.value,
			lastMatch:      idx.lastMatch,
			partialMatches: idx.partialMatches,
		}
		if next.HasReverseValue() {
			idx.lastMatch = next.value.ReversedPath
			idx.partialMatches = nil
		} else {
			idx.partialMatches = append([]token.Token{token.FromRune(r)}, idx.partialMatches...)
		}
	}
	return idx, true
}

// GetReverse seeks the value at the reversed path of t, returning the
// IndexedReplacement and true if the value is found.
//
// Note: the value itself is not reversed, but the path is. For example, a
// search for "nsoj" would return an IndexedReplacement with the value
// "JSON"/"Json" (assuming it exists)
func (idx *Index) GetReverse(t token.Token) (IndexedReplacement, bool) {
	if t.Len() == 0 {
		return idx.value.ReversedPath, idx.value.ReversedPath.HasValue()
	}
	node := idx
	var ok bool
	for _, r := range t.Lower() {
		if node, ok = node.nodes[r]; !ok {
			return IndexedReplacement{}, false
		}
	}
	return node.value.ReversedPath, node.value.ReversedPath.HasValue()
}

// GetForwrad searches the index for the t, returning the IndexedReplacement and true
// if found.
//
// To GetForward a reversed value, use GetReverse.
func (idx *Index) GetForward(t token.Token) (IndexedReplacement, bool) {
	if t.Len() == 0 {
		return idx.value.ForwardPath, idx.value.ForwardPath.HasValue()
	}
	node := idx
	var ok bool
	for _, r := range t.Lower() {
		if node, ok = node.nodes[r]; !ok {
			return IndexedReplacement{}, false
		}
	}
	return node.value.ForwardPath, node.value.ForwardPath.HasValue()
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
		if n.HasForwardValue() {
			values = append(values, n.value.ForwardPath)
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
		Lower:     token.FromString(camel.Lower()),
	}
	var exists bool
	var er IndexedReplacement
	var ok bool
	if er, ok = idx.GetForward(ir.Screaming); ok {
		exists = true
		idx.Delete(er.Camel)
		idx.Delete(er.Screaming)
	}
	if er, ok = idx.GetForward(ir.Camel); ok {
		exists = true
		idx.Delete(er.Screaming)
		idx.Delete(er.Camel)
	}
	if er, ok = idx.GetReverse(ir.Screaming.Reverse()); ok {
		exists = true
		idx.Delete(er.Camel)
		idx.Delete(er.Screaming)
	}
	if er, ok = idx.GetReverse(ir.Camel.Reverse()); ok {
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
			}
		}
		node = node.nodes[r]
	}
	node.value.ForwardPath = ir

	key = ir.Screaming.Reverse().LowerRunes()
	node = idx
	for _, r := range key {
		if _, ok = node.nodes[r]; !ok {
			node.nodes[r] = &Index{
				nodes: make(map[rune]*Index),
			}
		}
		node = node.nodes[r]
	}
	node.value.ReversedPath = ir

	return exists
}

func (idx *Index) deleteForward(k []rune) bool {
	node := idx
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
	node.value.ForwardPath = IndexedReplacement{}

	child := node
	for i := len(k) - 1; i >= 0; i-- {
		r = k[i]
		rstr := string(r)
		_ = rstr
		child = node
		node = nodes[i]
		if child.value.isEmpty() && len(child.nodes) == 0 {
			// safe to delete
			delete(node.nodes, r)
		} else {
			break
		}
	}
	return true
}

func (idx *Index) deleteReverse(k []rune) bool {
	node := idx
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
	node.value.ReversedPath = IndexedReplacement{}

	child := node
	for i := len(k) - 1; i >= 0; i-- {
		r = k[i]
		rstr := string(r)
		_ = rstr
		child = node
		node = nodes[i]
		if child.value.isEmpty() && len(child.nodes) == 0 {
			// safe to delete
			delete(node.nodes, r)
		} else {
			break
		}
	}
	return true
}

func (idx *Index) Delete(key token.Token) bool {
	tokstr := key.String()
	_ = tokstr
	if key.IsEmpty() {
		return false
	}
	rev := idx.deleteReverse(key.Reverse().LowerRunes())
	fwd := idx.deleteForward(key.LowerRunes())
	return rev || fwd
}
