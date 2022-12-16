package cli

import (
	"fmt"
	"regexp"
	"strings"
)

func Parse(input string) (Query, error) {
	qs, err := ParseMany([]string{input})
	if len(qs) == 0 {
		return Query{}, err
	}
	return qs[0], err
}

func ParseMany(sqls []string) ([]Query, error) {
	qs := []Query{}
	for _, sql := range sqls {
		q, err := parse(sql)
		if err != nil {
			return qs, err
		}
		qs = append(qs, q)
	}
	return qs, nil
}

func parse(input string) (Query, error) {
	parser := &parser{
		0,
		strings.TrimSpace(input),
		stepType,
		Query{},
		nil,
		""}
	return parser.parse()
}

//func (q *Query) solve() {
//	q.Solver(q.Params)
//}

func (p *parser) parse() (Query, error) {
	q, err := p.doParse()
	p.err = err
	if p.err == nil {
		p.err = p.validate()
	}
	p.logError()
	return q, p.err
}

func (p *parser) doParse() (Query, error) {
	for {
		if p.i >= len(p.input) {
			return p.query, p.err
		}
		switch p.step {
		case stepType:
			switch strings.ToUpper(p.peek()) {
			case "CREATE SOURCE":
				p.query.Type = Create
				p.query.Params = make(map[Param]interface{})
				p.step = stepCreateSource
				p.query.Solver = createSource
				p.pop()

				inputSourceName := p.pop()
				if !isIdentifierOrAsterisk(inputSourceName) {
					return p.query, fmt.Errorf("at CREATE SOURCE: expected source name")
				}
				p.query.Params[sourceName] = inputSourceName

				from := p.pop()
				if strings.ToUpper(from) != "FROM" {
					return p.query, fmt.Errorf("at CREATE SOURCE ${name}: expected FROM")
				}

				inputSourceUrl := p.pop()
				if len(inputSourceUrl) == 0 {
					return p.query, fmt.Errorf("at CREATE SOURCE ${name} FROM ${url}: expected valid URL")
				}
				p.query.Params[sourceUrl] = inputSourceUrl

				if strings.ToUpper(p.pop()) == "WHEN" {
					cronExp := p.pop()
					p.query.Params[schedule] = cronExp
				}
				param := p.pop()
				if param == ";" {
					continue
				}

			default:
				return p.query, fmt.Errorf("invalid query type")
			}

		case stepWhereOperator:
			operator := p.peek()
			currentCondition := p.query.Conditions[len(p.query.Conditions)-1]
			switch operator {
			case "=":
				currentCondition.Operator = Eq
			case ">":
				currentCondition.Operator = Gt
			case ">=":
				currentCondition.Operator = Gte
			case "<":
				currentCondition.Operator = Lt
			case "<=":
				currentCondition.Operator = Lte
			case "!=":
				currentCondition.Operator = Ne
			default:
				return p.query, fmt.Errorf("at WHERE: unknown operator")
			}
			p.query.Conditions[len(p.query.Conditions)-1] = currentCondition
			p.pop()
			p.step = stepWhereValue
		case stepWhereValue:
			currentCondition := p.query.Conditions[len(p.query.Conditions)-1]
			identifier := p.peek()
			if isIdentifier(identifier) {
				currentCondition.Operand2 = identifier
				currentCondition.Operand2IsField = true
			} else {
				quotedValue, ln := p.peekQuotedStringWithLength()
				if ln == 0 {
					return p.query, fmt.Errorf("at WHERE: expected quoted value")
				}
				currentCondition.Operand2 = quotedValue
				currentCondition.Operand2IsField = false
			}
			p.query.Conditions[len(p.query.Conditions)-1] = currentCondition
			p.pop()
			p.step = stepWhereAnd
		case stepWhereAnd:
			andRWord := p.peek()
			if strings.ToUpper(andRWord) != "AND" {
				return p.query, fmt.Errorf("expected AND")
			}
			p.pop()
			p.step = stepWhereField
		}
	}
}

func (p *parser) peek() string {
	peeked, _ := p.peekWithLength()
	return peeked
}

func (p *parser) pop() string {
	peeked, len := p.peekWithLength()
	p.i += len
	p.popWhitespace()
	return peeked
}

func (p *parser) popWhitespace() {
	for ; p.i < len(p.input) && p.input[p.i] == ' '; p.i++ {
	}
}

var reservedWords = []string{
	"(", ")", ">=", "<=", "!=", ",", "=", ">", "<", ";", "SELECT", "CREATE SOURCE", "VALUES", "UPDATE", "DELETE FROM",
	"WHERE", "FROM", "SET", "AS", "WHEN", "SCHEDULE",
}

func (p *parser) peekWithLength() (string, int) {
	if p.i >= len(p.input) {
		return "", 0
	}
	for _, rWord := range reservedWords {
		token := strings.ToUpper(p.input[p.i:min(len(p.input), p.i+len(rWord))])
		if token == rWord {
			return token, len(token)
		}
	}
	if p.input[p.i] == '\'' { // Quoted string
		return p.peekQuotedStringWithLength()
	}
	return p.peekIdentifierWithLength()
}

func (p *parser) peekQuotedStringWithLength() (string, int) {
	if len(p.input) < p.i || p.input[p.i] != '\'' {
		return "", 0
	}
	for i := p.i + 1; i < len(p.input); i++ {
		if p.input[i] == '\'' && p.input[i-1] != '\\' {
			return p.input[p.i+1 : i], len(p.input[p.i+1:i]) + 2 // +2 for the two quotes
		}
	}
	return "", 0
}

func (p *parser) peekIdentifierWithLength() (string, int) {
	for i := p.i; i < len(p.input); i++ {
		if matched, _ := regexp.MatchString(`[a-zA-Z0-9_*]`, string(p.input[i])); !matched {
			return p.input[p.i:i], len(p.input[p.i:i])
		}
	}
	return p.input[p.i:], len(p.input[p.i:])
}

func (p *parser) validate() error {
	if len(p.query.Conditions) == 0 && p.step == stepWhereField {
		return fmt.Errorf("at WHERE: empty WHERE clause")
	}
	if p.query.Type == UnknownType {
		return fmt.Errorf("query type cannot be empty")
	}
	//if len(p.query.Conditions) == 0 && p.query.Type == Delete {
	//	return fmt.Errorf("at WHERE: WHERE clause is mandatory for UPDATE & DELETE")
	//}
	for _, c := range p.query.Conditions {
		if c.Operator == UnknownOperator {
			return fmt.Errorf("at WHERE: condition without operator")
		}
		if c.Operand1 == "" && c.Operand1IsField {
			return fmt.Errorf("at WHERE: condition with empty left side operand")
		}
		if c.Operand2 == "" && c.Operand2IsField {
			return fmt.Errorf("at WHERE: condition with empty right side operand")
		}
	}
	return nil
}

func (p *parser) logError() {
	if p.err == nil {
		return
	}
	fmt.Println(p.input)
	fmt.Println(strings.Repeat(" ", p.i) + "^")
	fmt.Println(p.err)
}

func isIdentifier(s string) bool {
	for _, rw := range reservedWords {
		if strings.ToUpper(s) == rw {
			return false
		}
	}
	matched, _ := regexp.MatchString("[a-zA-Z_][a-zA-Z_0-9]*", s)
	return matched
}

func isIdentifierOrAsterisk(s string) bool {
	return isIdentifier(s) || s == "*"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
