# Implementation of generic Red-Black tree, with a keys constrained over orderable types

Red-Black tree is a binary search tree in which every node is colored with either red or black. It is a type of self balancing binary search tree. It has a good efficient worst case running time complexity.

Operations are based on the Introduction to Algorithms, but modified to omit sentinel node usage

rbt package has following exported functions:
```
Make[K, V]()
Insert(k K, v V)
Search(k K) (value V, exists bool)
Remove(k K)
Keys() []K
Traverse(func(k K, v V))
Size() int
```

Example of functions usage:

```go
package main

import (
	"fmt"

	"github.com/vkuksa/rbt"
)

func main() {
	tree := rbt.Make[int, int]()
	tree.Insert(1, 1)
	tree.Insert(-3, 2)
	tree.Insert(6, 3)
	tree.Insert(0, 4)
	tree.Insert(8, 5)

	if (!tree.IsValidRBTree()) {
		fmt.Printf("Tree size: %d\n", tree.Size())
		return
	}

	fmt.Printf("Tree size: %d\n", tree.Size())
	fmt.Printf("Tree keys: %d\n", tree.Keys())

	if v, found := tree.Search(5); found {
		fmt.Printf("Found value for a key %d: %d\n", 5, v)
	}

	tree.Traverse(func(k int, v int) bool {
		fmt.Printf("%d : %d\n", k, v)
		return true
	})
	tree.Remove(6)
}
```
