package hierarchy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user0608/goones/tools/hierarchy"
)

type testElement struct {
	Codigo      string  `json:"codigo" gorm:"primaryKey"`
	PadreCodigo *string `json:"padre_codigo"`
}

func (l testElement) GetCode() string        { return l.Codigo }
func (l testElement) GetParentCode() *string { return l.PadreCodigo }

func strPtr(s string) *string {
	return &s
}

func TestLotesTree_InsertAll(t *testing.T) {
	tests := []struct {
		name     string
		input    []testElement
		expected []hierarchy.NodeInfo[testElement]
	}{
		{
			name: "Single node with no parent",
			input: []testElement{
				{Codigo: "A"},
			},
			expected: []hierarchy.NodeInfo[testElement]{
				{
					Item:      testElement{Codigo: "A"},
					Hierarchy: []string{},
				},
			},
		},
		{
			name: "Multiple nodes, one hierarchy",
			input: []testElement{
				{Codigo: "A"},
				{Codigo: "B", PadreCodigo: strPtr("A")},
				{Codigo: "C", PadreCodigo: strPtr("B")},
			},
			expected: []hierarchy.NodeInfo[testElement]{
				{
					Item:      testElement{Codigo: "A"},
					Hierarchy: []string{},
				},
				{
					Item:      testElement{Codigo: "B", PadreCodigo: strPtr("A")},
					Hierarchy: []string{"A"},
				},
				{
					Item:      testElement{Codigo: "C", PadreCodigo: strPtr("B")},
					Hierarchy: []string{"A", "B"},
				},
			},
		},
		{
			name: "Two separate trees",
			input: []testElement{
				{Codigo: "A"},
				{Codigo: "B", PadreCodigo: strPtr("A")},
				{Codigo: "X"},
				{Codigo: "Y", PadreCodigo: strPtr("X")},
			},
			expected: []hierarchy.NodeInfo[testElement]{
				{
					Item:      testElement{Codigo: "A"},
					Hierarchy: []string{},
				},
				{
					Item:      testElement{Codigo: "B", PadreCodigo: strPtr("A")},
					Hierarchy: []string{"A"},
				},
				{
					Item:      testElement{Codigo: "X"},
					Hierarchy: []string{},
				},
				{
					Item:      testElement{Codigo: "Y", PadreCodigo: strPtr("X")},
					Hierarchy: []string{"X"},
				},
			},
		},
		{
			name: "Deep hierarchy",
			input: []testElement{
				{Codigo: "Root"},
				{Codigo: "Child1", PadreCodigo: strPtr("Root")},
				{Codigo: "Child2", PadreCodigo: strPtr("Child1")},
				{Codigo: "Child3", PadreCodigo: strPtr("Child2")},
			},
			expected: []hierarchy.NodeInfo[testElement]{
				{
					Item:      testElement{Codigo: "Root"},
					Hierarchy: []string{},
				},
				{
					Item:      testElement{Codigo: "Child1", PadreCodigo: strPtr("Root")},
					Hierarchy: []string{"Root"},
				},
				{
					Item:      testElement{Codigo: "Child2", PadreCodigo: strPtr("Child1")},
					Hierarchy: []string{"Root", "Child1"},
				},
				{
					Item:      testElement{Codigo: "Child3", PadreCodigo: strPtr("Child2")},
					Hierarchy: []string{"Root", "Child1", "Child2"},
				},
			},
		},
		{
			name: "Unsorted input, same result",
			input: []testElement{
				{Codigo: "C", PadreCodigo: strPtr("B")},
				{Codigo: "A"},
				{Codigo: "B", PadreCodigo: strPtr("A")},
				{Codigo: "D", PadreCodigo: strPtr("C")},
			},
			expected: []hierarchy.NodeInfo[testElement]{
				{
					Item:      testElement{Codigo: "A"},
					Hierarchy: []string{},
				},
				{
					Item:      testElement{Codigo: "B", PadreCodigo: strPtr("A")},
					Hierarchy: []string{"A"},
				},
				{
					Item:      testElement{Codigo: "C", PadreCodigo: strPtr("B")},
					Hierarchy: []string{"A", "B"},
				},
				{
					Item:      testElement{Codigo: "D", PadreCodigo: strPtr("C")},
					Hierarchy: []string{"A", "B", "C"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := hierarchy.NewTree[testElement]()
			if err := tree.InsertAll(tt.input); err != nil {
				t.Error(err)
			}
			got := tree.GetFlattenedItems()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestLotesTree_CyclicDependency(t *testing.T) {
	input := []testElement{
		{Codigo: "A", PadreCodigo: strPtr("B")},
		{Codigo: "B", PadreCodigo: strPtr("A")},
	}

	tree := hierarchy.NewTree[testElement]()
	err := tree.InsertAll(input)

	assert.Error(t, err, "expected an error due to cyclic dependency")
	assert.Equal(t, "cyclic dependency detected", err.Error())

	got := tree.GetFlattenedItems()
	assert.Empty(t, got, "no items should have been inserted due to the cycle")
}
