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
	*util.Stream
	endToken string
}

// 实例化一个语法分析器
func NewLexer(r *io.Reader, endToken string) *Lexer {
	s := util.NewStream(r, EndToken)
	return &Lexer {
		Stream: s,
		endToken: endToken
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
func Analyse(source string)   {
	return NewLexer(bytes.NewBufferString(source), endToken)
}

// Lexer函数方法