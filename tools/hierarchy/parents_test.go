package hierarchy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetParentCodes(t *testing.T) {
	tests := []struct {
		name     string
		itemMap  map[string]TestNode
		nodeCode string
		expected []string
	}{
		{
			name:     "Nodo raíz (sin padres)",
			itemMap:  map[string]TestNode{"A": {code: "A", parentCode: nil}},
			nodeCode: "A",
			expected: nil,
		},
		{
			name: "Nodo con un solo padre",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
			},
			nodeCode: "B",
			expected: []string{"A"},
		},
		{
			name: "Nodo con múltiples niveles de padres",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
			},
			nodeCode: "D",
			expected: []string{"C", "B", "A"},
		},
		{
			name: "Nodo con referencia a padre no existente",
			itemMap: map[string]TestNode{
				"B": {code: "B", parentCode: stringPtr("X")}, // "X" no está en el mapa
			},
			nodeCode: "B",
			expected: []string{"X"}, // Se incluye el padre inexistente, pero se detiene ahí
		},
		{
			name: "Nodo con referencia cíclica",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: stringPtr("C")},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")}, // Ciclo: A → B → C → A
			},
			nodeCode: "C",
			expected: []string{"B", "A"}, // Se corta antes del ciclo infinito
		},
		{
			name: "Nodo con padre en camino diferente",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"X": {code: "X", parentCode: stringPtr("A")}, // Camino diferente
				"Y": {code: "Y", parentCode: stringPtr("X")},
			},
			nodeCode: "C",
			expected: []string{"B", "A"}, // No debe incluir X o Y
		},
		{
			name: "Nodo con gran profundidad pero sin ciclo",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
				"E": {code: "E", parentCode: stringPtr("D")},
				"F": {code: "F", parentCode: stringPtr("E")},
				"G": {code: "G", parentCode: stringPtr("F")},
				"H": {code: "H", parentCode: stringPtr("G")},
				"I": {code: "I", parentCode: stringPtr("H")},
				"J": {code: "J", parentCode: stringPtr("I")},
				"K": {code: "K", parentCode: stringPtr("J")},
				"L": {code: "L", parentCode: stringPtr("K")},
				"M": {code: "M", parentCode: stringPtr("L")},
				"N": {code: "N", parentCode: stringPtr("M")},
				"O": {code: "O", parentCode: stringPtr("N")},
				"P": {code: "P", parentCode: stringPtr("O")},
				"Q": {code: "Q", parentCode: stringPtr("P")},
				"R": {code: "R", parentCode: stringPtr("Q")},
				"S": {code: "S", parentCode: stringPtr("R")},
				"T": {code: "T", parentCode: stringPtr("S")},
				"U": {code: "U", parentCode: stringPtr("T")},
				"V": {code: "V", parentCode: stringPtr("U")},
				"W": {code: "W", parentCode: stringPtr("V")},
				"X": {code: "X", parentCode: stringPtr("W")},
				"Y": {code: "Y", parentCode: stringPtr("X")},
				"Z": {code: "Z", parentCode: stringPtr("Y")},
			},
			nodeCode: "Z",
			expected: []string{"Y", "X", "W", "V", "U", "T", "S", "R", "Q", "P", "O", "N", "M", "L", "K", "J", "I", "H", "G", "F", "E", "D", "C", "B", "A"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parentCodes := GetParentCodes(tt.nodeCode, tt.itemMap)
			assert.Equal(t, tt.expected, parentCodes)
		})
	}
}
