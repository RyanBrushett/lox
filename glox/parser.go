package glox

import "errors"

type parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.statement())
	}

	return statements
}

func (p *parser) parseExpression() Expr {
	return p.expression()
}

func (p *parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()

	// Error handling in here somewhere
}

func (p *parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *parser) printStatement() Stmt {
	value := p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		panic(err)
	}
	return NewPrint(value)
}

func (p *parser) expressionStatement() Stmt {
	expr := p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		panic(err)
	}

	return NewExpression(expr)
}

func (p *parser) expression() Expr {
	return p.ternary()
}

func (p *parser) varDeclaration() Stmt {
	name, err := p.consume(IDENTIFIER, "Expect a variable name.")
	if err != nil {
		panic(err)
	}
	var initializer Expr

	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return NewVar(name, initializer)
}

func (p *parser) ternary() Expr {
	var expr Expr
	expr = p.equality()

	if p.match(QUESTION) {
		leftOperator := p.previous() // grab the QUESTION since that's the left operator.
		middle := p.expression()
		rightOperator, err := p.consume(COLON, "Expect ':' after expression")
		if err != nil {
			panic(err)
		}
		right := p.expression()
		expr = NewTernary(expr, leftOperator, middle, rightOperator, right)
	}
	return expr
}

func (p *parser) equality() Expr {
	var expr Expr
	expr = p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *parser) comparison() Expr {
	var expr Expr
	expr = p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *parser) term() Expr {
	var expr Expr
	expr = p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *parser) factor() Expr {
	var expr Expr
	expr = p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}

	return p.primary()
}

func (p *parser) primary() Expr {
	var err error = nil
	if p.match(FALSE) {
		return NewLiteral(false)
	}
	if p.match(TRUE) {
		return NewLiteral(true)
	}
	if p.match(NIL) {
		return NewLiteral(nil)
	}

	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().literal)
	}

	if p.match(IDENTIFIER) {
		return NewVariable(p.previous())
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			panic(err)
		}
		return NewGrouping(expr)
	}

	pe := ParseError(p.peek(), errors.New("expect expression"))
	panic(pe)
}

func (p *parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *parser) previous() *Token {
	if p.current-1 < 0 {
		return nil
	}
	return p.tokens[p.current-1]
}

func (p *parser) advance() *Token {
	if p.isAtEnd() {
		return nil
	} else {
		p.current++
	}
	return p.previous()
}

func (p *parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *parser) match(tokenTypes ...TokenType) bool {
	for _, tt := range tokenTypes {
		if p.check(tt) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *parser) consume(tokenType TokenType, msg string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return nil, ParseError(p.peek(), errors.New(msg))
}

// Apparently we'll need this for later
//
// func (p *parser) synchronise() {
// 	p.advance()

// 	for !p.isAtEnd() {
// 		if p.previous().tokenType == SEMICOLON {
// 			return
// 		}

// 		switch p.peek().tokenType {
// 		case CLASS:
// 		case FUN:
// 		case VAR:
// 		case FOR:
// 		case IF:
// 		case WHILE:
// 		case PRINT:
// 		case RETURN:
// 			return
// 		}

// 		p.advance()
// 	}
// }
