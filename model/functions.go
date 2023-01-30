package model

import (
	"fmt"
	"github.com/margostino/babeldb/common"
	"strings"
)

func (t *ExpressionTree) InOrderPrint(node *ExpressionNode) {
	if node != nil {
		t.InOrderPrint(node.Left)
		fmt.Printf("InOrder: %s\n", node.GetKey())
		t.InOrderPrint(node.Right)
	}
}

func (t *ExpressionTree) Insert(key string) {
	if t.Root == nil {
		t.Root = &ExpressionNode{Key: &ExpressionNodeKey{Value: key}}
	} else {
		t.Root.Insert(key)
	}
}

func (t *ExpressionTree) GetParamNode(node *ExpressionNode) *ExpressionNode {

	if node.isParamValue() {
		return node
	}

	return t.GetParamNode(node.Right)
}

func (t *ExpressionTree) Match(node *ExpressionNode, section *Section) bool {
	var match bool
	if node.IsComparisonOperatorNode() {
		text := strings.ToLower(section.Text)
		value := common.NewString(node.Right.GetKey()).
			ToLower().
			ReplaceAll("%", "").
			Value()
		// TODO: like operator logic
		switch node.GetOperator() {
		case EqualOperator:
			if value == "*" {
				match = true
			} else {
				match = text == value
			}
		case LikeOperator:
			match = strings.Contains(text, value)
		case NotLikeOperator:
			match = !strings.Contains(text, value)
		}

		return match
	}

	leftMatch := t.Match(node.Left, section)
	rightMatch := t.Match(node.Right, section)

	if node.GetOperator() == AndOperator {
		match = leftMatch && rightMatch
	} else if node.GetOperator() == OrOperator {
		match = leftMatch || rightMatch
	}

	return match
}

func (n *ExpressionNode) GetKey() string {
	return n.Key.Value
}

func (n *ExpressionNode) GetType() Type {
	return n.Key.Type
}

func (n *ExpressionNode) GetOperator() Operator {
	return n.Key.Operator
}

func (n *ExpressionNode) IsComparisonOperatorNode() bool {
	return isComparisonOperator(n.GetKey())
}

func (n *ExpressionNode) IsInode() bool {
	return isInode(n.GetKey())
}

func (n *ExpressionNode) IsValueFieldNode() bool {
	return !isFieldNode(n.GetKey()) && !n.IsInode()
}

func (n *ExpressionNode) IsLogicalOperatorNode() bool {
	return isLogicalOperator(n.GetKey())
}

func (n *ExpressionNode) isParamValue() bool {
	return len(n.GetKey()) == 2 && n.GetKey()[0:1] == ":"
}

func (n *ExpressionNode) Insert(key string) {
	if isLeaf(key) {
		if n.Left == nil {
			n.Left = &ExpressionNode{Key: &ExpressionNodeKey{Value: key}}
		} else if n.Right == nil {
			n.Right = &ExpressionNode{Key: &ExpressionNodeKey{Value: key}}
			if isComparisonOperator(n.GetKey()) && n.Left.GetKey() == TypeField {
				n.Right.Key.Type = TokenType
			}
		} else {
			n.Right.Insert(key)
		}
	} else {
		operator := GetOperator(key)

		if isLogicalOperator(key) {
			leafNode := *n
			node := &ExpressionNode{
				Key: &ExpressionNodeKey{
					Value:    key,
					Operator: operator,
				},
				Left:  &leafNode,
				Right: nil,
			}
			n.Key = node.Key
			n.Left = node.Left
			n.Right = node.Right
		} else {
			leafNode := *n
			if isLogicalOperator(n.GetKey()) {
				n.Right = &ExpressionNode{
					Key: &ExpressionNodeKey{
						Value:    key,
						Operator: operator,
					},
					Left:  n.Right,
					Right: nil,
				}
			} else {
				n.Key = &ExpressionNodeKey{
					Value:    key,
					Operator: operator,
				}
				n.Left = &leafNode
				n.Right = nil
			}
		}
	}

}
