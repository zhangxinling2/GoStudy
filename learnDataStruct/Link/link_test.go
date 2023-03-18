package link

import "testing"

type LinkNode struct {
	next *LinkNode
	data int
}

func (root *LinkNode) DeleteNode(n int) *LinkNode {
	if root.data == n {
		root = root.next
		return root
	}
	tmpPre := root
	tmp := root.next
	for tmp != nil {
		if tmp.data == n {
			tmpPre.next = tmp.next
			tmp.next = nil
			return root
		}
		tmpPre = tmp
		tmp = tmp.next
	}
	return root

}

func (root *LinkNode) AddNode(node *LinkNode) *LinkNode {
	tmp := root
	for tmp.next != nil {
		tmp = tmp.next
	}
	tmp.next = node
	return root
}

func (root *LinkNode) InsertNode(n int, node *LinkNode) *LinkNode {
	tmp := root
	if tmp.data == n {
		node.next = tmp.next
		tmp.next = node
		return root
	}
	for tmp.next != nil {
		tmp = tmp.next
		if tmp.data == n {
			node.next = tmp.next
			tmp.next = node
			return root
		}
	}
	tmp.next = node
	return root
}
func TestNode(t *testing.T) {
	n1 := &LinkNode{
		data: 1,
		next: nil,
	}
	n2 := &LinkNode{
		data: 2,
		next: nil,
	}
	n3 := &LinkNode{
		data: 3,
		next: nil,
	}
	n4 := &LinkNode{
		data: 4,
		next: nil,
	}
	n5 := &LinkNode{
		data: 5,
		next: nil,
	}
	n6 := &LinkNode{
		data: 6,
		next: nil,
	}
	n1.AddNode(n2)
	n1.AddNode(n3)
	n1.AddNode(n4)
	n1.InsertNode(4, n5)
	n1.AddNode(n6)
	root := n1
	for root != nil {
		t.Log(root.data)
		root = root.next
	}
}
