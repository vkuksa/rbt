package rbt

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NodeSuite struct {
	suite.Suite
}

func (suite *NodeSuite) TestGrandparent() {
	var n, p, gp *Node[int, int]
	assert.Nil(suite.T(), n.grandparent())
	n = &Node[int, int]{}
	assert.Nil(suite.T(), n.grandparent())
	p = &Node[int, int]{}
	assert.Nil(suite.T(), n.grandparent())

	gp = &Node[int, int]{value: 1}
	p.parent = gp
	n.parent = p
	assert.Equal(suite.T(), gp.value, n.grandparent().value)
}

func (suite *NodeSuite) TestMinimum() {
	var n, m *Node[int, int]
	assert.Nil(suite.T(), n.minimum())

	n = &Node[int, int]{value: 1}
	assert.Equal(suite.T(), 1, n.minimum().value)

	m = &Node[int, int]{value: -1}
	n.left = &Node[int, int]{value: 0}
	n.left.left = m
	m.parent = n.left
	n.left.parent = n

	assert.Equal(suite.T(), m.value, n.minimum().value)
}

func TestNodeSuite(t *testing.T) {
	suite.Run(t, new(NodeSuite))
}

type RedBlackTreeSuite struct {
	suite.Suite
}

func (suite *RedBlackTreeSuite) TestIsValidRBTree() {
	tree := Make[int, int]()

	assert.False(suite.T(), tree.IsValidRBTree())

	tree.root = MakeNode(5, 5, red)
	assert.False(suite.T(), tree.IsValidRBTree())

	tree.root = MakeNode(5, 5, black)
	assert.True(suite.T(), tree.IsValidRBTree())

	tree.root.left = MakeNode(2, 2, red)
	tree.root.left.left = MakeNode(4, 4, red) // Violates BST and Red Property
	tree.root.left.left.parent = tree.root.left
	assert.False(suite.T(), tree.IsValidRBTree())

	tree.root.right = MakeNode(7, 7, red)
	tree.root.right.right = MakeNode(9, 9, black)
	tree.root.right.right.parent = tree.root.right

	tree.root.left.left.key = 1 // Still violates Red Property
	assert.False(suite.T(), tree.IsValidRBTree())

	tree.root.left.left.color = black
	assert.False(suite.T(), tree.IsValidRBTree()) // Violates Depth Property

	// Add missing leaves to make tree have same black depth
	tree.root.left.right = MakeNode(3, 3, black)
	tree.root.right.left = MakeNode(6, 6, black)
	assert.True(suite.T(), tree.IsValidRBTree())
}

func (suite *RedBlackTreeSuite) TestCreateEmpty() {
	tree := Make[int, string]()
	assert.Equal(suite.T(), 0, tree.Size())
}

func (suite *RedBlackTreeSuite) TestInsert() {
	tree := Make[int, int]()

	// Insert a new node with a value of 8
	tree.Insert(8, 1)

	// Expected tree view at this moment:
	// 				8(B)

	// Tree should contain single root element
	assert.Equal(suite.T(), 1, tree.Size())
	// Assert if element is stored and valid
	assert.Equal(suite.T(), 8, tree.root.key)
	assert.Equal(suite.T(), 1, tree.root.value)
	assert.Equal(suite.T(), black, tree.root.color)

	tree.Insert(18, 2)
	// Expected tree view at this moment:
	// 				8(B)
	//			  		 \
	// 					 18(R)

	// Check if tree correctly assigns new key node as a right branch of a root
	assert.Equal(suite.T(), 2, tree.Size())
	assert.Equal(suite.T(), 18, tree.root.right.key)
	assert.Equal(suite.T(), 2, tree.root.right.value)
	assert.Equal(suite.T(), red, tree.root.right.color)

	tree.Insert(5, 3)
	// Expected tree view at this moment:
	// 				8(B)
	//			  /		 \
	// 			5(R)	 18(R)

	assert.Equal(suite.T(), 3, tree.Size())

	// After this insertion tree should perform recolor
	tree.Insert(15, 4)
	// Expected tree view at this moment:
	// 				8(B)
	//			  /		 \
	// 			5(B)	 18(B)
	//					/
	//				  15(R)

	assert.Equal(suite.T(), 15, tree.root.right.left.key)
	assert.Equal(suite.T(), red, tree.root.right.left.color)
	assert.Equal(suite.T(), 5, tree.root.left.key)
	assert.Equal(suite.T(), 18, tree.root.right.key)
	assert.Equal(suite.T(), black, tree.root.left.color)
	assert.Equal(suite.T(), black, tree.root.right.color)

	// After this insertion tree should perform right rotation and recolor
	tree.Insert(17, 5)
	// Expected tree view at this moment:
	// 				8(B)
	//			  /		 \
	// 			5(B)	 17(B)
	//					/	  \
	//				  15(R)  18(R)

	assert.Equal(suite.T(), 8, tree.root.key)
	assert.Equal(suite.T(), black, tree.root.color)
	assert.Equal(suite.T(), 5, tree.root.left.key)
	assert.Equal(suite.T(), black, tree.root.left.color)
	assert.Equal(suite.T(), 17, tree.root.right.key)
	assert.Equal(suite.T(), black, tree.root.right.color)
	assert.Equal(suite.T(), 15, tree.root.right.left.key)
	assert.Equal(suite.T(), red, tree.root.right.left.color)
	assert.Equal(suite.T(), 18, tree.root.right.right.key)
	assert.Equal(suite.T(), red, tree.root.right.right.color)

	// After this insertion tree should perform recolor
	tree.Insert(25, 6)
	// Expected tree view at this moment:
	// 				8(B)
	//			  /		 \
	// 			5(B)	 17(R)
	//					/	  \
	//				  15(B)  18(B)
	//							\
	//							25(R)

	assert.Equal(suite.T(), 17, tree.root.right.key)
	assert.Equal(suite.T(), red, tree.root.right.color)
	assert.Equal(suite.T(), 15, tree.root.right.left.key)
	assert.Equal(suite.T(), black, tree.root.right.left.color)
	assert.Equal(suite.T(), 18, tree.root.right.right.key)
	assert.Equal(suite.T(), black, tree.root.right.right.color)
	assert.Equal(suite.T(), 25, tree.root.right.right.right.key)
	assert.Equal(suite.T(), red, tree.root.right.right.right.color)

	// After this insertion tree should perform left rotation and recolor
	tree.Insert(40, 7)
	// Expected tree view at this moment:
	// 				8(B)
	//			   /	\
	// 			5(B)	17(R)
	//				   /	 \
	//				  15(B)  25(B)
	//						 /	  \
	//					   18(R) 40(R)

	assert.Equal(suite.T(), 25, tree.root.right.right.key)
	assert.Equal(suite.T(), black, tree.root.right.right.color)
	assert.Equal(suite.T(), 18, tree.root.right.right.left.key)
	assert.Equal(suite.T(), red, tree.root.right.right.left.color)
	assert.Equal(suite.T(), 40, tree.root.right.right.right.key)
	assert.Equal(suite.T(), red, tree.root.right.right.right.color)

	// After this insertion tree should perform left rotation and recolor
	tree.Insert(80, 8)
	// Expected tree view at this moment:
	// 				17(B)
	//			  /	     \
	// 		   8(R)	    25(R)
	//		  /   \     /	\
	//	   5(B) 15(B) 18(B) 40(B)
	//						    \
	//					   		80(R)

	assert.Equal(suite.T(), 17, tree.root.key)
	assert.Equal(suite.T(), black, tree.root.color)

	assert.Equal(suite.T(), 8, tree.root.left.key)
	assert.Equal(suite.T(), red, tree.root.left.color)
	assert.Equal(suite.T(), 5, tree.root.left.left.key)
	assert.Equal(suite.T(), black, tree.root.left.left.color)
	assert.Equal(suite.T(), 15, tree.root.left.right.key)
	assert.Equal(suite.T(), black, tree.root.left.right.color)

	assert.Equal(suite.T(), 25, tree.root.right.key)
	assert.Equal(suite.T(), red, tree.root.right.color)
	assert.Equal(suite.T(), 18, tree.root.right.left.key)
	assert.Equal(suite.T(), black, tree.root.right.left.color)
	assert.Equal(suite.T(), 40, tree.root.right.right.key)
	assert.Equal(suite.T(), black, tree.root.right.right.color)
	assert.Equal(suite.T(), 80, tree.root.right.right.right.key)
	assert.Equal(suite.T(), red, tree.root.right.right.right.color)

	// After these insertions tree should perform multiple rotations and recolors
	tree.Insert(3, 9)
	tree.Insert(1, 10)
	tree.Insert(-3, 11)
	tree.Insert(60, 12)
	// Expected tree view at this moment:
	// 				17(B)
	//			  /	     \
	// 		   8(B)	    25(B)
	//		  /   \     /	\
	//	   3(R) 15(B) 18(B) 60(B)
	//	  /	  \			    /	\
	//	 1(B) 5(B)		40(R)	 80(R)
	//   /
	// -3(R)

	assert.Equal(suite.T(), 12, tree.Size())
	assert.True(suite.T(), tree.IsValidRBTree())

	assert.Equal(suite.T(), 17, tree.root.key)
	assert.Equal(suite.T(), black, tree.root.color)

	assert.Equal(suite.T(), 15, tree.root.left.right.key)
	assert.Equal(suite.T(), black, tree.root.left.right.color)
	assert.Equal(suite.T(), 5, tree.root.left.left.right.key)
	assert.Equal(suite.T(), black, tree.root.left.left.right.color)
	assert.Equal(suite.T(), -3, tree.root.left.left.left.left.key)
	assert.Equal(suite.T(), red, tree.root.left.left.left.left.color)

	assert.Equal(suite.T(), 18, tree.root.right.left.key)
	assert.Equal(suite.T(), black, tree.root.right.left.color)
	assert.Equal(suite.T(), 40, tree.root.right.right.left.key)
	assert.Equal(suite.T(), red, tree.root.right.right.left.color)
	assert.Equal(suite.T(), 80, tree.root.right.right.right.key)
	assert.Equal(suite.T(), red, tree.root.right.right.right.color)
}

func (suite *RedBlackTreeSuite) TestDelete() {
	data := []struct {
		key   int
		value int
	}{{8, 1}, {18, 2}, {5, 3}, {15, 4}, {17, 5}, {25, 6}, {40, 7}, {80, 8}, {3, 9}, {1, 10}, {-3, 11}, {60, 12}}

	tree := Make[int, int]()

	for _, el := range data {
		tree.Insert(el.key, el.value)
	}

	assert.Equal(suite.T(), len(data), tree.Size())
	assert.True(suite.T(), tree.IsValidRBTree())

	// Expected tree view at this moment:
	// 				17(B)
	// 			  /	     \
	// 		   8(B)	    25(B)
	// 		  /   \     /	\
	// 	   3(R) 15(B) 18(B) 60(B)
	// 	  /	  \			    /	\
	// 	 1(B) 5(B)		40(R)	 80(R)
	//   /
	// -3(R)

	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{-3, 1, 3, 5, 8, 15, 17, 18, 25, 40, 60, 80}, tree.Keys())
	assert.Equal(suite.T(), 8, tree.root.left.key)
	assert.Equal(suite.T(), black, tree.root.left.color)

	tree.Remove(8)
	// Expected tree view at this moment:
	// 				17(B)
	//			  /	     \
	// 		   3(B)    25(B)
	//		  /   \     /	\
	//	   1(R) 15(B) 18(B) 60(B)
	//	  /	  	/		    /	\
	//	-3(B)  5(R)	     40(R)	 80(R)

	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{-3, 1, 3, 5, 15, 17, 18, 25, 40, 60, 80}, tree.Keys())
	assert.Equal(suite.T(), 3, tree.root.left.key)
	assert.Equal(suite.T(), black, tree.root.left.color)

	tree.Remove(80)
	// Expected tree view at this moment:
	// 				17(B)
	//			  /	     \
	// 		   3(B)    25(B)
	//		  /   \     /	\
	//	   1(R) 15(B) 18(B) 60(B)
	//	  /	  	/		    /
	//	-3(B)  5(R)	     40(R)

	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{-3, 1, 3, 5, 15, 17, 18, 25, 40, 60}, tree.Keys())
	assert.Equal(suite.T(), 60, tree.root.right.right.key)
	assert.Equal(suite.T(), black, tree.root.right.right.color)

	tree.Remove(60)
	// Expected tree view at this moment:
	// 				17(B)
	//			  /	     \
	// 		   3(B)    25(B)
	//		  /   \     /	\
	//	   1(R) 15(B) 18(B) 40(B)
	//	  /	  	/
	//	-3(B) 5(R)

	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{-3, 1, 3, 5, 15, 17, 18, 25, 40}, tree.Keys())
	assert.Equal(suite.T(), 17, tree.root.key)
	assert.Equal(suite.T(), black, tree.root.color)
	assert.Equal(suite.T(), 25, tree.root.right.key)
	assert.Equal(suite.T(), black, tree.root.right.color)
	assert.Equal(suite.T(), 3, tree.root.left.key)
	assert.Equal(suite.T(), black, tree.root.left.color)

	tree.Remove(5)
	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{-3, 1, 3, 15, 17, 18, 25, 40}, tree.Keys())

	tree.Remove(1)
	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{-3, 3, 15, 17, 18, 25, 40}, tree.Keys())

	tree.Remove(-3)
	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{3, 15, 17, 18, 25, 40}, tree.Keys())

	tree.Remove(3)
	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{15, 17, 18, 25, 40}, tree.Keys())

	// Expected tree view at this moment:
	// 				25(B)
	//			  /	     \
	// 		   17(B)    40(B)
	//		  /   \
	//	   15(R) 18(R)

	assert.True(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), []int{15, 17, 18, 25, 40}, tree.Keys())

	tree.Remove(15)
	tree.Remove(17)
	tree.Remove(18)
	tree.Remove(25)
	tree.Remove(40)

	assert.False(suite.T(), tree.IsValidRBTree())
	assert.Equal(suite.T(), 0, tree.Size())
}

func (suite *RedBlackTreeSuite) TestWithSyntheticValues() {
	const n = 1000
	c := 0
	tree := Make[int64, string]()

	for i := 0; i < n; i++ {
		tree.Insert(int64(i), strconv.Itoa(i))
		if !tree.IsValidRBTree() {
			c++
		}
	}

	for i := n * 2; i > n; i-- {
		tree.Insert(int64(i), strconv.Itoa(i))
		if !tree.IsValidRBTree() {
			c++
		}
	}

	assert.Equal(suite.T(), n*2, tree.size)
	assert.Equal(suite.T(), 0, c, "non-valid tree occurence amount during insertion")
	c = 0

	for i := 0; i < n; i++ {
		tree.Remove(int64(i))
		if !tree.IsValidRBTree() {
			c++
		}
	}
	for i := n * 2; i > n; i-- {
		tree.Remove(int64(i))
		if !tree.IsValidRBTree() {
			c++
		}
	}
	assert.Equal(suite.T(), 0, tree.size)
	// 1 because evaluation after last deletion would result in non-valid tree
	assert.Equal(suite.T(), 1, c, "non-valid tree occurence amount during deletion")
}

func (suite *RedBlackTreeSuite) TestUpdate() {
	tree := Make[int, int]()

	tree.Insert(8, 1)

	assert.Equal(suite.T(), 1, tree.Size())
	assert.Equal(suite.T(), 1, tree.root.value)

	tree.Insert(8, 2)
	assert.Equal(suite.T(), 1, tree.Size())
	assert.Equal(suite.T(), 2, tree.root.value)
}

func (suite *RedBlackTreeSuite) TestSearch() {
	data := []struct {
		key   int
		value int
	}{{8, 1}, {18, 2}, {5, 3}, {15, 4}, {17, 5}, {25, 6}, {40, 7}, {80, 8}, {3, 9}, {1, 10}, {-3, 11}, {60, 12}}

	tree := Make[int, int]()

	_, exists := tree.Search(8)
	assert.False(suite.T(), exists)

	for _, el := range data {
		tree.Insert(el.key, el.value)
	}

	for _, el := range data {
		result, exists := tree.Search(el.key)
		assert.True(suite.T(), exists)
		assert.Equal(suite.T(), el.value, result)
	}

	_, exists = tree.Search(-99)
	assert.False(suite.T(), exists)
}

func (suite *RedBlackTreeSuite) TestPreorder() {
	expected := []struct {
		value int
		color Color
	}{{17, black}, {8, red}, {5, black}, {15, black}, {25, red}, {18, black}, {40, black}, {80, red}}
	tree := Make[int, int]()

	for _, e := range expected {
		tree.Insert(e.value, e.value)
	}

	assert.Equal(suite.T(), len(expected), tree.Size())

	for i, n := range tree.Preorder() {
		assert.Equal(suite.T(), expected[i].value, n.value, fmt.Sprintf("Failed assertion of node value %s, expected %d", &n, expected[i].value))
		assert.Equal(suite.T(), expected[i].color, n.color, fmt.Sprintf("Failed assertion of node color %s, expected opposite", &n))
	}
}

func TestRedBlackTreeSuite(t *testing.T) {
	suite.Run(t, new(RedBlackTreeSuite))
}

const (
	iMin = 5
	iMax = 15
)

func makeFilledTree(x, y int) *RedBlackTree[int, int] {
	tree := Make[int, int]()

	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			tree.Insert(i*j+j, 0)
		}
	}

	return tree
}

func BenchmarkRedBlackTreeInsert(b *testing.B) {
	for i := iMin; i <= iMax; i++ {
		n := 1 << i
		b.Run(fmt.Sprintf("size_%d", n), func(b *testing.B) {
			tree := Make[int, int]()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StartTimer()
				for j := 0; j < n; j++ {
					tree.Insert(i*j+j, 0)
				}
				b.StopTimer()
			}
		})
	}
}

func BenchmarkRedBlackTreeSearch(b *testing.B) {
	for i := iMin; i <= iMax; i++ {
		n := 1 << i
		b.Run(fmt.Sprintf("size_%d", n), func(b *testing.B) {
			tree := makeFilledTree(b.N, n)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StartTimer()
				for j := 0; j < n; j++ {
					tree.Search(i*j + j)
				}
				b.StopTimer()
			}
		})
	}
}

func BenchmarkRedBlackTreeDelete(b *testing.B) {
	for i := iMin; i <= iMax; i++ {
		n := 1 << i
		b.Run(fmt.Sprintf("size_%d", n), func(b *testing.B) {
			tree := makeFilledTree(b.N, n)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StartTimer()
				for j := 0; j < n; j++ {
					tree.Remove(i*j + j)
				}
				b.StopTimer()
			}
		})
	}
}
