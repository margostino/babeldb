package cli

import "github.com/margostino/babeldb/engine"

type step int
type Param int
type Type int
type Operator int

var TypeString = []string{
	"UnknownType",
	"Select",
	"Create",
}

var OperatorString = []string{
	"UnknownOperator",
	"Eq",
	"Ne",
	"Gt",
	"Lt",
	"Gte",
	"Lte",
}

type parser struct {
	i               int
	input           string
	step            step
	query           Query
	err             error
	nextUpdateField string
}

type Query struct {
	Type       Type
	Conditions []*Condition
	Params     map[Param]interface{}
	Fields     []string
	Aliases    map[string]string
	Solver     func(*engine.Engine, map[Param]interface{})
}

type Condition struct {
	// Operand1 is the left hand side operand
	Operand1 string
	// Operand1IsField determines if Operand1 is a literal or a field name
	Operand1IsField bool
	// Operator is e.g. "=", ">"
	Operator Operator
	// Operand1 is the right hand side operand
	Operand2 string
	// Operand2IsField determines if Operand2 is a literal or a field name
	Operand2IsField bool
}

const (
	stepType step = iota
	stepCreateSource
	stepCreateSourceFrom
	stepSelectFromSource
	stepCreateSourceWhen
	stepCreateSourceWith
	stepCreateSourceWithConfig
	stepSchedule
	stepWhere
	stepWhereField
	stepWhereOperator
	stepWhereValue
	stepWhereAnd
)

const (
	sourceName Param = iota
	sourceUrl
	schedule
	linkLevel
)

const (
	UnknownType Type = iota
	Select
	Create
)

const (
	// UnknownOperator is the zero value for an Operator
	UnknownOperator Operator = iota
	// Eq -> "="
	Eq
	// Ne -> "!="
	Ne
	// Gt -> ">"
	Gt
	// Lt -> "<"
	Lt
	// Gte -> ">="
	Gte
	// Lte -> "<="
	Lte
)
