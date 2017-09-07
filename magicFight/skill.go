package main

import "github.com/robertkrimen/otto"

type Skill struct {
	ID          string
	Name        string
	ScriptID    string // 技能脚本的路径
	Description string
}

func (s *Skill) Use(ca *GameBattler, b *BattleControl) (error, string) {
	jsCode, err := loadJs(scriptSkillPrefix + s.ScriptID + ".js")
	if err != nil {
		return err, ""
	}
	vm := otto.New()
	jsBattleControl, err := vm.ToValue(b)
	if err != nil {
		return err, ""
	}
	jsCa, err := vm.ToValue(ca)
	if err != nil {
		return err, ""
	}
	vm.Set("bc", jsBattleControl)
	vm.Set("ca", jsCa)
	_, err = vm.Run(jsCode)
	if err != nil {
		return err, ""
	}
	result, err := vm.Get("result")
	if err != nil {
		return err, ""
	}
	r, err := result.Export()
	if err != nil || r == nil {
		return err, ""
	}
	return nil, r.(string)
}
