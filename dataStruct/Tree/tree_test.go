package tree

import (
	"fmt"
	"testing"
)

type Node struct {
	left  *Node
	right *Node
	root  *Node
	data  int
}

func buildTree() *TreeNode {
	n1 := &TreeNode{val: 51}
	n2 := &TreeNode{val: 35}
	n3 := &TreeNode{val: 65}
	n1.left = n2
	n1.right = n3
	return n1
}
func insertNode(root *Node, newNode *Node) *Node {
	if root == nil {
		return newNode
	}
	if root.data == newNode.data {
		return root
	}
	if root.data > newNode.data {
		if root.left == nil {
			root.left = newNode
			newNode.root = root
		} else {
			insertNode(root.left, newNode)
		}
	}
	if root.data < newNode.data {
		if root.right == nil {
			root.right = newNode
			newNode.root = root
		} else {
			insertNode(root.right, newNode)
		}
	}
	return root
}

// 中间用于帮助理解删除的是最后叶子节点的情况
func deleteNodeLeaf(root *Node, v int) *Node {
	leftRoot := root
	if leftRoot.data == v && leftRoot.left == nil && leftRoot.right == nil {
		leftRoot = leftRoot.root
		right := root
		if leftRoot.left == right {
			// 删除左边叶子
			leftRoot.left = nil
			right.root = nil
			return leftRoot
		} else {
			// 删除右边叶子
			leftRoot.right = nil
			right.root = nil
			return leftRoot
		}
	}
	return root
}

func deleteNode(root *Node, v int) *Node {
	if v < root.data {
		deleteNode(root.left, v)
	} else if v > root.data {
		deleteNode(root.right, v)
	} else {
		// 现在root指向要删除的节点
		leftNextGen := findNextGenFromLeft(root.left)
		rightNextGen := findNextGenFromRight(root.right)
		if leftNextGen == nil && rightNextGen == nil {
			// 现在要删除的是叶子结点，即最底部的节点
			top := root.root
			down := root
			if top.left == down {
				// 表示是左子树
				top.left = nil
				down.root = nil
				return nil
			} else {
				// 表示是右子树
				top.right = nil
				down.root = nil
				return nil
			}
		} else if leftNextGen != nil {
			root.data = leftNextGen.data
			deleteNode(leftNextGen, leftNextGen.data)
		} else if rightNextGen != nil {
			root.data = rightNextGen.data
			deleteNode(rightNextGen, rightNextGen.data)
		}
	}
	return root
}
func findNextGenFromLeft(root *Node) *Node {
	if root == nil {
		return nil
	}
	tmpNode := root
	for {
		if tmpNode.right != nil {
			tmpNode = tmpNode.right
		} else {
			break
		}
	}
	return tmpNode
}

func findNextGenFromRight(root *Node) *Node {
	if root == nil {
		return nil
	}
	tmpNode := root
	for {
		if tmpNode.left != nil {
			tmpNode = tmpNode.left
		} else {
			break
		}
	}
	return tmpNode
}

type TreeNode struct {
	val   int
	left  *TreeNode
	right *TreeNode
}

func insertTreeNode(t *TreeNode, newNode *TreeNode) {
	if t == nil {
		t = newNode
	}
	if t.val < newNode.val {
		if t.right == nil {
			t.right = newNode
		} else {
			insertTreeNode(t.right, newNode)
		}
	} else if t.val > newNode.val {
		if t.left == nil {
			t.left = newNode
		} else {
			insertTreeNode(t.left, newNode)
		}
	} else {
		return
	}
}

func deleteTreeNode(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return nil
	}
	if root.val < val {
		root.right = deleteTreeNode(root.right, val)
		return root
	}
	if root.val > val {
		root.left = deleteTreeNode(root.left, val)
		return root
	}
	//这两个把两个都为nil的也处理了
	if root.left == nil {
		return root.right
	}
	if root.right == nil {
		return root.left
	}
	leaf := findLeft(root.left)
	root.left = deleteTreeNode(root.left, leaf.val)
	leaf.left = root.left
	leaf.right = root.right
	root = leaf
	return root
}
func findLeft(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	for root.right != nil {
		root = root.right
	}
	return root
}
func TestTree(t *testing.T) {
	root := buildTree()
	insertTreeNode(root, &TreeNode{val: 43})
	insertTreeNode(root, &TreeNode{val: 68})
	n := deleteTreeNode(root, 43)
	fmt.Println(n)
}
