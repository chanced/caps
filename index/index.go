/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Chance Dinkins <chanceusc@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, Subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or Substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// Package index contains a trie index used by Converter to lookup Replacements.
package index

import (
	"strings"

	"github.com/chanced/caps/token"
)

// IndexedReplacement is a node in an Index
// created from a Replacement
type IndexedReplacement struct {
	Screaming string
	Camel     string
	Lower     string
}

// IsEmpty reports whether or not ir is empty
func (ir IndexedReplacement) IsEmpty() bool {
	return len(ir.Screaming) == 0
}

func (ir IndexedReplacement) HasValue() bool {
	return !ir.IsEmpty()
}

// Index is a trie index used by Converter to lookup Replacements.
type Index struct {
	value          IndexedReplacement
	nodes          map[rune]*Index
	lastMatch      IndexedReplacement
	partialMatches strings.Builder
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

func (idx Index) Contains(s string) bool {
	_, ok := idx.Get(s)
	return ok
}

// func (idx Index) HasNode(s string) bool {
// 	if len(s) == 0 {
// 		return false
// 	}
// 	node := &idx
// 	for _, r := range token.ToLower(idx.caser, s) {
// 		if n, ok := node.nodes[r]; ok {
// 			node = n
// 		} else {
// 			return false
// 		}
// 	}
// 	return true
// }

func (idx Index) PartialMatches() string {
	return idx.partialMatches.String()
}

func (idx Index) HasPartialMatches() bool {
	return idx.partialMatches.Len() > 0
}

func (idx Index) HasMatched() bool {
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
func (idx Index) Match(s string) (Index, bool) {
	var ok bool
	if len(s) == 0 {
		return idx, false
	}

	next := &idx
	for _, r := range token.ToLower(idx.caser, s) {
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
			idx.partialMatches.Reset()
		} else {
			idx.partialMatches.WriteRune(r)
		}
	}
	return idx, true
}

// Get searches the index for the t, returning the IndexedReplacement and true
// if found.
//
// To GetForward a reversed value, use GetReverse.
func (idx *Index) Get(s string) (IndexedReplacement, bool) {
	if len(s) == 0 {
		return idx.value, idx.value.HasValue()
	}
	node := idx
	var ok bool
	for _, r := range token.ToLower(idx.caser, s) {
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
func (idx *Index) Add(camel string, screaming string) bool {
	ir := IndexedReplacement{
		Screaming: screaming,
		Camel:     camel,
		Lower:     token.ToLower(idx.caser, camel),
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
	node := idx
	for _, r := range ir.Lower {
		if _, ok = node.nodes[r]; !ok {
			node.nodes[r] = &Index{
				nodes: make(map[rune]*Index),
				caser: idx.caser,
			}
		}
		node = node.nodes[r]
	}
	skey := token.ToLower(idx.caser, ir.Screaming)
	if ir.Lower != skey {
		for _, r := range skey {
			if _, ok = node.nodes[r]; !ok {
				node.nodes[r] = &Index{
					nodes: make(map[rune]*Index),
					caser: idx.caser,
				}
			}
			node = node.nodes[r]
		}
	}

	node.value = ir

	return exists
}

func (idx *Index) Delete(key string) bool {
	node := idx

	nodes := make([]*Index, len(key))
	var ok bool
	var i int
	var r rune

	for i, r = range token.ToLower(idx.caser, key) {
		nodes[i] = node
		if node, ok = node.nodes[r]; !ok || node == nil {
			return false
		}
	}
	node.value = IndexedReplacement{}

	child := node
	k := []rune(key)
	for i := len(key) - 1; i >= 0; i-- {
		r = k[i]
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
