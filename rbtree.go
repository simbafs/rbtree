package main

import (
	"fmt"
	"log"
)

type Color int

const (
	Red Color = iota
	Black
)

func (c Color) String() string {
	switch c {
	case Red:
		return "Red"
	case Black:
		return "Black"
	default:
		return "Unknown"
	}
}

type RBTree struct {
	Root *RBNode
	Nil  *RBNode
}

type RBNode struct {
	Value  int
	Left   *RBNode
	Right  *RBNode
	Parent *RBNode
	Color  Color
}

func (node *RBNode) SetLeft(left *RBNode) {
	if node == nil {
		return
	}
	node.Left = left
	if left != nil {
		left.Parent = node
	}
}

func (node *RBNode) SetRight(right *RBNode) {
	if node == nil {
		return
	}
	node.Right = right
	if right != nil {
		right.Parent = node
	}
}

func (node *RBNode) InOrder() string {
	if node == nil {
		return ""
	}
	return fmt.Sprintf("%s %d %s", node.Left.InOrder(), node.Value, node.Right.InOrder())
}

func (node *RBNode) Str() string {
	if node == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%d", node.Value)
}

func (node *RBNode) String() string {
	if node == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Node %d, left: %d, right: %d, parent: %d, color: %s", node.Value, node.Left.Value, node.Right.Value, node.Parent.Value, node.Color)
}

func NewRBTree() *RBTree {
	Nil := &RBNode{
		Value: -1,
		Color: Black,
	}
	return &RBTree{
		Nil:  Nil,
		Root: Nil,
	}
}

func (t *RBTree) Insert(val int) {
	if t.Root == t.Nil {
		t.Root = &RBNode{
			Value:  val,
			Color:  Black,
			Left:   t.Nil,
			Right:  t.Nil,
			Parent: t.Nil,
		}
		return
	}

	curr := t.Root
	for curr != t.Nil {
		if val > curr.Value {
			if curr.Right == t.Nil {
				curr.Right = &RBNode{
					Value:  val,
					Color:  Red,
					Left:   t.Nil,
					Right:  t.Nil,
					Parent: curr,
				}
				t.Fix(curr.Right)
				return
			}
			curr = curr.Right
		} else {
			if curr.Left == t.Nil {
				curr.Left = &RBNode{
					Value:  val,
					Color:  Red,
					Left:   t.Nil,
					Right:  t.Nil,
					Parent: curr,
				}
				t.Fix(curr.Left)
				return
			}
			curr = curr.Left
		}
	}
}

func (t *RBTree) Fix(node *RBNode) {
	for node != t.Nil {
		// node is root
		if node.Parent == t.Nil {
			node.Color = Black
			return
		}

		// parent is root
		if node.Parent.Parent == t.Nil {
			return
		}

		// if parent is black, no need to fix
		if node.Parent.Color == Black {
			return
		}

		// now, node has parent and grandparent
		uncle := node.Parent.Parent.Left
		if uncle == node.Parent {
			uncle = node.Parent.Parent.Right
		}

		if uncle.Color == Red {
			// case 1: uncle is red
			log.Printf("node %d is case 1", node.Value)
			node.Parent.Color = Black
			uncle.Color = Black
			node.Parent.Parent.Color = Red
			node = node.Parent.Parent
		} else {
			// uncle is black
			if node.Parent == node.Parent.Parent.Left && node == node.Parent.Right {
				// case 2: node is LR
				log.Printf("node %d is case 2", node.Value)
				t.RRotate(node.Parent)
				node = node.Left
			} else if node.Parent == node.Parent.Parent.Right && node == node.Parent.Left {
				// case 2 mirror: node is RL
				log.Printf("node %d is case 2 mirror", node.Value)
				t.LRotate(node.Parent)
				node = node.Right
			}

			if node == node.Parent.Left {
				// case 3: node is LL
				log.Printf("node %d is case 3", node.Value)
				node.Parent.Color = Black
				node.Parent.Parent.Color = Red
				t.RRotate(node.Parent.Parent)
			} else if node == node.Parent.Right {
				// case 3 mirror: node is RR
				log.Printf("node %d is case 3 mirror", node.Value)
				node.Parent.Color = Black
				node.Parent.Parent.Color = Red
				t.LRotate(node.Parent.Parent)
			} else {
				log.Printf("unknown case for node %d\n", node.Value)
			}
		}
	}
}

func (t *RBTree) RRotate(node *RBNode) {
	if node == t.Nil || node.Left == t.Nil {
		return
	}

	left := node.Left
	p := node.Parent

	node.SetLeft(left.Right)
	left.SetRight(node)

	if p == t.Nil {
		// node is root
		t.Root = left
	} else if p.Left == node {
		// node is the left child of its parent
		p.SetLeft(left)
	} else {
		// node is the right child of its parent
		p.SetRight(left)
	}
}

func (t *RBTree) LRotate(node *RBNode) {
	if node == t.Nil || node.Right == t.Nil {
		return
	}

	right := node.Right
	p := node.Parent

	node.SetRight(right.Left)
	right.SetLeft(node)

	if p == t.Nil {
		// node is root
		t.Root = right
	} else if p.Left == node {
		// node is the left child of its parent
		p.SetLeft(right)
	} else {
		// node is the right child of its parent
		p.SetRight(right)
	}
}

func (t *RBTree) InOrder() string {
	return t.Root.InOrder()
}

func (t *RBTree) Find(val int) *RBNode {
	curr := t.Root
	for curr != t.Nil && curr.Value != val {
		if val < curr.Value {
			curr = curr.Left
		} else {
			curr = curr.Right
		}
	}
	return curr
}
