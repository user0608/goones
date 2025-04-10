package hierarchy

import "log/slog"

func HasCyclicDependency[T TreeNode](existing []T, candidate T) bool {
	m := make(map[string]*string)
	for _, item := range existing {
		m[item.GetCode()] = item.GetParentCode()
	}
	m[candidate.GetCode()] = candidate.GetParentCode()

	visited := make(map[string]bool)
	return detectCycle(candidate.GetCode(), m, visited)
}

func detectCycle(code string, parents map[string]*string, visited map[string]bool) bool {
	if visited[code] {
		return true
	}
	visited[code] = true
	parent := parents[code]
	if parent == nil {
		return false
	}
	return detectCycle(*parent, parents, visited)
}

func GetParentCodes[T TreeNode](nodeCode string, itemMap map[string]T) []string {
	node, exists := itemMap[nodeCode]
	if !exists {
		return nil
	}

	parentCode := node.GetParentCode()
	if parentCode == nil {
		return nil
	}

	const maxDepth = 100
	var parentCodes []string
	currentCode := *parentCode

	var visited = map[string]struct{}{nodeCode: {}}
	for range maxDepth {
		if _, ok := visited[currentCode]; ok {
			break
		}
		visited[currentCode] = struct{}{}
		parentCodes = append(parentCodes, currentCode)

		parentNode, exists := itemMap[currentCode]
		if !exists {
			slog.Warn("Parent not found in map", "code", currentCode)
			break
		}

		nextParent := parentNode.GetParentCode()
		if nextParent == nil {
			break
		}
		currentCode = *nextParent
	}

	return parentCodes
}

func GetChildrenCodes[T TreeNode](nodeCode string, itemMap map[string]T) []string {
	var children []string
	visited := map[string]struct{}{nodeCode: {}}
	queue := []string{nodeCode}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for code, node := range itemMap {
			parentCode := node.GetParentCode()
			if parentCode != nil && *parentCode == current {
				if _, seen := visited[code]; !seen {
					visited[code] = struct{}{}
					children = append(children, code)
					queue = append(queue, code) // Añade a la cola para seguir explorando
				}
			}
		}
	}

	return children
}
