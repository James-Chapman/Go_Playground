package tree

import (
	"fmt"
)

type Node struct {
	value int
	left *Node
	right *Node
}

type BinaryTree struct {
	root *Node
}

func (t *BinaryTree) Insert(value int) {
	if t.root == nil {
		fmt.Println(t)
		t.root = &Node{value, nil, nil}
		fmt.Println(t)
	} else {
		t.root.Insert(value)
	}
}

func (n *Node) Insert(value int) {
	if value == n.value {
		return
	} else if value < n.value {
		if n.left == nil {
			n.left = &Node{value, nil, nil}
		} else {
			n.left.Insert(value)
		}
	} else if value > n.value {
		if n.right == nil {
			n.right = &Node{value, nil, nil}
		} else {
			n.right.Insert(value)
		}
	}
}

func (n *Node) Print() {
	fmt.Printf("%d\n", n.value)
	if n.left != nil {
		n.left.Print()
	}
	if n.right != nil {
		n.right.Print()
	}
}

