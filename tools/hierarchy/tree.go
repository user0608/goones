package hierarchy

import (
	"errors"
	"log/slog"
	"slices"
)

type TreeNode interface {
	GetCode() string
	GetParentCode() *string
}

type NodeInfo[T TreeNode] struct {
	Item      T
	Hierarchy []string
}

type treeNode[T TreeNode] struct {
	Item     T
	Children []*treeNode[T]
}

func (n *treeNode[T]) AddChild(node *treeNode[T]) bool {
	parentCode := node.Item.GetParentCode()
	if n.Item.GetCode() == *parentCode {
		n.Children = append(n.Children, node)
		return true
	}
	for _, child := range n.Children {
		if child.AddChild(node) {
			return true
		}
	}
	return false
}

func (n *treeNode[T]) Flatten(hierarchy []string) []NodeInfo[T] {
	items := []NodeInfo[T]{{Item: n.Item, Hierarchy: hierarchy}}
	newHierarchy := append(slices.Clone(hierarchy), n.Item.GetCode())
	for _, child := range n.Children {
		items = append(items, child.Flatten(newHierarchy)...)
	}
	return items
}

type Tree[T TreeNode] struct {
	Roots []*treeNode[T]
}

func NewTree[T TreeNode]() *Tree[T] {
	return &Tree[T]{Roots: []*treeNode[T]{}}
}

func (t *Tree[T]) AddNode(item T) {
	if item.GetParentCode() == nil {
		t.Roots = append(t.Roots, &treeNode[T]{Item: item, Children: []*treeNode[T]{}})
		return
	}
	node := &treeNode[T]{Item: item, Children: []*treeNode[T]{}}
	for _, root := range t.Roots {
		if root.AddChild(node) {
			return
		}
	}
	slog.Warn(
		"Node insertion failed: possible hierarchy inconsistency",
		"node_code", node.Item.GetCode(),
		"parent_code", *node.Item.GetParentCode(),
	)
}

func (t *Tree[T]) GetFlattenedItems() []NodeInfo[T] {
	var items []NodeInfo[T]
	for _, root := range t.Roots {
		items = append(items, root.Flatten([]string{})...)
	}
	return items
}

func (t *Tree[T]) InsertAll(items []T) error {

	itemMap := make(map[string]T)
	for _, item := range items {
		itemMap[item.GetCode()] = item
	}

	memo := make(map[string]int)

	type itemEntry struct {
		Item  T
		Depth int
	}

	var entries []itemEntry
	for _, item := range items {
		if t.detectCycle(&item, itemMap, make(map[string]bool)) {
			return errors.New("cyclic dependency detected")
		}
		var depth = t.calculateDepth(item.GetCode(), itemMap, memo)
		entries = append(entries, itemEntry{
			Item:  item,
			Depth: depth,
		})
	}
	slices.SortFunc(entries, func(a, b itemEntry) int {
		return a.Depth - b.Depth
	})
	for _, entry := range entries {
		t.AddNode(entry.Item)
	}
	return nil
}

func (t *Tree[T]) detectCycle(node *T, itemMap map[string]T, visited map[string]bool) bool {
	if node == nil {
		return false
	}
	code := (*node).GetCode()
	if visited[code] {
		return true
	}
	visited[code] = true

	if (*node).GetParentCode() == nil {
		return false
	}

	parent, exists := itemMap[*(*node).GetParentCode()]
	if !exists {
		return false
	}
	return t.detectCycle(&parent, itemMap, visited)
}

func (t *Tree[T]) calculateDepth(code string, itemMap map[string]T, memo map[string]int) int {
	const maxDepth = 100

	if val, exists := memo[code]; exists {
		return val
	}

	node, ok := itemMap[code]
	if !ok {
		return -1
	}

	parentCode := node.GetParentCode()
	if parentCode == nil {
		memo[code] = 0
		return 0
	}

	if val, exists := memo[*parentCode]; exists {
		memo[code] = val + 1
		return val + 1
	}

	count := 1
	currentCode := *parentCode
	for count < maxDepth {
		parentNode, exists := itemMap[currentCode]
		if !exists {
			slog.Warn("Resource not found", "code", currentCode)
			break
		}

		nextParent := parentNode.GetParentCode()
		if nextParent == nil {
			break
		}
		currentCode = *nextParent
		count++

		if val, exists := memo[*nextParent]; exists {
			count += val
			break
		}
	}
	memo[code] = count
	return count
}
