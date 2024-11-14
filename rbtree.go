package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
}

type RBNode struct {
	Value  int
	Left   *RBNode
	Right  *RBNode
	Parent *RBNode
	Color  Color
}

var Nil = &RBNode{
	Value: -1,
	Color: Black,
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
	if node == nil || node == Nil {
		return ""
	}
	return fmt.Sprintf("%s %d %s", node.Left.InOrder(), node.Value, node.Right.InOrder())
}

func (node *RBNode) Predecessor() *RBNode {
	curr := node.Left

	for curr.Right != Nil {
		curr = curr.Right
	}

	return curr
}

func (node *RBNode) Successor() *RBNode {
	curr := node.Right

	for curr.Left != Nil {
		curr = curr.Left
	}

	return curr
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
	if node == Nil {
		return "<Nil>"
	}
	return fmt.Sprintf("Node %d, left: %d, right: %d, parent: %d, color: %s", node.Value, node.Left.Value, node.Right.Value, node.Parent.Value, node.Color)
}

func NewRBTree() *RBTree {
	return &RBTree{
		Root: Nil,
	}
}

func (t *RBTree) Insert(val int) {
	if t.Root == Nil {
		t.Root = &RBNode{
			Value:  val,
			Color:  Black,
			Left:   Nil,
			Right:  Nil,
			Parent: Nil,
		}
		return
	}

	curr := t.Root
	for curr != Nil {
		if val > curr.Value {
			if curr.Right == Nil {
				curr.Right = &RBNode{
					Value:  val,
					Color:  Red,
					Left:   Nil,
					Right:  Nil,
					Parent: curr,
				}
				t.Fix(curr.Right)
				return
			}
			curr = curr.Right
		} else {
			if curr.Left == Nil {
				curr.Left = &RBNode{
					Value:  val,
					Color:  Red,
					Left:   Nil,
					Right:  Nil,
					Parent: curr,
				}
				t.Fix(curr.Left)
				return
			}
			curr = curr.Left
		}
	}
}

func (t *RBTree) Delete(node *RBNode) {
	if node == nil || node == Nil {
		return
	}

	log.Printf("delete(%d)", node.Value)

	// get child
	child := node.Left
	if child == Nil {
		child = node.Right
	}

	oneChild := node.Left == Nil || node.Right == Nil

	if node.Parent == Nil {
		// node is root
		if oneChild {
			log.Printf("node is root, it has at most one child")
			child.Parent = Nil
			t.Root = child
		} else {
			log.Printf("node is root, it has two children")
			suc := node.Successor()

			suc.SetLeft(node.Left)

			suc.Parent = Nil
			t.Root = suc
		}
	} else {
		if oneChild {
			log.Printf("node is not root, it has at most one child")
			if node.Parent.Left == node {
				log.Printf("node is its parent's left child")

				node.Parent.SetLeft(child)
			} else {
				log.Printf("node is its parent's right child")

				node.Parent.SetRight(child)
			}
		} else {
			log.Printf("node is not root, it has two children")
			suc := node.Successor()
			log.Println("suc", suc)

			// TODO: Delete node by moving not replacing
			node.Value = suc.Value
			node.Color = suc.Color

			t.Delete(suc)

			// log.Println(1)
			// suc.SetRight(node.Right)
			// log.Println(2)
			// suc.SetLeft(node.Left)
			// log.Println(3)
			//
			// if node.Parent.Left == node {
			// 	log.Println(4)
			// 	node.Parent.SetLeft(suc)
			// 	log.Println(5)
			// } else {
			// 	log.Println(6)
			// 	node.Parent.SetRight(suc)
			// 	log.Println(7)
			// }
			// log.Println(8)
		}
	}
}

func (t *RBTree) Fix(node *RBNode) {
	for node != Nil {
		// node is root
		if node.Parent == Nil {
			node.Color = Black
			return
		}

		// parent is root
		if node.Parent.Parent == Nil {
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
				t.LRotate(node.Parent)
				node = node.Left
			} else if node.Parent == node.Parent.Parent.Right && node == node.Parent.Left {
				// case 2 mirror: node is RL
				log.Printf("node %d is case 2 mirror", node.Value)
				t.RRotate(node.Parent)
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

// func (t *RBTree) RRotate(node **RBNode) {
// 	if *node == Nil || (*node).Left == Nil {
// 		return
// 	}
//
// 	left := (*node).Left
// 	p := (*node).Parent
//
// 	(*node).SetLeft(left.Right)
// 	left.SetRight(*node)
//
// 	left.Parent = p
//
// 	*node = left
// }

func (t *RBTree) RRotate(node *RBNode) {
	if node == Nil || node.Left == Nil {
		return
	}

	left := node.Left
	p := node.Parent

	node.SetLeft(left.Right)
	left.SetRight(node)

	left.Parent = p

	if p == Nil {
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
	if node == Nil || node.Right == Nil {
		return
	}

	right := node.Right
	p := node.Parent

	node.SetRight(right.Left)
	right.SetLeft(node)

	right.Parent = p

	if p == Nil {
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
	for curr != Nil && curr.Value != val {
		if val < curr.Value {
			curr = curr.Left
		} else {
			curr = curr.Right
		}
	}
	return curr
}

func (t *RBTree) Query(query string) (*RBNode, error) {
	parts := strings.Split(query, ".")

	curr := t.Root

	for _, part := range parts {
		if curr == nil {
			return nil, fmt.Errorf("Node is nil at path: %s", strings.Join(parts, "."))
		}
		switch part {
		case "left":
			curr = curr.Left
		case "right":
			curr = curr.Right
		case "parent":
			curr = curr.Parent
		case "root", "":
			curr = t.Root
		default:
			if val, err := strconv.Atoi(part); err == nil {
				curr = t.Find(val)
			} else {
				return nil, fmt.Errorf("Unknown path: %s", part)
			}
		}

	}

	return curr, nil
}
