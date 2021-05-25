package lexer

import (
	"bd/lexer/util"
	// "path/filepath"
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
func NewLexer(r *io.Reader, endToken string) *Lexer {
	s := util.NewStream(r, EndToken)
	return &Lexer {
		Stream: s,
		endToken: endToken,
	}
}


// 文件传入数据
func FileLexer(relatePath string) *Lexer {
	absPath, err := filepath(relatePath)
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
	return NewLexer(bytes.NewBufferString(source), endToken)
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
			l.PutBack(c)
			tokens = append(tokens, l.MakeVarOrKeyword())
			continue
		}

	}
}

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

func IsLetter()  {
	
}