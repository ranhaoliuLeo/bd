package lexer

import "fmt"

type TokenType int

// 枚举常量
const (
	KEYWORD  TokenType = 1
	VARIABLE TokenType = 2
	OPERATOR TokenType = 3
	BRACKET  TokenType = 4
	STRING   TokenType = 5
	FLOAT    TokenType = 6
	BOOLEAN  TokenType = 7
	INTEGER  TokenType = 8
)

// 定义默认的string输出结果
func (tt TokenType) String() string {
	switch tt {
		case KEYWORD:
			return "keyword"
		case VARIABLE:
			return "variable"
		case OPERATOR:
			return "operator"
		case BRACKET:
			return "bracket"
		case STRING:
			return "string`"
		case FLOAT:
			return "float"
		case BOOLEAN:
			return "boolean"
		case INTEGER:
			return "integer"
	}
	panic("Unkown Token-Type")
}

// 单词数据结构
type Token struct {
	Typ   TokenType
	Value string
}

func NewToken(t TokenType, v string) *Token {
	return &Token {
		Typ: t,
		Value: v,
	}
}

func (this *Token) IsVariable() bool  {
	return this.Typ == VARIABLE
}

func (this *Token) IsScalar() bool {
	return this.Typ == FLOAT || this.Typ == BOOLEAN || this.Typ == INTEGER || this.Typ == STRING
}

func (this *Token) IsNumber() bool {
	return this.Typ == INTEGER || this.Typ == FLOAT
}

func (this *Token) IsOperator() bool {
	return this.Typ == OPERATOR
}

func (this *Token) String() string {
	return fmt.Sprintf("Type:%v, Value:%s", this.Typ, this.Value)
}

func (this *Token) IsValue() bool {
	return this.IsVariable() || this.IsScalar()
}

func (this *Token) IsType() bool {
	switch this.Value {
		case "bool", "int", "float", "void", "string":
			return true
	}
	return false
}




