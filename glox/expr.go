package glox

type Expr interface {
	Accept(VisitorExpr) (interface{}, error)
}

type VisitorExpr interface {
	visitBinaryExpr(*Binary) (interface{}, error)
	visitGroupingExpr(*Grouping) (interface{}, error)
	visitLiteralExpr(*Literal) (interface{}, error)
	visitUnaryExpr(*Unary) (interface{}, error)
	visitTernaryExpr(*Ternary) (interface{}, error)
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewBinary(left Expr, operator *Token, right Expr) Expr {
	return &Binary{left, operator, right}
}

func (b *Binary) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{expression}
}

func (g *Grouping) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) Expr {
	return &Literal{value}
}

func (l *Literal) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitLiteralExpr(l)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func NewUnary(operator *Token, right Expr) Expr {
	return &Unary{operator, right}
}

func (u *Unary) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitUnaryExpr(u)
}

type Ternary struct {
	Left          Expr
	LeftOperator  *Token
	Middle        Expr
	RightOperator *Token
	Right         Expr
}

func NewTernary(left Expr, lo *Token, middle Expr, ro *Token, right Expr) Expr {
	return &Ternary{left, lo, middle, ro, right}
}

func (t *Ternary) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitTernaryExpr(t)
}
