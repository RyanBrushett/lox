package glox

type Stmt interface {
	Accept(VisitorStmt) (interface{}, error)
}

type VisitorStmt interface {
	visitExpressionStmt(*Expression) (interface{}, error)
	visitPrintStmt(*Print) (interface{}, error)
	visitVarStmt(*Var) (interface{}, error)
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) Stmt {
	return &Expression{expression}
}

func (e *Expression) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitExpressionStmt(e)
}

type Print struct {
	Expression Expr
}

func NewPrint(expression Expr) Stmt {
	return &Print{expression}
}

func (p *Print) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitPrintStmt(p)
}

type Var struct {
	Name        *Token
	Initializer Expr
}

func NewVar(name *Token, initializer Expr) Stmt {
	return &Var{name, initializer}
}

func (v *Var) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitVarStmt(v)
}
