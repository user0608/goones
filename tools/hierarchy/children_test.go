package hierarchy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChildrenCodes(t *testing.T) {
	tests := []struct {
		name     string
		itemMap  map[string]TestNode
		nodeCode string
		expected []string
	}{
		{
			name: "Nodo sin hijos",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
			},
			nodeCode: "A",
			expected: nil,
		},
		{
			name: "Nodo con un solo hijo",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
			},
			nodeCode: "A",
			expected: []string{"B"},
		},
		{
			name: "Nodo con múltiples hijos directos",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("A")},
			},
			nodeCode: "A",
			expected: []string{"B", "C"},
		},
		{
			name: "Nodo con hijos en múltiples niveles",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
				"E": {code: "E", parentCode: stringPtr("D")},
			},
			nodeCode: "A",
			expected: []string{"B", "C", "D", "E"}, // Debe incluir descendientes
		},
		{
			name: "Nodo con ciclo detectado",
			itemMap: map[string]TestNode{
				"A":  {code: "A", parentCode: nil},
				"B":  {code: "B", parentCode: stringPtr("A")},
				"C":  {code: "C", parentCode: stringPtr("B")},
				"A2": {code: "A2", parentCode: stringPtr("C")}, // Ciclo: A → B → C → A2 → C
			},
			nodeCode: "A",
			expected: []string{"B", "C", "A2"}, // Se corta antes del ciclo infinito
		},
		{
			name: "Nodo con múltiples ramas de hijos",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("A")},
				"D": {code: "D", parentCode: stringPtr("B")},
				"E": {code: "E", parentCode: stringPtr("B")},
				"F": {code: "F", parentCode: stringPtr("C")},
				"G": {code: "G", parentCode: stringPtr("E")},
			},
			nodeCode: "A",
			expected: []string{"B", "C", "D", "E", "F", "G"},
		},
		{
			name: "Nodo con múltiples ramas de hijos (2)",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("A")},
				"D": {code: "D", parentCode: stringPtr("B")},
				"E": {code: "E", parentCode: stringPtr("B")},
				"F": {code: "F", parentCode: stringPtr("C")},
				"G": {code: "G", parentCode: stringPtr("E")},
			},
			nodeCode: "B",
			expected: []string{"D", "E", "G"},
		},
		{
			name: "Nodo con hijo no listado en itemMap",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"X": {code: "X", parentCode: stringPtr("Y")}, // "Y" no está en el mapa
			},
			nodeCode: "A",
			expected: []string{"B"},
		},
		{
			name: "Nodo inexistente en el mapa",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
			},
			nodeCode: "Z", // No existe en itemMap
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			childrenCodes := GetChildrenCodes(tt.nodeCode, tt.itemMap)
			assert.ElementsMatch(t, tt.expected, childrenCodes)
		})
	}
}
