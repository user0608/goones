package hierarchy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user0608/goones/tools/hierarchy"
)

func TestHasCyclicDependency(t *testing.T) {
	tests := []struct {
		name        string
		existing    []testElement
		candidate   testElement
		expectCycle bool
	}{
		{
			name:        "Single node, no cycle",
			existing:    []testElement{},
			candidate:   testElement{Codigo: "A"},
			expectCycle: false,
		},
		{
			name: "Two nodes, no cycle",
			existing: []testElement{
				{Codigo: "A"},
			},
			candidate:   testElement{Codigo: "B", PadreCodigo: strPtr("A")},
			expectCycle: false,
		},
		{
			name: "Simple cycle: A -> B -> A",
			existing: []testElement{
				{Codigo: "A", PadreCodigo: strPtr("B")},
			},
			candidate:   testElement{Codigo: "B", PadreCodigo: strPtr("A")},
			expectCycle: true,
		},
		{
			name: "Deeper cycle: A -> B -> C -> A",
			existing: []testElement{
				{Codigo: "A", PadreCodigo: strPtr("B")},
				{Codigo: "B", PadreCodigo: strPtr("C")},
			},
			candidate:   testElement{Codigo: "C", PadreCodigo: strPtr("A")},
			expectCycle: true,
		},
		{
			name: "No cycle with deeper hierarchy",
			existing: []testElement{
				{Codigo: "X"},
				{Codigo: "A", PadreCodigo: strPtr("X")},
				{Codigo: "B", PadreCodigo: strPtr("A")},
			},
			candidate:   testElement{Codigo: "C", PadreCodigo: strPtr("B")},
			expectCycle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hierarchy.HasCyclicDependency(tt.existing, tt.candidate)
			if tt.expectCycle {
				assert.True(t, result)
			} else {
				assert.False(t, result)
			}
		})
	}
}
