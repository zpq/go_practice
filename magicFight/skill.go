package main

import (
	"github.com/robertkrimen/otto"
)

type Skill struct {
	ID          string
	Name        string
	ScriptID    string // 技能脚本的路径
	Description string
}

func (s *Skill) Use(b *Battle) error {
	jsCode, err := loadJs(s.ScriptID)
	if err != nil {
		return err
	}
	vm := otto.New()
	jsBattle, err := vm.ToValue(b)
	if err != nil {
		return err
	}
	vm.Set("battle", jsBattle)
	_, err = vm.Run(jsCode)
	if err != nil {
		return err
	}
	return nil
}
