package main

type Buff struct {
	ID       string
	Name     string
	Category int
	IsCanPly bool // false means buff's ply alway equal one or zero
	ply      int
}

func (b *Buff) cast() {
	//运行js脚本
}
