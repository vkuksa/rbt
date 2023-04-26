Implementation of generic Red-Black tree, with a keys constrained over orderable types

Operations are based on Introduction to algorithms book, but modified to omit sentinel node usage

Tree has following exported functions:
Insert(k K, v V)
Search(k K) (value V, exists bool)
Remove(k K)
Keys() []K
Preorder() []Node[K,V]
Size() int
IsValidRBTree() bool

