package linkdoubledir

import "testing"

type LinkNode struct {
	previous *LinkNode
	next     *LinkNode
	data     int
}

func insertNode(root *LinkNode, newNode *LinkNode) *LinkNode {
	tmpNode := root

	// 整个链表是空的情况，新增
	if root == nil {
		return newNode
	}

	// 在链表的头，添加节点
	if root.data >= newNode.data {
		newNode.next = tmpNode
		tmpNode.previous = newNode

		return newNode
	}

	for {
		if tmpNode.next == nil {
			// 已经到头，追加节点即可
			tmpNode.next = newNode
			newNode.previous = tmpNode
			return root
		} else {
			if tmpNode.next.data >= newNode.data {
				// 找到位置，在此插入新节点
				newNode.previous = tmpNode
				newNode.next = tmpNode.next

				tmpNode.next = newNode
				newNode.next.previous = newNode

				return root
			}
		}
		tmpNode = tmpNode.next
	}
}

func deleteNode(root *LinkNode, v int) *LinkNode {
	if root == nil {
		return nil
	}

	if root.data == v {
		// 要删除的数据在第一个节点
		leftHand := root
		root = root.next

		leftHand.next = nil
		root.previous = nil

		// todo 需要解决只有一个节点的情况

		return root
	}

	tmpNode := root
	for {
		if tmpNode.next == nil {
			// 走到链表的尾部，仍然没有找到要删除的数据，直接返回原root
			return root
		} else {
			if tmpNode.next.data == v {
				// 找到节点，开始删除，删除完成后返回原root
				rightHand := tmpNode.next
				tmpNode.next = rightHand.next
				rightHand.next.previous = tmpNode

				// 清理掉右手上的link，保证GC正常回收
				rightHand.next = nil
				rightHand.previous = nil

				return root
			}
		}
		tmpNode = tmpNode.next
	}
}
func TestNode(t *testing.T) {
	n1 := &LinkNode{data: 1}
	n2 := &LinkNode{data: 5}
	n3 := &LinkNode{data: 10}
	insertNode(n1, n2)
	insertNode(n1, n3)
	deleteNode(n1, 5)
}
