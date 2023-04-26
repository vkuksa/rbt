package rbt

import (
	"fmt"

	"github.com/vkuksa/rbt/utils"
	"golang.org/x/exp/constraints"
)

type Color uint8

const (
	black = Color(iota)
	red
)

func (c Color) String() string {
	switch c {
	case red:
		return "RED"
	case black:
		return "BLACK"
	default:
		panic("invalid node color")
	}
}

type Node[K constraints.Ordered, V any] struct {
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
	key    K
	value  V
	color  Color
}

func MakeNode[K constraints.Ordered, V any](k K, v V, c Color) *Node[K, V] {
	return &Node[K, V]{
		key:   k,
		value: v,
		color: c,
	}
}

func (node *Node[K, V]) String() string {
	return fmt.Sprintf("%v(%s)", node.key, node.color)
}

func (node *Node[K, V]) Color() Color {
	if node == nil {
		return black
	}

	return node.color
}

func (node *Node[K, V]) grandparent() *Node[K, V] {
	if node == nil || node.parent == nil {
		return nil
	} else {
		return node.parent.parent
	}
}

func (node *Node[K, V]) inorder(closure func(node *Node[K, V])) {
	if node != nil {
		node.left.inorder(closure)
		closure(node)
		node.right.inorder(closure)
	}
}

func (node *Node[K, V]) preorder(closure func(node *Node[K, V])) {
	if node != nil {
		closure(node)
		node.left.preorder(closure)
		node.right.preorder(closure)
	}
}

func (node *Node[K, V]) minimum() *Node[K, V] {
	if node == nil {
		return nil
	}

	for node.left != nil {
		node = node.left
	}
	return node
}

func (node *Node[K, V]) isBinarySearchTree(minKey, maxKey K) bool {
	if node == nil {
		return true
	}
	if node.key <= minKey || node.key >= maxKey {
		return false
	}
	return node.left.isBinarySearchTree(minKey, node.key) &&
		node.right.isBinarySearchTree(node.key, maxKey)
}

func (node *Node[K, V]) hasSameBlackHeight() bool {
	if node == nil {
		return true
	}

	leftBlackHeight := node.left.getBlackHeight()
	rightBlackHeight := node.right.getBlackHeight()
	if leftBlackHeight != rightBlackHeight {
		return false
	}
	return node.left.hasSameBlackHeight() && node.right.hasSameBlackHeight()
}

func (node *Node[K, V]) getBlackHeight() int {
	if node == nil {
		return 0
	}

	height := 0
	for node != nil {
		if node.color == black {
			height++
		}
		node = node.left
	}
	return height
}

func (node *Node[K, V]) containsConsecutiveRedNodes() bool {
	if node == nil {
		return false
	}

	if node.color == red && ((node.left != nil && node.left.color == red) || (node.right != nil && node.right.color == red)) {
		return true
	}
	return node.left.containsConsecutiveRedNodes() || node.right.containsConsecutiveRedNodes()
}

type RedBlackTree[K constraints.Ordered, V any] struct {
	root *Node[K, V]
	size int
}

// Creates empty instance of a tree
func Make[K constraints.Ordered, V any]() *RedBlackTree[K, V] {
	return &RedBlackTree[K, V]{}
}

// Inserts a value into a tree for a key with a given key
// If key already exists - updates it's value
func (tree *RedBlackTree[K, V]) Insert(k K, v V) {
	n := MakeNode(k, v, red) // New node to be inserted

	if tree.root == nil {
		tree.root = n // Tree is empty
	} else {
		current := tree.root // Node to be compared against during lookup

	loop:
		for {
			switch {
			case k < current.key:
				if current.left == nil {
					current.left = n
					break loop
				} else {
					current = current.left
				}
			case k > current.key:
				if current.right == nil {
					current.right = n
					break loop
				} else {
					current = current.right
				}
			default:
				// key already exists, update the value
				current.value = v
				return
			}
		}

		n.parent = current
	}

	tree.insertFixup(n)
	tree.size++
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Returns second bool parameter that indicates whether value was found in a map
func (tree *RedBlackTree[K, V]) Search(k K) (value V, exists bool) {
	if n := tree.search(k); n == nil {
		exists = false
	} else {
		exists = true
		value = n.value
	}

	return
}

// Removes given key/value pair from a tree
func (tree *RedBlackTree[K, V]) Remove(k K) {
	if n := tree.search(k); n != nil {
		tree.delete(n)
	}

	return
}

// Prints ordered collection of keys stored in a tree
func (tree *RedBlackTree[K, V]) Keys() []K {
	// Gathering copies of keys
	result := make([]K, 0, tree.Size())
	tree.root.inorder(func(n *Node[K, V]) {
		result = append(result, *&n.key)
	})
	return result
}

// Returns a copy of a tree nodes in a preordered traversal
func (tree *RedBlackTree[K, V]) Preorder() []Node[K, V] {
	// Gathering copies of nodes, to not expose the tree for editing from within
	result := make([]Node[K, V], 0, tree.Size())
	tree.root.preorder(func(n *Node[K, V]) {
		result = append(result, *n)
	})
	return result
}

// Returns a number of elements stored in a tree
func (tree *RedBlackTree[K, V]) Size() int {
	return tree.size
}

// A red-black tree satisfies the following properties:
//
//	Red/Black Property: Every node is colored, either red or black.
//	Root Property: The root is black.
//	Leaf Property: Every leaf (NIL) is black.
//	Red Property: If a red node has children then, the children are always black.
//	Depth Property: For each node, any simple path from this node to any of its descendant leaf has the same black-depth (the number of black nodes).

// Note: as this implementation does not use sentinel nodes,
// when tree is empty - it is not considered as valid
func (tree *RedBlackTree[K, V]) IsValidRBTree() bool {
	// Check if root is nil or not black
	if tree.root == nil || tree.root.color != black {
		return false
	}
	// Recursively check if a tree is a valid BST
	if !tree.root.isBinarySearchTree(utils.MinValue[K](), utils.MaxValue[K]()) {
		return false
	}
	// Recursively check if there are any consecutive red nodes
	if tree.root.containsConsecutiveRedNodes() {
		return false
	}
	// Recursively check that all paths from root to leaf nodes have the same black height
	if !tree.root.hasSameBlackHeight() {
		return false
	}
	// If all checks pass, the tree is valid
	return true
}

func (tree *RedBlackTree[K, V]) transplant(x, y *Node[K, V]) {
	if x.parent == nil {
		tree.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	if y != nil {
		y.parent = x.parent
	}
}

func (tree *RedBlackTree[K, V]) leftRotate(n *Node[K, V]) {
	r := n.right
	n.right = r.left
	if r.left != nil {
		r.left.parent = n
	}
	r.parent = n.parent

	tree.transplant(n, r)

	r.left = n
	n.parent = r
}

func (tree *RedBlackTree[K, V]) rightRotate(n *Node[K, V]) {
	l := n.left
	n.left = l.right
	if l.right != nil {
		l.right.parent = n
	}
	l.parent = n.parent

	tree.transplant(n, l)

	l.right = n
	n.parent = l
}

func (tree *RedBlackTree[K, V]) search(k K) *Node[K, V] {
	for x := tree.root; x != nil; {
		switch {
		case k == x.key:
			return x
		case k < x.key:
			x = x.left
		case k > x.key:
			x = x.right
		}
	}

	return nil
}

func (tree *RedBlackTree[K, V]) insertFixup(n *Node[K, V]) {
	for n != tree.root && n.parent.color == red {
		switch n.parent {
		// If our parent is a left child of our grandparent
		case n.grandparent().left:
			if u := n.grandparent().right; u != nil && u.color == red {
				n.parent.color = black
				u.color = black
				n.grandparent().color = red
				n = n.grandparent()
			} else {
				if n == n.parent.right {
					n = n.parent
					tree.leftRotate(n)
				}

				n.parent.color = black
				n.grandparent().color = red
				tree.rightRotate(n.grandparent())
			}
		// If our parent is a right child of our grandparent
		case n.grandparent().right:
			if u := n.grandparent().left; u != nil && u.color == red {
				n.parent.color = black
				u.color = black
				n.grandparent().color = red
				n = n.grandparent()
			} else {
				if n == n.parent.left {
					n = n.parent
					tree.rightRotate(n)
				}

				n.parent.color = black
				n.grandparent().color = red
				tree.leftRotate(n.grandparent())
			}
		default:
			panic("parent has no grandparent")
		}
	}

	tree.root.color = black
}

func (tree *RedBlackTree[K, V]) delete(node *Node[K, V]) {
	var x, y *Node[K, V]

	// If the node to be deleted has at most one child, that node is spliced out.
	// Otherwise, find the node's minimum and splice out that node instead.
	if node.left == nil || node.right == nil {
		y = node
	} else {
		y = node.right.minimum()
	}

	// Determine the value of x
	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	// Set node's replacement
	if x != nil {
		x.parent = y.parent
	}

	if y.parent == nil {
		tree.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	// If the spliced-out node is not the same as the node to be deleted,
	// overwrite the key of the node to be deleted with the key of the spliced-out node
	if y != node {
		node.key = y.key
		node.value = y.value
	}

	// Perform fixup of the colors after deletion
	if y.color == black {
		tree.deleteFixup(x, y.parent)
	}
	tree.size--
}

func (tree *RedBlackTree[K, V]) deleteFixup(n *Node[K, V], parent *Node[K, V]) {
	var sib *Node[K, V]

	for n != tree.root && n.Color() == black {
		if n == parent.left {
			sib = parent.right
			if sib.color == red {
				// Case 1: sibling is red
				sib.color = black
				parent.color = red
				tree.leftRotate(parent)
				sib = parent.right
			}
			if sib.left.Color() == black && sib.right.Color() == black {
				// Case 2: sibling and its children are black
				sib.color = red
				n = parent
				parent = n.parent
			} else {
				if sib.right.Color() == black {
					// Case 3: sibling is black, sibling's left child is red and sibling's right child is black
					sib.left.color = black
					sib.color = red
					tree.rightRotate(sib)
					sib = parent.right
				}
				// Case 4: sibling is black, sibling's right child is red
				sib.color = parent.color
				parent.color = black
				sib.right.color = black
				tree.leftRotate(parent)
				n = tree.root
			}
		} else {
			sib = parent.left
			if sib.Color() == red {
				// Case 1: sibling is red
				sib.color = black
				parent.color = red
				tree.rightRotate(parent)
				sib = parent.left
			}
			if sib.left.Color() == black && sib.right.Color() == black {
				// Case 2: sibling and its children are black
				sib.color = red
				n = parent
				parent = n.parent
			} else {
				if sib.left.Color() == black {
					// Case 3: sibling is black, sibling's right child is red and sibling's left child is black
					sib.right.color = black
					sib.color = red
					tree.leftRotate(sib)
					sib = parent.left
				}
				// Case 4: sibling is black, sibling's left child is red
				sib.color = parent.color
				parent.color = black
				sib.left.color = black
				tree.rightRotate(parent)
				n = tree.root
			}
		}
	}

	if n != nil {
		n.color = black
	}
}
