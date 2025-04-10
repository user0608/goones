package hierarchy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestNode struct {
	code       string
	parentCode *string
}

func (n TestNode) GetCode() string        { return n.code }
func (n TestNode) GetParentCode() *string { return n.parentCode }

func TestCalculateDepth(t *testing.T) {
	tests := []struct {
		name     string
		itemMap  map[string]TestNode
		nodeCode string
		expected int
		memo     map[string]int
	}{
		{
			name:     "Nodo raíz (nil parentCode)",
			itemMap:  map[string]TestNode{},
			expected: -1,
			memo:     make(map[string]int),
		},
		{
			name: "Nodo con un solo nivel de profundidad",
			itemMap: map[string]TestNode{
				"B": {code: "B", parentCode: stringPtr("A")},
			},
			nodeCode: "B",
			expected: 1,
			memo:     make(map[string]int),
		},
		{
			name: "Nodo con múltiples niveles de profundidad",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
			},
			nodeCode: "D",
			expected: 3,
			memo:     make(map[string]int),
		},
		{
			name: "Nodo con referencia cíclica",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: stringPtr("C")},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")}, // Ciclo: A → B → C → A
			},
			nodeCode: "C",
			expected: 100, // Se debe cortar en la primera detección de ciclo
			memo:     make(map[string]int),
		},
		{
			name: "Nodo con un padre que no existe en el mapa",
			itemMap: map[string]TestNode{
				"B": {code: "B", parentCode: stringPtr("X")}, // "X" no está en itemMap
			},
			nodeCode: "B",
			expected: 1,
			memo:     make(map[string]int),
		},
		{
			name: "Nodo con un solo nivel de profundidad con memoria",
			itemMap: map[string]TestNode{
				"B": {code: "B", parentCode: stringPtr("A")},
			},
			nodeCode: "B",
			expected: 1,
			memo:     map[string]int{"B": 1},
		},
		{
			name: "Nodo con múltiples niveles de profundidad con memoria",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
				"X": {code: "X", parentCode: stringPtr("D")},
			},
			nodeCode: "X",
			expected: 4,
			memo:     map[string]int{"D": 3},
		},
		{
			name: "Nodo con múltiples niveles de profundidad con memoria 2",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
				"X": {code: "X", parentCode: stringPtr("D")},
			},
			nodeCode: "X",
			expected: 4,
			memo:     map[string]int{"B": 1},
		},
		{
			name: "Nodo con memoria optimizando cálculo en un árbol profundo",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
				"E": {code: "E", parentCode: stringPtr("D")},
			},
			nodeCode: "E",
			expected: 4,
			memo:     map[string]int{"C": 2}, // Usa el valor memorizado de "C" para reducir cálculos
		},
		{
			name: "Nodo con memoria optimizando cálculo en la mitad del camino",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
				"E": {code: "E", parentCode: stringPtr("D")},
				"F": {code: "F", parentCode: stringPtr("E")},
			},
			nodeCode: "F",
			expected: 5,
			memo:     map[string]int{"D": 3}, // "D" ya tiene su profundidad memorizada
		},
		{
			name: "Nodo con memoria almacenando el resultado tras la ejecución",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
			},
			nodeCode: "D",
			expected: 3,
			memo:     map[string]int{}, // Al inicio está vacío, pero después de ejecutar debería contener "D": 3
		},
		{
			name: "Nodo con memoria previa de otro camino",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"X": {code: "X", parentCode: stringPtr("A")}, // Camino diferente
				"Y": {code: "Y", parentCode: stringPtr("X")},
			},
			nodeCode: "C",
			expected: 2,
			memo:     map[string]int{"X": 1}, // Memoria de un nodo de otro camino, no debe afectar
		},
		{
			name: "Nodo con memoria que cubre parte del cálculo",
			itemMap: map[string]TestNode{
				"A": {code: "A", parentCode: nil},
				"B": {code: "B", parentCode: stringPtr("A")},
				"C": {code: "C", parentCode: stringPtr("B")},
				"D": {code: "D", parentCode: stringPtr("C")},
				"E": {code: "E", parentCode: stringPtr("D")},
				"F": {code: "F", parentCode: stringPtr("E")},
			},
			nodeCode: "F",
			expected: 5,
			memo:     map[string]int{"C": 2, "E": 4}, // Ambos valores deben ser aprovechados
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := &Tree[TestNode]{}

			depth := tree.calculateDepth(tt.nodeCode, tt.itemMap, tt.memo)
			assert.Equal(t, tt.expected, depth)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
