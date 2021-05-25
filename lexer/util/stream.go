package util

import (
	"bufio"
	"container/list"
	"io"
)

type Stream struct {
	scanner    *bufio.Scanner
	queueCache *list.List
	endToken   string
	isEnd      bool
}

func NewStream(r io.Reader, endToken string) *Stream {
	// 将基础的io.reader转换为scanner-reader
	scanner := bufio.NewScanner(r)
	// 以utf-8为界限划分字节流
	scanner.Split(bufio.ScanRunes)
	// 返回stream实例
	return &Stream {
		scanner: scanner,
		queueCache: list.New(),
		endToken: endToken,
		isEnd: false,
	}
}


// 如果有缓存，读取缓存列表内下一个token
// 如果没缓存，则调用扫描器进行扫描
// 扫描器出错，则停止，返回自定义的终止符号
// 消化token
func (this *Stream) Next() string {
	if this.queueCache.Len() != 0 {
		e := this.queueCache.Front()
		return this.queueCache.Remove(e).(string)
	}

	if this.scanner.Scan() {
		return this.scanner.Text()
	}

	this.isEnd = true
	return this.endToken
}

// 从几个方面判断是否还有后续字节流
func (this *Stream) HasNext() bool {
	if this.queueCache.Len() != 0 {
		return true
	}

	if this.scanner.Scan() {
		this.queueCache.PushBack(this.scanner.Text())
		return true
	}

	if !this.isEnd {
		return true
	}

	return false
}

// 主要是读取，不会消耗Cache
// 还能有优化空间
func (this *Stream) Peek() string {
	if this.queueCache.Len() != 0 {
		return this.queueCache.Front().Value.(string)
	}

	if this.scanner.Scan() {
		token := this.scanner.Text()
		this.queueCache.PushBack(token)
		return token
	}

	return this.endToken
}

func (this *Stream) PutBack(token string) {
	this.queueCache.PushFront(token)
}