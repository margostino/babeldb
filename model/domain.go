package model

import "fmt"

type Operator int32
type VarType int32
type QueryType int32

const (
	EqualOperator Operator = iota
	LikeOperator
	NotLikeOperator
	StringType VarType   = iota
	SelectType QueryType = iota
	CreateType
)

type QueryVar struct {
	Operator Operator
	VarType  VarType
	Field    string
	Value    interface{}
}

type ExpressionTree struct {
	Root *ExpressionNode
}

type ExpressionNode struct {
	Key     string
	VarType VarType
	Left    *ExpressionNode
	Right   *ExpressionNode
}

type Query struct {
	Source     string
	QueryType  QueryType
	Vars       []*QueryVar
	Expression *ExpressionTree
}

func (t *ExpressionTree) InOrderPrint(node *ExpressionNode) {
	if node != nil {
		t.InOrderPrint(node.Left)
		fmt.Printf("InOrder: %s\n", node.Key)
		t.InOrderPrint(node.Right)
	} else {
		//println("JKSBJKBKJBSJBS")
	}
}

func (t *ExpressionTree) Insert(key string) {
	if t.Root == nil {
		t.Root = &ExpressionNode{Key: key}
	} else {
		t.Root.Insert(key)
	}
}

func (n *ExpressionNode) Insert(key string) {
	if isLeaf(key) {
		if n.Left == nil {
			n.Left = &ExpressionNode{Key: key}
		} else if n.Right == nil {
			n.Right = &ExpressionNode{Key: key}
		}
	} else if isInode(key) {
		leafNode := *n
		n.Key = key
		n.Left = &leafNode
		n.Right = nil
	} else {
		if n.Right == nil {
			n.Right = &ExpressionNode{Key: key}
		} else {
			n.Right.Insert(key)
		}
	}

	//if isInode(key) {
	//	leafNode := *n
	//	n.Key = key
	//	n.Left = &leafNode
	//	n.Right = nil
	//} else if isLeaf(key) {
	//	if n.Left == nil {
	//		n.Left = &ExpressionNode{Key: key}
	//	} else if n.Right == nil {
	//		n.Right = &ExpressionNode{Key: key}
	//	}
	//} else {
	//	if n.Right == nil {
	//		n.Right = &ExpressionNode{Key: key}
	//	} else {
	//		n.Right.Insert(key)
	//	}
	//}
}
