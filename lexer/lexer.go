package lexer

import (
	"bd/lexer/util"
	"path/filepath"
	"os"
	"io"
	"bytes"
)

const EndToken = "$"

type Lexer struct {
	// 此处嵌套了Stream
	*util.Stream
	endToken string
}

// 实例化一个语法分析器
func NewLexer(r io.Reader, endToken string) *Lexer {
	s := util.NewStream(r, EndToken)
	return &Lexer {
		Stream: s,
		endToken: endToken,
	}
}


// 文件传入数据
func FileLexer(relatePath string) *Lexer {
	absPath, err := filepath.Abs(relatePath)
	if err != nil {
		panic(err)
	}
	f, err := os.Open(absPath);
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return NewLexer(f, EndToken)
}

// 从字面量数据传入(test)
func Analyse(source string) []*Token  {
	return NewLexer(bytes.NewBufferString(source), EndToken).Analyse()
}

// Lexer函数方法

// 分析字节流/返回token的切片
func (this *Lexer) Analyse() []*Token {
	tokens := make([]*Token, 0)
	for this.HasNext() {
		// 获取一个字符(会resume)
		c := this.Next()
		if c == EndToken {
			break
		}
		// 查看下一个字符（不会resume）
		lookahead := this.Peek()

		// 空格|换行 跳过
		if c == " " || c == "\n" || c == "\t" {
			continue
		}

		
		if "/" == c {
			// 注释解析器实现(单行注释与多行注释 这里解释器将忽略跳过注释内容)
			// 查看下一个字符综合判定
			if "/" == lookahead {
				// 如果换行了，就不在是注释范围内了
				// 取消循环
				for this.HasNext() {
					if "\n" == this.Next() {
						break
					}
				}
				// 多行注释
			} else if "*" == lookahead {
				valid := false
				for this.HasNext() {
					p := this.Next()
					if("*" == p && "/" == this.Peek()) {
						this.Next()
						valid = true
						break
					}
				}
				if !valid {
					panic("mutiple comment syntax error!")
				}
				
			}
			// 如果不是注释，跳过(讨论是否该panic)
			continue
		}

		// 括号解析
		if c == "{" || c == "}" || c == "(" || c == ")" {
			tokens = append(tokens, NewToken(BRACKET, c))
			continue
		}
		
		// 字符串解析
		if c == `"` || c == `'` {
			// 如果碰到了则将该token放回cache
			// 全权由下一步函数解析
			this.PutBack(c)
			tokens = append(tokens, this.MakeString())
			continue
		}

		if IsLetter(c) {
			// 判定是否为正确的单词，放回，由后续函数接手
			this.PutBack(c)
			tokens = append(tokens, this.MakeVarOrKeyword())
			continue
		}
		// 数字字面量处理
		if IsNumber(c) {
			this.PutBack(c)
			tokens = append(tokens, this.MakeNumber())
		}

		// 优先判定表达式，若非表达式，再定义为单独符号
		// 3+5 .5 3 * 5
		if (c == "+" || c == "-" || c == ".") && IsNumber(lookahead) {
			var lastToken *Token = nil
			if len(tokens) > 0 {
				lastToken = tokens[len(tokens)-1]
			}

			if nil == lastToken || !lastToken.IsValue() || lastToken.IsOperator() {
				this.PutBack(c)
				tokens = append(tokens, this.MakeNumber())
				continue
			}
		}

		if IsOperator(c) {
			this.PutBack(c)
			tokens = append(tokens, this.MakeOp())
			continue
		}
		panic("unexpected character" + c)
	}
	return tokens
}

// 字符串处理函数
func (this *Lexer) MakeString() *Token{
	s := ""
	state := 0
	for this.HasNext() {
		c := this.Next()
		switch state {
			case 0:
				if c == `'` {
					state = 1
				} else {
					state = 2
				}
				s += c
			case 1: 
				if c == `'` {
					return NewToken(STRING, s + c)
				} else {
					s += c
				}
			case 2:
				if c == `"` {
					return NewToken(STRING, s + c)
				} else {
					s += c
				}
		}
	}
	panic("String syntax error")
}

// 单词处理函数（可能为关键字或变量）
func (this *Lexer) MakeVarOrKeyword() *Token {
	s := ""
	for this.HasNext() {
		lookahead := this.Peek()
		if IsLiteral(lookahead) {
			s += lookahead
			// this.Next()
		} else {
			break
		}
		this.Next()
	}

	if(IsKeyword(s)) {
		return NewToken(KEYWORD, s)
	}

	if "true" == s || "false" == s {
		return NewToken(BOOLEAN, s)
	}

	return NewToken(VARIABLE, s)
}

// 数字字面量处理函数(此处兼容了字面量表达式的处理)
func (this *Lexer) MakeNumber() *Token {
	state := 0
	s := ""
	for this.HasNext() {
		lookahead := this.Peek()
		switch state {
		// 初始化case
		case 0:
			if "0" == lookahead {
				state = 1
			} else if IsNumber(lookahead) {
				state = 2
			} else if "+" == lookahead || "-" == lookahead {
				state = 3
			} else if "." == lookahead {
				state = 5
			}
		// 如果数字以0开头，必须要跳过所有0
		case 1:
			if "0" == lookahead {
				state = 1
			} else if IsNumber(lookahead) {
				state = 2
			} else if lookahead == "." {
				state = 4
			} else {
				return NewToken(INTEGER, s)
			}
		// 如果是数字
		case 2:
			if IsNumber(lookahead) {
				state = 2
			} else if lookahead == "." {
				state = 4
			} else {
				return NewToken(INTEGER, s)
			}
		// 
		case 3:
			if IsNumber(lookahead) {
				state = 2
			} else if lookahead == "." {
				state = 5
			} else {
				panic("unexpected number " + lookahead)
			}
		
		case 4:
			if lookahead == "." {
				panic("unexpected number " + lookahead)
			} else if IsNumber(lookahead) {
				state = 20
			} else {
				return NewToken(FLOAT, s)
			}
		case 5:
			if IsNumber(lookahead) {
				state = 20
			} else {
				panic("unexpected number " + lookahead)
			}

		case 20:
			if IsNumber(lookahead) {
				state = 20
			} else if "." == lookahead {
				panic("unexpected number " + lookahead)
			} else {
				return NewToken(FLOAT, s)
			}
			
		}

		s += lookahead
		this.Next()
	}

	panic("make number fail")
}

func (this *Lexer) MakeOp() *Token {
	state := 0

	for this.HasNext() {
		// 这里不需要peek预检，因为不存在二义性，直接消耗字符即可
		lookahead := this.Next()
		switch state {
		case 0:
			switch lookahead {
			case "+":
				state = 1
			case "-":
				state = 2
			case "*":
				state = 3
			case `/`:
				state = 4
			case `>`:
				state = 5
			case `<`:
				state = 6
			case `=`:
				state = 7
			case `!`:
				state = 8
			case `&`:
				state = 9
			case `|`:
				state = 10
			case `^`:
				state = 11
			case `%`:
				state = 12
			case ",":
				return NewToken(OPERATOR, ",")
			case ";":
				return NewToken(OPERATOR, ";")
			}
		case 1:
			switch lookahead {
			case "+":
				return NewToken(OPERATOR, "++")
			case "=":
				return NewToken(OPERATOR, "+=")
			default: 
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "+")
			}
		case 2:
			switch lookahead {
			case "-":
				return NewToken(OPERATOR, "--")
			case "=":
				return NewToken(OPERATOR, "-=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "-")
			}
		case 3: 
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, "*=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "=")
			}
		case 4:
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, "/=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "/")
			}
		case 5:
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, ">=")
			case ">":
				return NewToken(OPERATOR, ">>")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, ">")
			}
		case 6:
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, "<=")
			case "<":
				return NewToken(OPERATOR, "<<")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "<")
			}
		case 7:
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, "==")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "=")
			}
		case 8:
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, "!=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "=")
			}
		case 9:
			switch lookahead {
			case "&":
				return NewToken(OPERATOR, "&&")
			case "=":
				return NewToken(OPERATOR, "&=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "&")
			}
		case 10:
			switch lookahead {
			case "|":
				return NewToken(OPERATOR, "||")
			case "=":
				return NewToken(OPERATOR, "|=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "|")
			}
		case 11:
			switch lookahead {
			case "^":
				return NewToken(OPERATOR, "^^")
			case "=":
				return NewToken(OPERATOR, "^=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "^")
			}
		case 12:
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, "%=")
			default:
				this.PutBack(lookahead)
				return NewToken(OPERATOR, "%")
			}
		}
	}
	panic("error with parse token")
} 

