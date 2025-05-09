package main

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/legamerdc/gas/ds"
)

func init() {
	registerSelector()
	registerModifier1()
	registerRunning()
	registerModifier2()
	registerAbility()
}

func registerSelector() {
	j1 := `{"id":1,"range":100,"count":1,"ally":false}`
	j2 := `{"id":2,"range":100,"count":1,"ally":true}`
	j3 := `{"id":3,"self":true}`
	var c1, c2, c3 SelectorConfig
	assert(jsoniter.Unmarshal([]byte(j1), &c1))
	assert(jsoniter.Unmarshal([]byte(j2), &c2))
	assert(jsoniter.Unmarshal([]byte(j3), &c3))
	ds.RegisterProxy[selectorConfigI](&c1)
	ds.RegisterProxy[selectorConfigI](&c2)
	ds.RegisterProxy[selectorConfigI](&c3)
}

func registerModifier1() {
	j1 := `{"id":1,"type":1,"value":3}` // 伤害3
	j2 := `{"id":2,"type":2,"value":2}` // 治疗2
	var c1, c2 ModifierConfig
	assert(jsoniter.Unmarshal([]byte(j1), &c1))
	assert(jsoniter.Unmarshal([]byte(j2), &c2))
	ds.RegisterProxy[modifierConfigI](&c1)
	ds.RegisterProxy[modifierConfigI](&c2)
}

func registerModifier2() {
	j3 := `{"id":3,"type":3,"running":1}` // 伤害
	j4 := `{"id":4,"type":3,"running":2}` // 治疗
	var c3, c4 ModifierConfig
	assert(jsoniter.Unmarshal([]byte(j3), &c3))
	assert(jsoniter.Unmarshal([]byte(j4), &c4))
	ds.RegisterProxy[modifierConfigI](&c3)
	ds.RegisterProxy[modifierConfigI](&c4)
}

func registerRunning() {
	j1 := `{"id":1,"gap":300,"duration":901,"selector":1,"modifier":1,"repeat":true}`
	j2 := `{"id":2,"gap":250,"duration":1001,"selector":2,"modifier":2,"repeat":true}`
	var c1, c2 ChainLikeRunningConfig
	assert(jsoniter.Unmarshal([]byte(j1), &c1))
	assert(jsoniter.Unmarshal([]byte(j2), &c2))
	ds.RegisterProxy[runningConfigI](&c1)
	ds.RegisterProxy[runningConfigI](&c2)
}

func registerAbility() {
	j1 := `{"id":1,"cd":2000,"selector":3,"modifier":3}`               // 伤害
	j2 := `{"id":2,"cd":2000,"selector":3,"modifier":4,"repeat":true}` // 治疗
	var c1, c2 SimpleAbilityConfig
	assert(jsoniter.Unmarshal([]byte(j1), &c1))
	assert(jsoniter.Unmarshal([]byte(j2), &c2))
	ds.RegisterProxy[abilityConfigI](&c1)
	ds.RegisterProxy[abilityConfigI](&c2)
}

func assert(e error) {
	if e != nil {
		panic(e)
	}
}
