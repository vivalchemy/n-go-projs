package main

import (
	"fmt"
)

const (
	SIZE = 5
)

type Node struct {
	Left  *Node
	Val   string
	Right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type HashMap map[string]*Node

type Cache struct {
	Queue   Queue
	HashMap HashMap
}

func NewCache() *Cache {
	return &Cache{Queue: NewQueue(), HashMap: HashMap{}}
}

func (c *Cache) Check(str string) {
	node := &Node{}
	if val, ok := c.HashMap[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: str}
	}
	c.Add(node)
}

func (c *Cache) Add(node *Node) {
	if c.Queue.Length == SIZE {
		c.Remove(c.Queue.Tail.Left)
	}

	fmt.Println("adding node", node.Val)
	firstElement := c.Queue.Head.Right
	c.Queue.Head.Right = node
	firstElement.Left = node

	node.Left = c.Queue.Head
	node.Right = firstElement

	c.Queue.Length++

	c.HashMap[node.Val] = node
}

func (c *Cache) Remove(node *Node) *Node {
	fmt.Println("removing node", node.Val)

	leftNode := node.Left
	rightNode := node.Right

	leftNode.Right = rightNode
	rightNode.Left = leftNode

	c.Queue.Length--

	delete(c.HashMap, node.Val)

	return node
}

func (c *Cache) Display() {
	fmt.Print("[\033[0;33mhead<->")
	for node := c.Queue.Head.Right; node != c.Queue.Tail; node = node.Right {
		fmt.Print(node.Val, "<->")
	}
	fmt.Println("tail\033[0m]")
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head
	return Queue{Head: head, Tail: tail}
}

func main() {
	fmt.Println("starting cache")
	cache := NewCache()
	for _, word := range []string{"parrot", "avocado", "dragonfruit", "tree", "potato", "tomato", "tree", "dog"} {
		cache.Check(word)
		cache.Display()
	}

}
