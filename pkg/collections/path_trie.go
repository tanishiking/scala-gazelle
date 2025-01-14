package collections

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/dghubble/trie"
)

// PathTrie is a copy of https://github.com/dghubble/trie/blob/main/path_trie.go
// with a slightly modified API: this version returns the PathTrie node rather
// than the interface value on the node.

// PathTrie is a trie of string keys and interface{} values. Internal nodes
// have nil values so stored nil values cannot be distinguished and are
// excluded from walks. By default, PathTrie will segment keys by forward
// slashes with PathSegmenter (e.g. "/a/b/c" -> "/a", "/b", "/c"). A custom
// StringSegmenter may be used to customize how strings are segmented into
// nodes. A classic trie might segment keys by rune (i.e. unicode points).
type PathTrie struct {
	segmenter trie.StringSegmenter // key segmenter, must not cause heap allocs
	separator string
	value     interface{}
	children  map[string]*PathTrie
}

// PathTrieConfig for building a path trie with different segmenter
type PathTrieConfig struct {
	Segmenter trie.StringSegmenter
	Separator string
}

// NewPathTrieWithConfig allocates and returns a new *PathTrie with the given *PathTrieConfig
func NewPathTrieWithConfig(config *PathTrieConfig) *PathTrie {
	segmenter := trie.PathSegmenter
	separator := ""
	if config != nil {
		if config.Segmenter != nil {
			segmenter = config.Segmenter
		}
		if config.Separator != "" {
			separator = config.Separator
		}
	}
	return &PathTrie{
		segmenter: segmenter,
		separator: separator,
	}
}

// newPathTrieFromTrie returns new trie while preserving its config
func (trie *PathTrie) newPathTrie() *PathTrie {
	return &PathTrie{
		segmenter: trie.segmenter,
		separator: trie.separator,
	}
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *PathTrie) Get(key string) *PathTrie {
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		node = node.children[part]
		if node == nil {
			return nil
		}
	}
	return node
}

// Put inserts the value into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
func (trie *PathTrie) Put(key string, value interface{}) (*PathTrie, bool) {
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		child := node.children[part]
		if child == nil {
			if node.children == nil {
				node.children = map[string]*PathTrie{}
			}
			child = trie.newPathTrie()
			node.children[part] = child
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.value == nil
	node.value = value
	return node, isNewVal
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the trie.
func (trie *PathTrie) Delete(key string) bool {
	var path []nodeStr // record ancestors to check later
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		path = append(path, nodeStr{part: part, node: node})
		node = node.children[part]
		if node == nil {
			// node does not exist
			return false
		}
	}
	// delete the node value
	node.value = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(path) - 1; i >= 0; i-- {
			parent := path[i].node
			part := path[i].part
			delete(parent.children, part)
			if !parent.isLeaf() {
				// parent has other children, stop
				break
			}
			parent.children = nil
			if parent.value != nil {
				// parent has a value, stop
				break
			}
		}
	}
	return true // node (internal or not) existed and its value was nil'd
}

// Walk iterates over each key/value stored in the trie and calls the given
// walker function with the key and value. If the walker function returns
// an error, the walk is aborted.
// The traversal is depth first with no guaranteed order.
func (trie *PathTrie) Walk(walker trie.WalkFunc) error {
	return trie.walk("", 0, walker)
}

// WalkPath iterates over each key/value in the path in trie from the root to
// the node at the given key, calling the given walker function for each
// key/value. If the walker function returns an error, the walk is aborted.
func (trie *PathTrie) WalkPath(key string, walker trie.WalkFunc) error {
	// Get root value if one exists.
	if trie.value != nil {
		if err := walker("", trie.value); err != nil {
			return err
		}
	}
	for part, i := trie.segmenter(key, 0); ; part, i = trie.segmenter(key, i) {
		if trie = trie.children[part]; trie == nil {
			return nil
		}
		if trie.value != nil {
			var k string
			if i == -1 {
				k = key
			} else {
				k = key[0:i]
			}
			if err := walker(k, trie.value); err != nil {
				return err
			}
		}
		if i == -1 {
			break
		}
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (trie *PathTrie) String() string {
	var buf strings.Builder
	trie.Fprint(&buf, true, "")
	return buf.String()
}

// Fprint prints a tree structure to the given writer.
func (trie *PathTrie) Fprint(w io.Writer, root bool, padding string) {
	if trie == nil {
		return
	}

	keys := make([]string, 0, len(trie.children))
	for key := range trie.children {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	index := 0
	for _, k := range keys {
		v := trie.children[k]
		name := fmt.Sprintf("%s %v", k, v.value)
		fmt.Fprintf(w, "%s%s\n", padding+getBoxPadding(root, getBoxType(index, len(trie.children))), name)
		v.Fprint(w, false, padding+getBoxPadding(root, getBoxTypeExternal(index, len(trie.children))))
		index++
	}
}

// PathTrie node and the part string key of the child the path descends into.
type nodeStr struct {
	node *PathTrie
	part string
}

func (trie *PathTrie) walk(key string, depth int, walker trie.WalkFunc) error {
	if trie.value != nil {
		if err := walker(key, trie.value); err != nil {
			return err
		}
	}
	keys := make([]string, 0, len(trie.children))
	for key := range trie.children {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, part := range keys {
		child := trie.children[part]
		k := key
		if depth != 0 {
			k += trie.separator
		}
		k += part
		if err := child.walk(k, depth+1, walker); err != nil {
			return err
		}
	}
	return nil
}

func (trie *PathTrie) isLeaf() bool {
	return len(trie.children) == 0
}
