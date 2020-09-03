package tree

import (
	//"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBinaryTree(t *testing.T) {
	
	var tree = BinaryTree{}
	s := []int{7,9,4,1,0,4,6,8,2,3,6}
	
	for _, v := range s {
		tree.Insert(v)
	}

	tree.root.Print()

	assert.Equal(t, 6, tree.root.left.right.value, "Values should match")
}

